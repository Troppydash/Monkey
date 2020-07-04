package evaluator

import (
	"Monkey/object"
	"Monkey/token"
)

var Integer = map[string]InfixFn{
	"+": func(token token.Token, left object.Object, right object.Object) object.Object {
		leftVal := left.(*object.Integer)
		switch right.(type) {
		case *object.Integer:
			rightVal := right.(*object.Integer)
			return &object.Integer{Value: leftVal.Value + rightVal.Value}
		default:
			return nil
		}
	},
	"-": func(token token.Token, left object.Object, right object.Object) object.Object {
		leftVal := left.(*object.Integer)
		switch right.(type) {
		case *object.Integer:
			rightVal := right.(*object.Integer)
			return &object.Integer{Value: leftVal.Value - rightVal.Value}
		default:
			return nil
		}
	},
	"*": func(token token.Token, left object.Object, right object.Object) object.Object {
		leftVal := left.(*object.Integer)
		switch right.(type) {
		case *object.Integer:
			rightVal := right.(*object.Integer)
			return &object.Integer{Value: leftVal.Value * rightVal.Value}
		default:
			return nil
		}
	},
	"/": func(token token.Token, left object.Object, right object.Object) object.Object {
		leftVal := left.(*object.Integer)
		switch right.(type) {
		case *object.Integer:
			rightVal := right.(*object.Integer)
			return &object.Integer{Value: leftVal.Value / rightVal.Value}
		default:
			return nil
		}
	},
	"%": func(token token.Token, left object.Object, right object.Object) object.Object {
		leftVal := left.(*object.Integer)
		switch right.(type) {
		case *object.Integer:
			rightVal := right.(*object.Integer)
			return &object.Integer{Value: float64(int64(leftVal.Value) % int64(rightVal.Value))}
		default:
			return nil
		}
	},
	"<": func(token token.Token, left object.Object, right object.Object) object.Object {
		leftVal := left.(*object.Integer).Value
		switch right.(type) {
		case *object.Integer:
			rightVal := right.(*object.Integer).Value
			return NativeBoolToBooleanObject(leftVal < rightVal)
		default:
			return nil
		}
	},
	"<=": func(token token.Token, left object.Object, right object.Object) object.Object {
		leftVal := left.(*object.Integer).Value
		switch right.(type) {
		case *object.Integer:
			rightVal := right.(*object.Integer).Value
			return NativeBoolToBooleanObject(leftVal <= rightVal)
		default:
			return nil
		}
	},
	">": func(token token.Token, left object.Object, right object.Object) object.Object {
		leftVal := left.(*object.Integer).Value
		switch right.(type) {
		case *object.Integer:
			rightVal := right.(*object.Integer).Value
			return NativeBoolToBooleanObject(leftVal > rightVal)
		default:
			return nil
		}
	},
	">=": func(token token.Token, left object.Object, right object.Object) object.Object {
		leftVal := left.(*object.Integer).Value
		switch right.(type) {
		case *object.Integer:
			rightVal := right.(*object.Integer).Value
			return NativeBoolToBooleanObject(leftVal >= rightVal)
		default:
			return nil
		}
	},
	"==": func(token token.Token, left object.Object, right object.Object) object.Object {
		leftVal := left.(*object.Integer).Value
		switch right.(type) {
		case *object.Integer:
			rightVal := right.(*object.Integer).Value
			return NativeBoolToBooleanObject(leftVal == rightVal)
		default:
			return nil
		}
	},
	"!=": func(token token.Token, left object.Object, right object.Object) object.Object {
		leftVal := left.(*object.Integer).Value
		switch right.(type) {
		case *object.Integer:
			rightVal := right.(*object.Integer).Value
			return NativeBoolToBooleanObject(leftVal != rightVal)
		default:
			return nil
		}
	},
}
