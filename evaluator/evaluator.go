package evaluator

import (
	"Monkey/ast"
	"Monkey/object"
	"Monkey/options"
	"Monkey/parser"
	"Monkey/token"
	"fmt"
	"strings"
)

var (
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
	NULL  = &object.Null{}
)

// Create a new error node
func NewError(data *token.TokenData, format string, a ...interface{}) *object.Error {
	return &object.Error{
		Message:   fmt.Sprintf(format, a...),
		TokenData: data,
	}
}

// Modified here
// Is the object an error
func CheckError(obj object.Object) bool {
	if obj == nil {
		return false
	}
	isError := obj.Type() == object.ERROR_OBJ

	if isError {
		// If Fatal errors is set, we stop the exec
		if options.FatalErrors {
			fmt.Println(obj.Inspect())
			return true
		} else {
			// Else we treat error as a valid value
			return false
		}
	}

	return false
}

// Master function to determine if an object is true or not
func IsTruthful(obj object.Object) bool {
	switch {
	case obj == FALSE, obj == NULL:
		return false
	case obj.Type() == object.INTEGER_OBJ:
		integer := obj.(*object.Integer)
		return integer.Value != 0
	default:
		return true
	}
}

// Master Eval Function
func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {

	case *ast.Program:
		return EvalProgram(node, env)

	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)

	case *ast.PrintExpressionStatement:
		return EvalPrintExpressionStatement(node.Token, node.Expression, env)

	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}

	case *ast.Boolean:
		return NativeBoolToBooleanObject(node.Value)

	case *ast.PrefixExpression:
		right := Eval(node.Right, env)
		if CheckError(right) {
			return right
		}
		return EvalPrefixExpression(node.Operator, right, node.Token)

	case *ast.InfixExpression:
		// We need to short circuit AND or OR gates
		return EvalInfixExpression(node, env)

	case *ast.BlockStatement:
		return EvalBlockStatement(node, env)

	case *ast.IfExpression:
		return EvalIfExpression(node, env)

	case *ast.ReturnStatement:
		val := Eval(node.ReturnValue, env)
		if CheckError(val) {
			return val
		}
		return &object.ReturnValue{Value: val}

	case *ast.LetStatement:
		val := Eval(node.Value, env)
		if CheckError(val) {
			return val
		}

		env.Store(node.Name.Value, val)

	case *ast.Identifier:
		return EvalIdentifier(node, env)

	case *ast.FunctionLiteral:
		params := node.Parameters
		body := node.Body
		return &object.Function{Parameters: params, Env: env, Body: body}

	case *ast.CallExpression:
		function := Eval(node.Function, env)
		if CheckError(function) {
			return function
		}

		args := EvalExpressions(node.Arguments, env)
		if len(args) == 1 && CheckError(args[0]) {
			return args[0]
		}
		return ApplyFunction(node.Token, function, args)

	case *ast.StringLiteral:
		return &object.String{Value: node.Value}

	case *ast.ArrayLiteral:
		elements := EvalExpressions(node.Elements, env)
		if len(elements) == 1 && CheckError(elements[0]) {
			return elements[0]
		}
		return &object.Array{
			Elements: elements,
		}

	case *ast.IndexExpression:
		left := Eval(node.Left, env)
		if CheckError(left) {
			return left
		}

		start := Eval(node.Start, env)
		if CheckError(start) {
			return start
		}

		end := Eval(node.End, env)
		if CheckError(end) {
			return end
		}

		return EvalIndexExpression(left, start, end, node.Token, node.HasRange)
	}

	return NULL
}

// Eval Indexing
func EvalIndexExpression(exp object.Object, start object.Object, end object.Object, token token.Token, hasRange bool) object.Object {
	switch exp.(type) {
	case *object.Array:
		return EvalArrayIndexExpression(exp, start, end, token, hasRange)
	default:
		return NewError(token.ToTokenData(), "index operator not support for `%s`", exp.Type())
	}
}

func IsIndexInRange(index int64, length int64) bool {
	return index >= 0 && index < length
}

// Eval Array indexing
func EvalArrayIndexExpression(array object.Object, start object.Object, end object.Object, token token.Token, hasRange bool) object.Object {
	arrayObj := array.(*object.Array)
	length := int64(len(arrayObj.Elements))

	// Four Options
	switch {
	case !hasRange:
		s := int64(start.(*object.Integer).Value)
		if s < 0 {
			s = length + s
		}
		if !IsIndexInRange(s, length) {
			return NewError(token.ToTokenData(), "index out of range. got=%d, expected=%d-%d",
				s, 0, length-1)
		}
		return arrayObj.Elements[s]

	case hasRange:
		{
			var startIndex int64
			var endIndex int64

			switch {
			case start.Type() == object.INTEGER_OBJ && end.Type() == object.NULL_OBJ:
				startIndex = int64(start.(*object.Integer).Value)
				endIndex = length
			case start.Type() == object.NULL_OBJ && end.Type() == object.INTEGER_OBJ:
				startIndex = 0
				endIndex = int64(end.(*object.Integer).Value)
			case start.Type() == object.INTEGER_OBJ && end.Type() == object.INTEGER_OBJ:
				startIndex = int64(start.(*object.Integer).Value)
				endIndex = int64(end.(*object.Integer).Value)
			default:
				// Full Range
				return arrayObj
			}

			if startIndex < 0 {
				startIndex = length + startIndex
			}
			if endIndex < 0 {
				endIndex = length + endIndex
			}

			if !IsIndexInRange(startIndex, length+1) {
				return NewError(token.ToTokenData(), "startIndex out of range. got=%d, expected=%d-%d",
					startIndex, 0, length-1)
			}
			if !IsIndexInRange(endIndex, length+1) {
				return NewError(token.ToTokenData(), "endIndex out of range. got=%d, expected=%d-%d",
					endIndex, 0, length-1)
			}

			if startIndex > endIndex {
				return NewError(token.ToTokenData(), "startIndex larger than endIndex. startIndex=%d, endIndex=%d",
					startIndex, endIndex)
			}

			return &object.Array{
				Elements: arrayObj.Elements[startIndex:endIndex],
			}
		}

	default:
		return NewError(token.ToTokenData(), "parser probably failed, this should never happen")
	}
}

// Create a function and eval it
func ApplyFunction(token token.Token, function object.Object, args []object.Object) object.Object {
	switch fn := function.(type) {
	case *object.Function:
		extendedEnv := ExtendFunctionEnv(fn, args)
		evaluated := Eval(fn.Body, extendedEnv)
		return UnwrapReturnValue(evaluated)

	case *object.Builtin:
		return fn.Fn(token, args...)

	default:
		return NewError(token.ToTokenData(), "not a function: %s", function.Type())

	}
}

// Unwrap the return value for an object
func UnwrapReturnValue(evaluated object.Object) object.Object {
	if returnValue, ok := evaluated.(*object.ReturnValue); ok {
		return returnValue.Value
	}

	return evaluated
}

// Create the environment for the function
func ExtendFunctionEnv(function *object.Function, args []object.Object) *object.Environment {
	env := object.NewEnclosingEnvironment(function.Env)

	// Map the arguments to parameters
	for paramIdx, param := range function.Parameters {
		env.Store(param.Value, args[paramIdx])
	}

	return env
}

// Eval a list of expressions
func EvalExpressions(arguments []ast.Expression, env *object.Environment) []object.Object {
	var result []object.Object

	for _, e := range arguments {
		evaluated := Eval(e, env)
		if CheckError(evaluated) {
			return []object.Object{evaluated}
		}
		result = append(result, evaluated)
	}

	return result
}

// Fetch the value from env and return it
func EvalIdentifier(node *ast.Identifier, env *object.Environment) object.Object {
	if val, ok := env.Get(node.Value); ok {
		return val
	}

	if builtin, ok := builtins[node.Value]; ok {
		return builtin
	}

	return NewError(node.Token.ToTokenData(), "identifier not found: %s", node.Value)
}

// Eval a block statement
func EvalBlockStatement(block *ast.BlockStatement, env *object.Environment) object.Object {
	var result object.Object

	for _, statement := range block.Statements {
		result = Eval(statement, env)

		if result != nil {
			rt := result.Type()
			if rt == object.RETURN_VALUE_OBJ || rt == object.ERROR_OBJ {
				return result
			}
		}
	}

	return result
}

// Eval If Expression
func EvalIfExpression(ie *ast.IfExpression, env *object.Environment) object.Object {
	condition := Eval(ie.Condition, env)
	if CheckError(condition) {
		return condition
	}

	if IsTruthful(condition) {
		return Eval(ie.Consequence, env)
	} else if ie.Alternative != nil {
		return Eval(ie.Alternative, env)
	} else {
		return NULL
	}
}

func UnhandledOperationError(token token.Token, left object.Object, right object.Object, operator string) object.Object {
	return NewError(token.ToTokenData(), "unknown operation: %s %s %s",
		left.Type(), operator, right.Type())
}

func EvalOperatorExpression(token token.Token, operator string, left object.Object, right object.Object) object.Object {
	if fn, ok := InfixMap[left.Type()][operator]; ok {
		result := fn(token, left, right)
		if result != nil {
			return result
		}
	}

	// Implicit Handling
	switch {
	case operator == "==":
		return NativeBoolToBooleanObject(left == right)
	case operator == "!=":
		return NativeBoolToBooleanObject(left != right)
	case left.Type() != right.Type():
		return NewError(token.ToTokenData(), "type mismatch: %s %s %s",
			left.Type(), operator, right.Type())
	default:
		return UnhandledOperationError(token, left, right, operator)
	}
}

// Eval Infix Expression
func EvalInfixExpression(node *ast.InfixExpression, env *object.Environment) object.Object {
	operator := node.Operator
	left := Eval(node.Left, env)
	if CheckError(left) {
		return left
	}

	// Eval Short Circuits
	switch operator {
	case token.AND, token.OR, token.XOR:
		return EvalShortCircuitExpression(operator, left, node, env)
	}

	// Regular
	right := Eval(node.Right, env)
	if CheckError(right) {
		return right
	}

	return EvalOperatorExpression(node.Token, operator, left, right)
	//else if fn, ok = InfixMap[right.Type()][operator]; ok {
	//	return fn(node.Token, left, right)
	//}

	//switch {
	//case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
	//	return EvalIntegerInfixExpression(operator, left, right, node.Token)
	//
	//case left.Type() == object.STRING_OBJ || right.Type() == object.STRING_OBJ:
	//	return EvalStringInfixExpression(operator, left, right, node.Token)
	//
	//
	//case operator == "==":
	//	return NativeBoolToBooleanObject(IsTruthful(left) == IsTruthful(right))
	//case operator == "!=":
	//	return NativeBoolToBooleanObject(IsTruthful(left) != IsTruthful(right))
	//
	//case left.Type() != right.Type():
	//	return NewError(node.Token.ToTokenData(), "type mismatch: %s %s %s",
	//		left.Type(), operator, right.Type())
	//default:
	//	return NewError(node.Token.ToTokenData(), "unknown operator: %s %s %s",
	//		left.Type(), operator, right.Type())
	//}
}

func EvalShortCircuitExpression(operator string, left object.Object, node *ast.InfixExpression, env *object.Environment) object.Object {
	switch operator {
	case token.AND:
		return EvalAndExpression(left, node.Right, env)
	case token.OR:
		return EvalOrExpression(left, node.Right, env)
	default:
		right := Eval(node.Right, env)
		if CheckError(right) {
			return right
		}
		return NativeBoolToBooleanObject(IsTruthful(left) != IsTruthful(right))
	}
}

// All string operators
// Legacy
func EvalStringInfixExpression(operator string, left object.Object, right object.Object, token token.Token) object.Object {
	switch operator {
	case "+":
		{
			// Left is int, right is string
			if leftVal, ok := left.(*object.String); ok {
				if rightVal, ok := right.(*object.String); ok {
					return &object.String{
						Value: leftVal.Value + rightVal.Value,
					}
				} else if rightVal, ok := right.(*object.Integer); ok {
					// Left is string, right is int
					leftVal := left.(*object.String).Value
					return &object.String{
						Value: leftVal + parser.FormatFloat(rightVal.Value),
					}
				}
			} else if leftVal, ok := left.(*object.Integer); ok {
				rightVal := right.(*object.String).Value
				return &object.String{
					Value: parser.FormatFloat(leftVal.Value) + rightVal,
				}
			} else {
				// Error
			}
		}
	case "*":
		{
			if rightVal, ok := right.(*object.Integer); ok {
				leftVal := left.(*object.String)
				var out strings.Builder

				amount := int(rightVal.Value)

				for i := 0; i < amount; i++ {
					out.WriteString(leftVal.Value)
				}

				return &object.String{Value: out.String()}
			}
		}
	case "==":
		{
			if leftVal, ok := left.(*object.String); ok {
				if rightVal, ok := right.(*object.String); ok {
					return NativeBoolToBooleanObject(leftVal.Value == rightVal.Value)
				}
			}
		}
	case "!=":
		{
			if leftVal, ok := left.(*object.String); ok {
				if rightVal, ok := right.(*object.String); ok {
					return NativeBoolToBooleanObject(leftVal.Value != rightVal.Value)
				}
			}
		}

	}

	return NewError(token.ToTokenData(), "unknown operator: %s %s %s",
		left.Type(), operator, right.Type())
}

// Eval Or Expression
func EvalOrExpression(left object.Object, rightExp ast.Expression, env *object.Environment) object.Object {
	// If True
	if IsTruthful(left) {
		// Short circuit
		return NativeBoolToBooleanObject(true)
	}
	right := Eval(rightExp, env)
	return NativeBoolToBooleanObject(IsTruthful(right))
}

// Eval And Expression
func EvalAndExpression(left object.Object, rightExp ast.Expression, env *object.Environment) object.Object {
	// If false
	if !IsTruthful(left) {
		// Short circuit
		return NativeBoolToBooleanObject(false)
	}
	right := Eval(rightExp, env)
	return NativeBoolToBooleanObject(IsTruthful(right))
}

// Eval Integer Expression
// Legacy
func EvalIntegerInfixExpression(operator string, left object.Object, right object.Object, token token.Token) object.Object {
	leftVal := left.(*object.Integer).Value
	rightVal := right.(*object.Integer).Value

	switch operator {
	case "+":
		return &object.Integer{Value: leftVal + rightVal}
	case "-":
		return &object.Integer{Value: leftVal - rightVal}
	case "*":
		return &object.Integer{Value: leftVal * rightVal}
	case "/":
		return &object.Integer{Value: leftVal / rightVal}
	case "%":
		return &object.Integer{Value: float64(int64(leftVal) % int64(rightVal))}
	case "<":
		return NativeBoolToBooleanObject(leftVal < rightVal)
	case "<=":
		return NativeBoolToBooleanObject(leftVal <= rightVal)
	case ">":
		return NativeBoolToBooleanObject(leftVal > rightVal)
	case ">=":
		return NativeBoolToBooleanObject(leftVal >= rightVal)
	case "==":
		return NativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return NativeBoolToBooleanObject(leftVal != rightVal)
	default:
		return NewError(token.ToTokenData(), "unknown operator: %s %s %s",
			left.Type(), operator, right.Type())
	}
}

// Eval a prefix expression
func EvalPrefixExpression(operator string, right object.Object, token token.Token) object.Object {
	switch operator {
	case "!":
		return EvalBangOperatorExpression(right)
	case "-":
		return EvalMinusPrefixOperatorExpression(right, token)
	case "+":
		return EvalPlusPrefixOperatorExpression(right, token)
	default:
		return NewError(token.ToTokenData(), "unknown operation: %s%s",
			operator, right.Type())
	}
}

// Eval + infix operator
func EvalPlusPrefixOperatorExpression(right object.Object, token token.Token) object.Object {
	if right.Type() != object.INTEGER_OBJ {
		return NewError(token.ToTokenData(), "unknown operation: +%s", right.Type())
	}

	return right
}

// Eval - infix operator
func EvalMinusPrefixOperatorExpression(right object.Object, token token.Token) object.Object {
	// Not Integer
	if right.Type() != object.INTEGER_OBJ {
		return NewError(token.ToTokenData(), "unknown operation: -%s", right.Type())
	}

	value := right.(*object.Integer).Value
	return &object.Integer{Value: -value}
}

// Eval the bang/invert operator
func EvalBangOperatorExpression(right object.Object) object.Object {
	isTrue := IsTruthful(right)
	return NativeBoolToBooleanObject(!isTrue)
}

// Converter a native type to boolean object
func NativeBoolToBooleanObject(value bool) object.Object {
	if value {
		return TRUE
	}
	return FALSE
}

// Eval Print ExpressionStmt
func EvalPrintExpressionStatement(token token.Token, exp ast.Expression, env *object.Environment) object.Object {
	result := Eval(exp, env)
	builtins["writeLine"].Fn(token, result)
	//fmt.Println(result.Inspect())
	return NULL
}

// Init Infix Map
func InitInfix() {
	arrayInit()

	InfixMap = map[object.ObjectType]InfixObj{
		object.INTEGER_OBJ: Integer,
		object.STRING_OBJ:  String,
		object.ARRAY_OBJ:   Array,
	}
}

// Eval Statements
func EvalProgram(program *ast.Program, env *object.Environment) object.Object {
	var result object.Object

	for _, statement := range program.Statements {
		result = Eval(statement, env)

		switch result.(type) {
		case *object.ReturnValue:
			return result.(*object.ReturnValue).Value
		case *object.Error:
			return result
		}
	}

	return result
}
