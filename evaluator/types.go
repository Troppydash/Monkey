package evaluator

import (
	"Monkey/object"
	"Monkey/token"
)

type InfixFn func(token token.Token, left object.Object, right object.Object) object.Object

type InfixObj map[string]InfixFn

var InfixMap map[object.ObjectType]InfixObj
