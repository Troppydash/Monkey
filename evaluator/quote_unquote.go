package evaluator

import (
	"Monkey/ast"
	"Monkey/object"
	"Monkey/parser"
	"Monkey/token"
)

func EvalUnquote(token token.Token, arguments []ast.Expression, environment *object.Environment) object.Object {
	if len(arguments) < 1 {
		return NewFatalError(token.ToTokenData(), "unquote only takes one argument. got=%d", len(arguments))
	}

	argument := Eval(arguments[0], environment)

	return argument
}

func EvalQuote(token token.Token, arguments []ast.Expression, environment *object.Environment) object.Object {
	if len(arguments) != 1 {
		return NewFatalError(token.ToTokenData(), "quote only takes one argument. got=%d", len(arguments))
	}

	argument := evalUnquoteCalls(arguments[0], environment)

	return &object.Quote{
		Node: argument,
	}
}

func evalUnquoteCalls(quoted ast.Node, environment *object.Environment) ast.Node {
	return ast.Modify(quoted, func(node ast.Node) ast.Node {
		if !isUnquoteCall(node) {
			return node
		}
		call, ok := node.(*ast.CallExpression)
		if !ok {
			return node
		}

		if len(call.Arguments) != 1 {
			return node
		}

		unquoted := Eval(call.Arguments[0], environment)
		return convertObjectToASTNode(unquoted, call.Token)
	})
}

func convertObjectToASTNode(unquoted object.Object, t token.Token) ast.Node {
	switch obj := unquoted.(type) {
	case *object.Integer:
		tmpt := token.NewToken(
			token.INT,
			parser.FormatFloat(obj.Value),
			t.ToTokenData(),
		)
		//t := token.Token{
		//	Type:         token.INT,
		//	Literal:      parser.FormatFloat(obj.Value),
		//	RowNumber:    t.RowNumber,
		//	ColumnNumber: t.ColumnNumber,
		//	Filename:     t.Filename,
		//}
		return &ast.IntegerLiteral{Token: tmpt, Value: obj.Value}
	case *object.Boolean:
		var tmpt token.Token
		if obj.Value {
			tmpt = token.NewToken(token.TRUE, "true", t.ToTokenData())
		} else {
			tmpt = token.NewToken(token.FALSE, "false", t.ToTokenData())
		}
		return &ast.Boolean{Token: tmpt, Value: obj.Value}
	case *object.Quote:
		return obj.Node
	default:
		return &ast.Null{Token: t}
	}
}

func isUnquoteCall(node ast.Node) bool {
	callExpression, ok := node.(*ast.CallExpression)
	if !ok {
		return false
	}

	return callExpression.Function.TokenLiteral() == "unquote"
}
