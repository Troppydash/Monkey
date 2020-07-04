package evaluator

import (
	"Monkey/object"
	"Monkey/token"
	"strings"
)

var String = map[string]InfixFn{
	"+": func(token token.Token, left object.Object, right object.Object) object.Object {
		leftVal := left.(*object.String).Value
		switch right.(type) {
		case *object.String:
			rightVal := right.(*object.String).Value
			return &object.String{Value: leftVal + rightVal}
		default:
			return nil
		}
	},
	"*": func(token token.Token, left object.Object, right object.Object) object.Object {
		leftVal := left.(*object.String).Value
		switch right.(type) {
		case *object.Integer:
			rightVal := int(right.(*object.Integer).Value)

			var out strings.Builder

			for i := 0; i < rightVal; i++ {
				out.WriteString(leftVal)
			}
			return &object.String{Value: out.String()}
		default:
			return nil
		}
	},
	"==": func(token token.Token, left object.Object, right object.Object) object.Object {
		leftVal := left.(*object.String).Value
		switch right.(type) {
		case *object.String:
			rightVal := right.(*object.String).Value
			return NativeBoolToBooleanObject(leftVal == rightVal)

		default:
			return nil
		}
	},
	"!=": func(token token.Token, left object.Object, right object.Object) object.Object {
		leftVal := left.(*object.String).Value
		switch right.(type) {
		case *object.String:
			rightVal := right.(*object.String).Value
			return NativeBoolToBooleanObject(leftVal != rightVal)

		default:
			return nil
		}
	},
}
