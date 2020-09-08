package ast

type ModifierFunc func(Node) Node

// Modify walks the ast tree and apply the modifier func to each node
func Modify(node Node, modifier ModifierFunc) Node {
	switch node := node.(type) {
	case *Program:
		for i, statement := range node.Statements {
			node.Statements[i], _ = Modify(statement, modifier).(Statement)
		}
	case *ExpressionStatement:
		node.Expression, _ = Modify(node.Expression, modifier).(Expression)

	case *InfixExpression:
		node.Left, _ = Modify(node.Left, modifier).(Expression)
		node.Right, _ = Modify(node.Right, modifier).(Expression)

	case *PrefixExpression:
		node.Right, _ = Modify(node.Right, modifier).(Expression)

	case *IndexExpression:
		node.Left, _ = Modify(node.Left, modifier).(Expression)
		node.Start, _ = Modify(node.Start, modifier).(Expression)
		node.End, _ = Modify(node.End, modifier).(Expression)

	case *IfExpression:
		node.Condition, _ = Modify(node.Condition, modifier).(Expression)
		node.Consequence, _ = Modify(node.Consequence, modifier).(*BlockStatement)
		if node.Alternative != nil {
			node.Alternative, _ = Modify(node.Alternative, modifier).(*BlockStatement)
		}
	case *BlockStatement:
		for i, _ := range node.Statements {
			node.Statements[i], _ = Modify(node.Statements[i], modifier).(Statement)
		}
	case *ReturnStatement:
		node.ReturnValue, _ = Modify(node.ReturnValue, modifier).(Expression)

	case *LetStatement:
		node.Value, _ = Modify(node.Value, modifier).(Expression)

	case *FunctionLiteral:
		for i, _ := range node.Parameters {
			node.Parameters[i], _ = Modify(node.Parameters[i], modifier).(*Identifier)
		}
		node.Body, _ = Modify(node.Body, modifier).(*BlockStatement)

	case *ArrayLiteral:
		for i, _ := range node.Elements {
			node.Elements[i], _ = Modify(node.Elements[i], modifier).(Expression)
		}

	case *HashLiteral:
		newPairs := make(map[Expression]Expression)
		for key, val := range node.Pairs {
			newKey, _ := Modify(key, modifier).(Expression)
			newValue, _ := Modify(val, modifier).(Expression)
			newPairs[newKey] = newValue
		}
		node.Pairs = newPairs
	case *CallExpression:
		node.Function, _ = Modify(node.Function, modifier).(Expression)
		for i, _ := range node.Arguments {
			node.Arguments[i], _ = Modify(node.Arguments[i], modifier).(Expression)
		}

	}
	return modifier(node)
}
