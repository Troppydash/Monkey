package evaluator

import (
	"Monkey/object"
	"Monkey/token"
)

var prototypes map[object.ObjectType]*object.Hash

func init() {
	double := &object.String{
		Value: "double",
	}
	prototypes = map[object.ObjectType]*object.Hash{
		object.IntegerObj: {
			Pairs: map[object.HashKey]object.HashPair{
				double.HashKey(): object.HashPair{
					Key: double,
					Value: &object.Builtin{
						Fn: func(token token.Token, env *object.Environment, args ...object.Object) object.Object {
							self, _ := env.Get("this")
							integer, _ := self.(*object.Integer)
							integer.Value *= 2
							return integer
						},
						Parameters: 0,
						VarArgs:    false,
						Prototype:  true,
					},
				},
			},
		},
	}
}
