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
func Eval(node ast.Node) object.Object {
	switch node := node.(type) {

	case *ast.Program:
		return EvalProgram(node)

	case *ast.ExpressionStatement:
		return Eval(node.Expression)

	case *ast.PrintExpressionStatement:
		return EvalPrintExpressionStatement(node.Expression)

	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}

	case *ast.Boolean:
		return NativeBoolToBooleanObject(node.Value)

	case *ast.PrefixExpression:
		right := Eval(node.Right)
		return EvalPrefixExpression(node.Operator, right)

	case *ast.InfixExpression:
		// We need to short circuit AND or OR gates
		return EvalInfixExpression(node)

	case *ast.BlockStatement:
		return EvalBlockStatement(node)

	case *ast.IfExpression:
		return EvalIfExpression(node)

	case *ast.ReturnStatement:
		val := Eval(node.ReturnValue)
		return &object.ReturnValue{Value: val}
	}

	return NULL
}

// Eval a block statement
func EvalBlockStatement(block *ast.BlockStatement) object.Object {
	var result object.Object

	for _, statement := range block.Statements {
		result = Eval(statement)

		if result != nil && result.Type() == object.RETURN_VALUE_OBJ {
			return result
		}
	}

	return result
}

// Eval If Expression
func EvalIfExpression(ie *ast.IfExpression) object.Object {
	condition := Eval(ie.Condition)

	if IsTruthful(condition) {
		return Eval(ie.Consequence)
	} else if ie.Alternative != nil {
		return Eval(ie.Alternative)
	} else {
		return NULL
	}
}

// Eval Infix Expression
func EvalInfixExpression(node *ast.InfixExpression) object.Object {
	operator := node.Operator
	left := Eval(node.Left)

	switch operator {
	case token.AND:
		return EvalAndExpression(left, node.Right)
	case token.OR:
		return EvalOrExpression(left, node.Right)
	case token.XOR:
		right := Eval(node.Right)
		return NativeBoolToBooleanObject(IsTruthful(left) != IsTruthful(right))
	}

	right := Eval(node.Right)

	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return EvalIntegerInfixExpression(operator, left, right)
	case operator == "==":
		return NativeBoolToBooleanObject(IsTruthful(left) == IsTruthful(right))
	case operator == "!=":
		return NativeBoolToBooleanObject(IsTruthful(left) != IsTruthful(right))
	default:
		return NULL
	}
}

// Eval Or Expression
func EvalOrExpression(left object.Object, rightExp ast.Expression) object.Object {
	// If True
	if IsTruthful(left) {
		// Short circuit
		return NativeBoolToBooleanObject(true)
	}
	right := Eval(rightExp)
	return NativeBoolToBooleanObject(IsTruthful(right))
}

// Eval And Expression
func EvalAndExpression(left object.Object, rightExp ast.Expression) object.Object {
	// If false
	if !IsTruthful(left) {
		// Short circuit
		return NativeBoolToBooleanObject(false)
	}
	right := Eval(rightExp)
	return NativeBoolToBooleanObject(IsTruthful(right))
}

// Eval Integer Expression
func EvalIntegerInfixExpression(operator string, left object.Object, right object.Object) object.Object {
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
		return NULL
	}
}

// Eval a prefix expression
func EvalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return EvalBangOperatorExpression(right)
	case "-":
		return EvalMinusPrefixOperatorExpression(right)
	case "+":
		return EvalPlusPrefixOperatorExpression(right)
	default:
		return NULL
	}
}

// Eval + infix operator
func EvalPlusPrefixOperatorExpression(right object.Object) object.Object {
	if right.Type() != object.INTEGER_OBJ {
		return NULL
	}

	return right
}

// Eval - infix operator
func EvalMinusPrefixOperatorExpression(right object.Object) object.Object {
	// Not Integer
	if right.Type() != object.INTEGER_OBJ {
		return NULL
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
func EvalPrintExpressionStatement(exp ast.Expression) object.Object {
	result := Eval(exp)
	fmt.Println(result.Inspect())
	return result
}

// Eval Statements
func EvalProgram(program *ast.Program) object.Object {
	var result object.Object

	for _, statement := range program.Statements {
		result = Eval(statement)

		switch result.(type) {
		case *object.ReturnValue:
			return result.(*object.ReturnValue).Value
		}
	}

	return result
}
