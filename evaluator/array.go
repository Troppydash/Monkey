package evaluator

import (
	"Monkey/object"
	"Monkey/token"
)

var Array = map[string]InfixFn{
	"+": func(token token.Token, left object.Object, right object.Object) object.Object {
		//leftVal := left.(*object.Array)
		//switch right.(type) {
		//case *object.Array:
		//	rightVal := right.(*object.Array)
		//	return NULL
		//default:
		//	return nil
		//}
		return nil
	},
}
