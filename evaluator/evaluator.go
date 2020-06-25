package evaluator

import (
	"Monkey/ast"
	"Monkey/object"
	"fmt"
)

var (
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
	NULL  = &object.Null{}
)

// Master Eval Function
func Eval(node ast.Node) object.Object {
	switch node := node.(type) {

	case *ast.Program:
		return EvalStatements(node.Statements)

	case *ast.ExpressionStatement:
		return Eval(node.Expression)

	case *ast.PrintExpressionStatement:
		return EvalPrintExpressionStatement(node.ExpStmt)

	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}

	case *ast.Boolean:
		return NativeBoolToBooleanObject(node.Value)
	}

	return nil
}

func NativeBoolToBooleanObject(value bool) object.Object {
	if value {
		return TRUE
	}
	return FALSE
}

// Eval Print ExpressionStmt
func EvalPrintExpressionStatement(expStmt *ast.ExpressionStatement) object.Object {
	result := Eval(expStmt)
	fmt.Println(result.Inspect())
	return result
}

// Eval Statements
func EvalStatements(stmts []ast.Statement) object.Object {
	var result object.Object

	for _, statement := range stmts {
		result = Eval(statement)
	}

	return result
}
