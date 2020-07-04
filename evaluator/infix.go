package evaluator

import (
	"Monkey/object"
)

var InfixMap = map[object.ObjectType]InfixObj{
	object.INTEGER_OBJ: Integer,
	object.STRING_OBJ:  String,
	object.ARRAY_OBJ:   Array,
}
