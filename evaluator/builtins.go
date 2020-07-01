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

var builtins = map[string]*object.Builtin{
	// TODO: Math Functions

	// Array
	"len": {
		Fn: func(token token.Token, args ...object.Object) object.Object {
			if len(args) != 1 {
				return NewError(token.ToTokenData(), "wrong number of arguments. got=%d, expected=1",
					len(args))
			}

			switch arg := args[0].(type) {
			case *object.String:
				return &object.Integer{Value: float64(len(arg.Value))}
			default:
				return NewError(token.ToTokenData(), "argument to `len` not supported. got %s",
					args[0].Type())
			}
		},
	},

	// IO
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
				return NewError(token.ToTokenData(), "wrong number of arguments. got=%d, expected=0/1",
					len(args))
			}

			reader := bufio.NewReader(os.Stdin)

			if len(args) == 1 {
				fmt.Print(args[0].Inspect() + "> ")
			}
			text, _ := reader.ReadString('\n')
			return &object.String{
				Value: text[:len(text)-1],
			}
		},
	},
	"takeLine": {
		Fn: func(token token.Token, args ...object.Object) object.Object {
			if len(args) > 1 {
				return NewError(token.ToTokenData(), "wrong number of arguments. got=%d, expected=0/1",
					len(args))
			}

			reader := bufio.NewReader(os.Stdin)

			if len(args) == 1 {
				fmt.Println(args[0].Inspect() + "> ")
			}
			text, _ := reader.ReadString('\n')
			return &object.String{
				Value: text[:len(text)-1],
			}
		},
	},

	// Checking
	"error": {
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
	"null": {
		Fn: func(token token.Token, args ...object.Object) object.Object {
			if len(args) != 1 {
				return FALSE
			}

			return NativeBoolToBooleanObject(args[0].Type() == object.NULL_OBJ)
		},
	},

	// Casting
	"bool": {
		Fn: func(token token.Token, args ...object.Object) object.Object {
			if len(args) != 1 {
				return NewError(token.ToTokenData(), "wrong number of arguments. got=%d, expected=1",
					len(args))
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
	"string": {
		Fn: func(token token.Token, args ...object.Object) object.Object {
			if len(args) > 1 {
				return NewError(token.ToTokenData(), "wrong number of arguments. got=%d, expected=1",
					len(args))
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
	"number": {
		Fn: func(token token.Token, args ...object.Object) object.Object {
			if len(args) != 1 {
				return NewError(token.ToTokenData(), "wrong number of arguments. got=%d, expected=1",
					len(args))
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

			return NewError(token.ToTokenData(), "argument to `number` not supported. got %s",
				args[0].Type())
		},
	},
}
