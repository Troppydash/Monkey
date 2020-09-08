package evaluator

import (
	"Monkey/ast"
	"Monkey/lexer"
	"Monkey/object"
	"Monkey/options"
	"Monkey/parser"
	"testing"
)

func TestDefineMacros(t *testing.T) {
	input := `
	let number = 1
	let function = fn(x, y) { x + y }
	let mymacro = macro(x, y) { x + y }
	`

	env := object.NewEnvironment()
	program := testParseProgram(input)
	DefineMacros(program, env)

	// Didn't eval it, shouldn't be defined
	if len(program.Statements) != 2 {
		t.Fatalf("Wrong number of statement. got=%d",
			len(program.Statements))
	}

	_, ok := env.Get("number")
	if ok {
		t.Fatalf("number should not be defined")
	}

	_, ok = env.Get("function")
	if ok {
		t.Fatalf("function should not be defined")
	}

	// This should be defined from the call to DefineMacros
	obj, ok := env.Get("mymacro")
	if !ok {
		t.Fatalf("macro not in environment")
	}

	macro, ok := obj.(*object.Macro)
	if !ok {
		t.Fatalf("object is not Macro. got=%T (%+v)",
			obj, obj)
	}

	if macro.Parameters[0].Value != "x" {
		t.Fatalf("parameter is not 'x'. got=%q", macro.Parameters[0])
	}

	if macro.Parameters[1].Value != "y" {
		t.Fatalf("parameter is not 'y'. got=%q", macro.Parameters[1])
	}

	expectedBody := "{((x) + (y))}"

	if macro.Body.ToString() != expectedBody {
		t.Fatalf("body is not %q. got=%q", expectedBody, macro.Body.ToString())
	}
}

// TODO: Make parsing functions like this
// fn thisIsAFunction() {}

func testParseProgram(input string) *ast.Program {
	l := lexer.New(input, "TestParseProgram")
	p := parser.New(l)
	return p.ParseProgram()
}

func TestExpandMacros(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			`let infix = macro() { quote(1 + 2) }
				infix()`,
			`(1 + 2)`,
		},
		{
			`let reverse = macro(a, b) { quote(unquote(b) - unquote(a)) }
				reverse(2 + 2, 10 - 5)`,
			`(10 - 5) - (2 + 2)`,
		},
		{
			`
let unless = macro(cond, cons, alter) {
	quote(if (!(unquote(cond))) {
		unquote(cons)
	} else {
		unquote(alter)
	})
}
unless(10 > 5, write("not"), write("greater"))
`,
			`if (!(10 > 5)) { write("not") } else { write("greater") }`,
		},
		//		{
		//			`let ifnot$ = macro(cond, consq, alt) {
		//    quote(if (!(unquote(cond))) {
		//        if typeof(unquote(consq)) == FUNCTION {
		//			tmp()
		//        } else {
		//        // TODO: Fix this
		//            unquote(consq)
		//        }
		//    })
		//}
		//ifnot$(1 > 12, #{
		//                   writeLine("Smaller")
		//               }, #{
		//                      writeLine("Bigger")
		//                  })
		//`,
		//			`12`,
		//		},
	}

	options.NicerToString = true

	for _, tt := range tests {
		expected := testParseProgram(tt.expected)
		program := testParseProgram(tt.input)

		env := object.NewEnvironment()
		DefineMacros(program, env)
		expanded := ExpandMacros(program, env)

		if expanded.ToString() != expected.ToString() {
			t.Errorf("not equal. want=%q, got=%q",
				expanded.ToString(), expected.ToString())
		}
	}
}
