package evaluator

import (
	"Monkey/object"
	"Monkey/token"
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Argument not supported error
func ArgumentNotSupported(method string, typ interface{}, token token.Token) *object.Error {
	return NewError(token.ToTokenData(), "argument to `%s` not supported. got %s",
		method, typ)
}

// Incorrect argument amount error
func WrongArgumentsAmount(method string, got interface{}, expected interface{}, token token.Token) *object.Error {
	return NewError(token.ToTokenData(), "wrong number of arguments for method `%s`. got=%d, expected=%s",
		method, got, expected)
}

// Prohibited Value error
func ProhibitedValue(method string, value interface{}, reason interface{}, token token.Token) *object.Error {
	return NewError(token.ToTokenData(), "prohibited value of arguments for method `%s`. got=%v, reason=%s",
		method, value, reason)
}

var builtins map[string]*object.Builtin

func builtinsInit() {
	builtins = map[string]*object.Builtin{
		// TODO: Math Functions

		// TODO: Array Push and Pop

		// Array
		"len": {
			Fn: func(token token.Token, args ...object.Object) object.Object {
				if len(args) != 1 {
					return WrongArgumentsAmount("len", len(args), "1", token)
				}

				switch arg := args[0].(type) {
				case *object.String:
					return &object.Integer{Value: float64(len(arg.Value))}
				case *object.Array:
					return &object.Integer{Value: float64(len(arg.Elements))}
				default:
					return ArgumentNotSupported("len", args[0].Type(), token)
				}
			},
		},
		"range": {
			Fn: func(token token.Token, args ...object.Object) object.Object {

				switch len(args) {

				case 0:
					return &object.Array{Elements: []object.Object{}}
				case 1:
					amount, ok := args[0].(*object.Integer)
					if !ok {
						return ArgumentNotSupported("range", args[0].Type(), token)
					}

					var eles []object.Object
					for i := 0; i < int(amount.Value); i++ {
						eles = append(eles, &object.Integer{Value: float64(i)})
					}

					return &object.Array{Elements: eles}

				case 2:
					amount, ok := args[1].(*object.Integer)
					if !ok {
						return ArgumentNotSupported("range", args[0].Type(), token)
					}

					starting, ok := args[0].(*object.Integer)
					if !ok {
						return ArgumentNotSupported("range", args[0].Type(), token)
					}

					var eles []object.Object
					if amount.Value < starting.Value {
						for i := starting.Value; i > amount.Value; i -= 1 {
							eles = append(eles, &object.Integer{Value: i})
						}
					} else {
						for i := starting.Value; i < amount.Value; i += 1 {
							eles = append(eles, &object.Integer{Value: i})
						}
					}

					return &object.Array{Elements: eles}

				case 3:
					skip, ok := args[2].(*object.Integer)
					if !ok {
						return ArgumentNotSupported("range", args[2].Type(), token)
					}
					if skip.Value == 0 {
						return ProhibitedValue("range", skip.Value, "range would loop forever", token)
					}

					amount, ok := args[1].(*object.Integer)
					if !ok {
						return ArgumentNotSupported("range", args[1].Type(), token)
					}

					starting, ok := args[0].(*object.Integer)
					if !ok {
						return ArgumentNotSupported("range", args[0].Type(), token)
					}

					if skip.Value < 0 {
						return ProhibitedValue("range", skip.Value, "skip cannot be negative", token)
					}

					var eles []object.Object

					if amount.Value < starting.Value {
						for i := starting.Value; i > amount.Value; i -= skip.Value {
							eles = append(eles, &object.Integer{Value: i})
						}
					} else {
						for i := starting.Value; i < amount.Value; i += skip.Value {
							eles = append(eles, &object.Integer{Value: i})
						}
					}

					return &object.Array{Elements: eles}

				default:
					return WrongArgumentsAmount("range", len(args), "1-3", token)
				}
			},
		},
		"push": {
			Fn: func(token token.Token, args ...object.Object) object.Object {
				if len(args) != 2 {
					return WrongArgumentsAmount("push", len(args), "2", token)
				}

				if args[0].Type() != object.ARRAY_OBJ {
					return ArgumentNotSupported("push", args[0].Type(), token)
				}

				arr := args[0].(*object.Array)
				length := len(arr.Elements)

				newElements := make([]object.Object, length+1, length+1)
				copy(newElements, arr.Elements)
				newElements[length] = args[1]

				return &object.Array{Elements: newElements}
			},
		},
		// TODO: Pop, map, forEach, repeat

		"for": {
			// TODO: Implem better for loops
			Fn: func(token token.Token, args ...object.Object) object.Object {
				if len(args) != 1 {
					return WrongArgumentsAmount("repeat", len(args), "1", token)
				}

				if args[0].Type() != object.FUNCTION_OBJ {
					return ArgumentNotSupported("repeat", args[0].Type(), token)
				}

				fn := args[0].(*object.Function)

				var result object.Object = FALSE

				if len(fn.Parameters) == 1 {
					times := 0
					for !IsTruthful(result) {
						result = ApplyFunction(token, fn, []object.Object{
							&object.Integer{Value: float64(times)},
						})
						times += 1
					}
				} else {
					for !IsTruthful(result) {
						result = ApplyFunction(token, fn, []object.Object{})
					}
				}

				return result
			},
		},

		// IO
		"format": {
			// TODO: Implem
			Fn: func(token token.Token, args ...object.Object) object.Object {
				return NULL
			},
		},
		"write": {
			Fn: func(token token.Token, args ...object.Object) object.Object {
				var out []string
				for _, obj := range args {
					out = append(out, obj.Inspect())
				}

				fmt.Print(strings.Join(out, " "))
				return NULL
			},
		},
		"writeLine": {
			Fn: func(token token.Token, args ...object.Object) object.Object {
				var out []string
				for _, obj := range args {
					out = append(out, obj.Inspect())
				}

				fmt.Println(strings.Join(out, " "))
				return NULL
			},
		},
		"take": {
			Fn: func(token token.Token, args ...object.Object) object.Object {
				if len(args) > 1 {
					return WrongArgumentsAmount("take", len(args), "0-1", token)
				}

				reader := bufio.NewReader(os.Stdin)

				if len(args) == 1 {
					fmt.Print(args[0].Inspect() + " > ")
				}

				text, _, _ := reader.ReadLine()
				return &object.String{
					Value: string(text),
				}
			},
		},
		"takeLine": {
			Fn: func(token token.Token, args ...object.Object) object.Object {
				if len(args) > 1 {
					return WrongArgumentsAmount("takeLine", len(args), "0-1", token)
				}

				reader := bufio.NewReader(os.Stdin)

				if len(args) == 1 {
					fmt.Println(args[0].Inspect() + " > ")
				}
				text, _, _ := reader.ReadLine()
				return &object.String{
					Value: string(text),
				}
			},
		},

		// TODO: Add make error && panic/fatalError
		// Checking
		"error?": {
			Fn: func(token token.Token, args ...object.Object) object.Object {
				if len(args) != 1 {
					return TRUE
				}

				switch args[0].(type) {
				case *object.Error:
					return TRUE
				default:
					return FALSE
				}
			},
		},
		"null?": {
			Fn: func(token token.Token, args ...object.Object) object.Object {
				if len(args) != 1 {
					return FALSE
				}

				return NativeBoolToBooleanObject(args[0].Type() == object.NULL_OBJ)
			},
		},

		// Casting
		"bool!": {
			Fn: func(token token.Token, args ...object.Object) object.Object {
				if len(args) != 1 {
					return WrongArgumentsAmount("bool!", len(args), "1", token)
				}
				target := args[0]

				switch target.(type) {
				case *object.Boolean:
					return target
				default:
					return NativeBoolToBooleanObject(IsTruthful(args[0]))
				}
			},
		},
		"string!": {
			Fn: func(token token.Token, args ...object.Object) object.Object {
				if len(args) > 1 {
					return WrongArgumentsAmount("string!", len(args), "1", token)
				}

				target := args[0]

				switch target.(type) {
				case *object.Error:
					err := target.(*object.Error)
					return &object.String{
						Value: err.Message,
					}
				case *object.String:
					return target
				default:
					return &object.String{
						Value: target.Inspect(),
					}
				}
			},
		},
		"number!": {
			Fn: func(token token.Token, args ...object.Object) object.Object {
				if len(args) != 1 {
					return WrongArgumentsAmount("number!", len(args), "1", token)
				}
				target := args[0]

				switch target.(type) {
				case *object.Integer:
					return target

				case *object.Boolean:
					b := target.(*object.Boolean)
					if b.Value {
						return &object.Integer{
							Value: 1,
						}
					} else {
						return &object.Integer{
							Value: 0,
						}
					}

				case *object.String:
					s := target.(*object.String)
					v, err := strconv.ParseFloat(s.Value, 64)
					if err != nil {
						return NewError(token.ToTokenData(), "casting to number not successful. got=%s",
							s.Value)
					}

					return &object.Integer{
						Value: v,
					}
				}

				return ArgumentNotSupported("number", args[0].Type(), token)
			},
		},
	}
}
