package evaluator

import (
	"Monkey/object"
	"Monkey/token"
)

var prototypes map[object.ObjectType]*object.Hash

type keyHashKeyPair struct {
	keys map[string]*object.String
}

func NewKHKP() *keyHashKeyPair {
	return &keyHashKeyPair{
		keys: make(map[string]*object.String),
	}
}
func (kh *keyHashKeyPair) AddKey(key string) {
	kh.keys[key] = &object.String{Value: key}
}
func (kh *keyHashKeyPair) GetKey(key string) *object.String {
	result, _ := kh.keys[key]
	return result
}
func (kh *keyHashKeyPair) GetKeyHash(key string) object.HashKey {
	return kh.GetKey(key).HashKey()
}

func init() {
	khkp := NewKHKP()
	khkp.AddKey("double")
	khkp.AddKey("length")
	khkp.AddKey("keys")
	khkp.AddKey("values")
	khkp.AddKey("push")
	khkp.AddKey("pop")
	prototypes = map[object.ObjectType]*object.Hash{
		object.IntegerObj: {
			Pairs: map[object.HashKey]object.HashPair{
				khkp.GetKeyHash("double"): {
					Key: khkp.GetKey("double"),
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
						Eval:       true,
					},
				},
			},
		},
		object.StringObj: {
			Pairs: map[object.HashKey]object.HashPair{
				khkp.GetKeyHash("length"): {
					Key: khkp.GetKey("length"),
					Value: &object.Builtin{
						Fn: func(token token.Token, env *object.Environment, args ...object.Object) object.Object {
							self, _ := env.Get("this")
							str, _ := self.(*object.String)
							return &object.Integer{Value: float64(len(str.Value))}
						},
						Parameters: 0,
						VarArgs:    false,
						Prototype:  true,
						Eval:       true,
					},
				},
			},
		},
		object.ArrayObj: {
			Pairs: map[object.HashKey]object.HashPair{
				khkp.GetKeyHash("length"): {
					Key: khkp.GetKey("length"),
					Value: &object.Builtin{
						Fn: func(token token.Token, env *object.Environment, args ...object.Object) object.Object {
							self, _ := env.Get("this")
							array, _ := self.(*object.Array)
							return &object.Integer{Value: float64(len(array.Elements))}
						},
						Parameters: 0,
						VarArgs:    false,
						Prototype:  true,
						Eval:       true,
					},
				},
				khkp.GetKeyHash("pop"): {
					Key: khkp.GetKey("pop"),
					Value: &object.Builtin{
						Fn: func(token token.Token, env *object.Environment, args ...object.Object) object.Object {
							self, _ := env.Get("this")
							array, _ := self.(*object.Array)

							parser := object.NewParser(
								object.NewOptionalArgument(func(optional bool, arg object.Object) object.ParsedArgument {
									if optional {
										return object.ParsedArgument{
											Value: 1,
										}
									}
									integer, _ := arg.(*object.Integer)
									return object.ParsedArgument{
										Value: int(integer.Value),
									}
								}, object.IntegerObj),
							)
							result, err := parser.Parse(args)
							if err != nil {
								return NewFatalError(token.ToTokenData(), err.Error())
							}

							amount := result[0].Value.(int)

							length := len(array.Elements)
							if amount > length {
								return NewFatalError(token.ToTokenData(), "array index out of bounds")
							}
							newElements := make([]object.Object, length-amount, length-amount)
							copy(newElements, array.Elements[:length-amount])
							oldElements := array.Elements[length-amount:]
							array.Elements = newElements

							if len(args) == 0 {
								return oldElements[0]
							}

							return &object.Array{
								Elements: oldElements,
							}
						},
						Parameters: 1,
						VarArgs:    false,
						Prototype:  true,
						Eval:       false,
					},
				},
				khkp.GetKeyHash("push"): {
					Key: khkp.GetKey("push"),
					Value: &object.Builtin{
						Fn: func(token token.Token, env *object.Environment, args ...object.Object) object.Object {
							self, _ := env.Get("this")
							array, _ := self.(*object.Array)

							parser := object.NewParser(
								object.NewAnyOptionalArgument(
									func(optional bool, arg object.Object) object.ParsedArgument {
										if optional {
											return object.ParsedArgument{
												Value: &object.String{Value: ""},
											}
										} else {
											return object.ParsedArgument{
												Value: arg,
											}
										}
									}),
								object.NewAnyVarargsArgument(
									func(optional bool, arg object.Object) object.ParsedArgument {
										return object.ParsedArgument{Value: arg}
									}),
							)
							result, err := parser.Parse(args)
							if err != nil {
								return NewFatalError(token.ToTokenData(), err.Error())
							}

							length := len(array.Elements)
							var arrayToMerge []object.Object
							for _, res := range result {
								arrayToMerge = append(arrayToMerge, res.Value.(object.Object))
							}

							newElements := make([]object.Object, length+len(arrayToMerge), length+len(arrayToMerge))
							copy(newElements, array.Elements)

							for i, element := range arrayToMerge {
								newElements[i+length] = element
							}
							array.Elements = newElements
							return NULL
						},
						Parameters: 1,
						VarArgs:    true,
						Prototype:  true,
						Eval:       false,
					},
				},
			},
		},
		object.HashObj: {
			Pairs: map[object.HashKey]object.HashPair{
				khkp.GetKey("length").HashKey(): {
					Key: khkp.GetKey("length"),
					Value: &object.Builtin{
						Fn: func(token token.Token, env *object.Environment, args ...object.Object) object.Object {
							self, _ := env.Get("this")
							hash, _ := self.(*object.Hash)
							return &object.Integer{Value: float64(len(hash.Pairs))}
						},
						Parameters: 0,
						VarArgs:    false,
						Prototype:  true,
						Eval:       true,
					},
				},
				khkp.GetKey("keys").HashKey(): {
					Key: khkp.GetKey("keys"),
					Value: &object.Builtin{
						Fn: func(token token.Token, env *object.Environment, args ...object.Object) object.Object {
							self, _ := env.Get("this")
							hash, _ := self.(*object.Hash)

							keys := make([]object.Object, len(hash.Pairs))
							i := 0
							for _, v := range hash.Pairs {
								keys[i] = v.Key
								i++
							}

							return &object.Array{
								Elements: keys,
							}
						},
						Parameters: 0,
						VarArgs:    false,
						Prototype:  true,
						Eval:       true,
					},
				},
				khkp.GetKey("values").HashKey(): {
					Key: khkp.GetKey("values"),
					Value: &object.Builtin{
						Fn: func(token token.Token, env *object.Environment, args ...object.Object) object.Object {
							self, _ := env.Get("this")
							hash, _ := self.(*object.Hash)
							values := make([]object.Object, len(hash.Pairs))
							i := 0
							for _, v := range hash.Pairs {
								values[i] = v.Value
								i++
							}

							return &object.Array{
								Elements: values,
							}
						},
						Parameters: 0,
						VarArgs:    false,
						Prototype:  true,
						Eval:       true,
					},
				},
			},
		},
	}
}
