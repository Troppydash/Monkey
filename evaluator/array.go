package evaluator

import (
	"Monkey/object"
	"Monkey/token"
)

var Array map[string]InfixFn

func handleArray(operator string, token token.Token, left object.Object, right object.Object) object.Object {
	leftObj := left.(*object.Array)
	leftVal := leftObj.Elements

	switch right.(type) {
	case *object.Array:
		rightVal := right.(*object.Array).Elements

		for i := 0; i < len(leftVal); i++ {
			if i >= len(rightVal) {
				break
			}
			left := leftVal[i]
			right := rightVal[i]

			result := EvalOperatorExpression(token, operator, left, right)
			leftVal[i] = result
		}
		return leftObj
	case *object.Integer:
		rightVal := right.(*object.Integer)

		for i := 0; i < len(leftVal); i++ {
			left := leftVal[i]
			result := EvalOperatorExpression(token, operator, left, rightVal)
			leftVal[i] = result
		}
		return leftObj
	default:
		return nil
	}
}

// Array is special
func init() {
	Array = map[string]InfixFn{
		"+": func(token token.Token, left object.Object, right object.Object) object.Object {
			return handleArray("+", token, left, right)
		},
		"-": func(token token.Token, left object.Object, right object.Object) object.Object {
			return handleArray("-", token, left, right)
		},
		"*": func(token token.Token, left object.Object, right object.Object) object.Object {
			return handleArray("*", token, left, right)
		},
		"/": func(token token.Token, left object.Object, right object.Object) object.Object {
			return handleArray("/", token, left, right)
		},
		"%": func(token token.Token, left object.Object, right object.Object) object.Object {
			return handleArray("%", token, left, right)
		},
	}
}
