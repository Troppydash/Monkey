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
	return NewFatalError(token.ToTokenData(), "argument to `%s` not supported. got %s",
		method, typ)
}

// Incorrect argument amount error
func WrongArgumentsAmount(method string, got interface{}, expected interface{}, token token.Token) *object.Error {
	return NewFatalError(token.ToTokenData(), "wrong number of arguments for method `%s`. got=%d, expected=%s",
		method, got, expected)
}

// Prohibited Value error
func ProhibitedValue(method string, value interface{}, reason interface{}, token token.Token) *object.Error {
	return NewFatalError(token.ToTokenData(), "prohibited value of arguments for method `%s`. got=%v, reason=%s",
		method, value, reason)
}

var builtins map[string]*object.Builtin

func init() {
	builtins = map[string]*object.Builtin{
		// TODO: Math Functions

		"typeof": {
			Fn: func(token token.Token, env *object.Environment, args ...object.Object) object.Object {
				if len(args) != 1 {
					return WrongArgumentsAmount("typeof", len(args), "1", token)
				}

				return &object.String{
					Value: string(args[0].Type()),
				}
			},
		},

		// TODO: Reorder function parameter
		"include": {
			Fn: func(token token.Token, env *object.Environment, args ...object.Object) object.Object {
				if len(args) != 1 {
					return WrongArgumentsAmount("include", len(args), "1", token)
				}

				str, ok := args[0].(*object.String)
				if !ok {
					return NULL
				}

				filename := str.Value

				err := LinkAndEval(filename, env)
				//old := tmp.CurrentProcessingFileDirectory
				//abs := runner.GetInstance().ToAbsolute(filename)
				//p, e := runner.GetInstance().CompileAbs(abs)
				if err != nil {
					return NewFatalError(token.ToTokenData(), "Failed to compile file %q\n", filename)

				}
				//Eval(p, env)
				//runner.GetInstance().Pop(abs)
				//tmp.CurrentProcessingFileDirectory = old

				return NULL
			},
		},
		// Array
		"len": {
			Fn: func(token token.Token, env *object.Environment, args ...object.Object) object.Object {
				if len(args) != 1 {
					return WrongArgumentsAmount("len", len(args), "1", token)
				}

				switch arg := args[0].(type) {
				case *object.String:
					return &object.Integer{Value: float64(len(arg.Value))}
				case *object.Array:
					return &object.Integer{Value: float64(len(arg.Elements))}
				case *object.Hash:
					return &object.Integer{Value: float64(len(arg.Pairs))}
				default:
					return ArgumentNotSupported("len", args[0].Type(), token)
				}
			},
		},
		"keys": {
			Fn: func(token token.Token, env *object.Environment, args ...object.Object) object.Object {
				if len(args) != 1 {
					return WrongArgumentsAmount("keys", len(args), "1", token)
				}

				hash, ok := args[0].(*object.Hash)
				if !ok {
					return ArgumentNotSupported("keys", args[0].Type(), token)

				}
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
		},

		"range": {
			Fn: func(token token.Token, env *object.Environment, args ...object.Object) object.Object {

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
			Fn: func(token token.Token, env *object.Environment, args ...object.Object) object.Object {
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
		"add": {
			Fn: func(token token.Token, env *object.Environment, args ...object.Object) object.Object {
				if len(args) != 2 {
					return WrongArgumentsAmount("add", len(args), "2", token)
				}

				if args[0].Type() != object.ARRAY_OBJ {
					return ArgumentNotSupported("add", args[0].Type(), token)
				}

				arr := args[0].(*object.Array)

				arr.Elements = append(arr.Elements, args[1])
				return NULL
			},
		},
		// TODO: Pop, map, forEach, repeat

		"loop": {
			Fn: func(token token.Token, env *object.Environment, args ...object.Object) object.Object {
				if !(len(args) <= 2) {
					return WrongArgumentsAmount("loop", len(args), "1-2", token)
				}

				if args[0].Type() != object.FUNCTION_OBJ {
					return ArgumentNotSupported("loop", args[0].Type(), token)
				}

				fn := args[0].(*object.Function)

				t := &object.Integer{Value: float64(0)}

				switch len(args) {
				case 1:
					var result object.Object
					for result != BREAK {
						//env.Store("t", &object.Integer{Value: float64(t)})
						t.Value += 1
						result = ApplyFunction(token, fn, []object.Object{
							t,
						}, env)
					}

				case 2:
					val, ok := args[1].(*object.Integer)
					if !ok {
						return ArgumentNotSupported("loop", args[1], token)
					}

					times := val.Value

					for ; t.Value < times; t.Value++ {
						//env.Store("t", &object.Integer{Value: float64(t)})
						ApplyFunction(token, fn, []object.Object{
							t,
						}, env)
					}
				}
				return NULL
			},
		},

		"while": {
			Fn: func(token token.Token, env *object.Environment, args ...object.Object) object.Object {
				if !(len(args) == 2) {
					return WrongArgumentsAmount("while", len(args), "2", token)
				}

				if args[0].Type() != object.FUNCTION_OBJ {
					return ArgumentNotSupported("while", args[0].Type(), token)
				}
				if args[1].Type() != object.FUNCTION_OBJ {
					return ArgumentNotSupported("while", args[1].Type(), token)
				}

				fn := args[0].(*object.Function)
				exe := args[1].(*object.Function)

				result := ApplyFunction(token, fn, []object.Object{}, env)
				for IsTruthful(result) {
					ApplyFunction(token, exe, []object.Object{}, env)
					result = ApplyFunction(token, fn, []object.Object{}, env)
				}
				return NULL
			},
		},

		// IO
		"format": {
			// TODO: Implem
			Fn: func(token token.Token, env *object.Environment, args ...object.Object) object.Object {
				return NULL
			},
		},
		"write": {
			Fn: func(token token.Token, env *object.Environment, args ...object.Object) object.Object {
				var out []string
				for _, obj := range args {
					out = append(out, obj.Inspect())
				}

				fmt.Print(strings.Join(out, " "))
				return NULL
			},
		},
		"writeLine": {
			Fn: func(token token.Token, env *object.Environment, args ...object.Object) object.Object {
				var out []string
				for _, obj := range args {
					out = append(out, obj.Inspect())
				}

				fmt.Println(strings.Join(out, " "))
				return NULL
			},
		},
		"take": {
			Fn: func(token token.Token, env *object.Environment, args ...object.Object) object.Object {
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
			Fn: func(token token.Token, env *object.Environment, args ...object.Object) object.Object {
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
			Fn: func(token token.Token, env *object.Environment, args ...object.Object) object.Object {
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
			Fn: func(token token.Token, env *object.Environment, args ...object.Object) object.Object {
				if len(args) != 1 {
					return FALSE
				}

				return NativeBoolToBooleanObject(args[0].Type() == object.NULL_OBJ)
			},
		},

		// Casting
		"bool!": {
			Fn: func(token token.Token, env *object.Environment, args ...object.Object) object.Object {
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
			Fn: func(token token.Token, env *object.Environment, args ...object.Object) object.Object {
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
			Fn: func(token token.Token, env *object.Environment, args ...object.Object) object.Object {
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
