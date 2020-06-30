package evaluator

import (
	"Monkey/ast"
	"Monkey/object"
	"Monkey/token"
	"fmt"
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

// Is the object an error
func IsError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJ
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
		return EvalPrintExpressionStatement(node.Expression, env)

	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}

	case *ast.Boolean:
		return NativeBoolToBooleanObject(node.Value)

	case *ast.PrefixExpression:
		right := Eval(node.Right, env)
		if IsError(right) {
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
		if IsError(val) {
			return val
		}
		return &object.ReturnValue{Value: val}

	case *ast.LetStatement:
		val := Eval(node.Value, env)
		if IsError(val) {
			return val
		}

		env.Store(node.Name.Value, val)
	// Uncomment for let to return a value
	//return val

	case *ast.Identifier:
		return EvalIdentifier(node, env)

	case *ast.FunctionLiteral:
		params := node.Parameters
		body := node.Body
		return &object.Function{Parameters: params, Env: env, Body: body}

	case *ast.CallExpression:
		function := Eval(node.Function, env)
		if IsError(function) {
			return function
		}

		args := EvalExpression(node.Arguments, env)
		if len(args) == 1 && IsError(args[0]) {
			return args[0]
		}
		return ApplyFunction(node.Token, function, args)

	}

	return NULL
}

// Create a function and eval it
func ApplyFunction(token token.Token, function object.Object, args []object.Object) object.Object {
	fn, ok := function.(*object.Function)
	if !ok {
		return NewError(token.ToTokenData(), "not a function: %s", function.Type())
	}

	extendedEnv := ExtendFunctionEnv(fn, args)
	evaluated := Eval(fn.Body, extendedEnv)
	return UnwrapReturnValue(evaluated)
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
func EvalExpression(arguments []ast.Expression, env *object.Environment) []object.Object {
	var result []object.Object

	for _, e := range arguments {
		evaluated := Eval(e, env)
		if IsError(evaluated) {
			return []object.Object{evaluated}
		}
		result = append(result, evaluated)
	}

	return result
}

// Fetch the value from env
func EvalIdentifier(node *ast.Identifier, env *object.Environment) object.Object {
	val, ok := env.Get(node.Value)
	if !ok {
		return NewError(node.Token.ToTokenData(), "identifier not found: %s", node.Value)
	}
	return val
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
	if IsError(condition) {
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

// Eval Infix Expression
func EvalInfixExpression(node *ast.InfixExpression, env *object.Environment) object.Object {
	operator := node.Operator
	left := Eval(node.Left, env)

	if IsError(left) {
		return left
	}

	switch operator {
	case token.AND:
		return EvalAndExpression(left, node.Right, env)
	case token.OR:
		return EvalOrExpression(left, node.Right, env)
	case token.XOR:
		right := Eval(node.Right, env)
		if IsError(right) {
			return right
		}
		return NativeBoolToBooleanObject(IsTruthful(left) != IsTruthful(right))
	}

	right := Eval(node.Right, env)
	if IsError(right) {
		return right
	}

	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return EvalIntegerInfixExpression(operator, left, right, node.Token)
	case operator == "==":
		return NativeBoolToBooleanObject(IsTruthful(left) == IsTruthful(right))
	case operator == "!=":
		return NativeBoolToBooleanObject(IsTruthful(left) != IsTruthful(right))

	case left.Type() != right.Type():
		return NewError(node.Token.ToTokenData(), "type mismatch: %s %s %s",
			left.Type(), operator, right.Type())
	default:
		return NewError(node.Token.ToTokenData(), "unknown operator: %s %s %s",
			left.Type(), operator, right.Type())
	}
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
		return NewError(token.ToTokenData(), "unknown operator: %s%s",
			operator, right.Type())
	}
}

// Eval + infix operator
func EvalPlusPrefixOperatorExpression(right object.Object, token token.Token) object.Object {
	if right.Type() != object.INTEGER_OBJ {
		return NewError(token.ToTokenData(), "unknown operator: +%s", right.Type())
	}

	return right
}

// Eval - infix operator
func EvalMinusPrefixOperatorExpression(right object.Object, token token.Token) object.Object {
	// Not Integer
	if right.Type() != object.INTEGER_OBJ {
		return NewError(token.ToTokenData(), "unknown operator: -%s", right.Type())
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
func EvalPrintExpressionStatement(exp ast.Expression, env *object.Environment) object.Object {
	result := Eval(exp, env)
	fmt.Println(result.Inspect())
	return result
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
