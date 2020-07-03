package evaluator

import (
	"Monkey/lexer"
	"Monkey/object"
	"Monkey/options"
	"Monkey/parser"
	"testing"
)

// Test Indexing
func TestArrayIndexExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{
			`[1, 2, 3][0]`,
			1,
		},
		{
			`[1, 2, 3][1]`,
			2,
		},
		{
			`[1, 2, 3][2]`,
			3,
		},
		{
			`let i = 0; [1][i]`,
			1,
		},
		{
			`[1, 2, 3][1 + 1]`,
			3,
		},
		{
			"let myArray = [1, 2, 3] myArray[2]",
			3,
		},
		{
			"let myArray = [1, 2, 3] myArray[0] + myArray[1] + myArray[2]",
			6,
		},
		{
			"let myArray = [1, 2, 3] let i = myArray[0] myArray[i]",
			2,
		},
		{
			"[1, 2, 3][3]",
			nil,
		},
		{
			`[1, 2, 3, 4, 5][1:3]`,
			[]float64{2, 3},
		},
		{
			`[1, 2, 3, 4, 5][:3]`,
			[]float64{1, 2, 3},
		},
		{
			`[1, 2, 3, 4, 5][1:]`,
			[]float64{2, 3, 4, 5},
		},
		{
			`[1, 2, 3, 4, 5][:]`,
			[]float64{1, 2, 3, 4, 5},
		},
		{
			`[1, 2, 3, 4, 5][1:-1]`,
			[]float64{2, 3, 4},
		},
		{
			`[1, 2, 3, 4, 5][-2]`,
			4,
		},
	}

	for _, tt := range tests {
		evaluated := CheckEval(tt.input)
		switch tt.expected.(type) {
		case int:
			i := tt.expected.(int)
			CheckIntegerObject(t, evaluated, float64(i))
		case nil:
			if _, ok := evaluated.(*object.Error); !ok {
				t.Errorf("should be error. got=%T",
					evaluated)
			}
		case []float64:
			i := tt.expected.([]float64)
			CheckArrayObject(t, evaluated, i)
		}
	}
}

func CheckArrayObject(t *testing.T, obj object.Object, expected []float64) bool {
	result, ok := obj.(*object.Array)
	if !ok {
		t.Errorf("obj not type Array. got=%T",
			obj)
		return false
	}

	if len(expected) != len(result.Elements) {
		t.Errorf("array does not contain enough values. got=%d, expected=%d",
			len(result.Elements), len(expected))
		return false
	}

	for index, ele := range result.Elements {
		val, ok := ele.(*object.Integer)
		if !ok {
			t.Errorf("ele(%d) is not type Integer. got=%T",
				index, ele)
			continue
		}

		if expected[index] != val.Value {
			t.Errorf("ele(%d) is does not contain %f. got=%f",
				index, expected[index], val.Value)
			continue
		}
	}
	return true
}

// Test eval arrays
func TestArrayLiteral(t *testing.T) {
	input := `[1, 2 * 2, 3 + 3]`

	evaluated := CheckEval(input)
	result, ok := evaluated.(*object.Array)
	if !ok {
		t.Fatalf("object is not type Array. got=%T (%+v)",
			evaluated, evaluated)
	}

	if len(result.Elements) != 3 {
		t.Fatalf("array has wrong number of elements. got=%d",
			len(result.Elements))
	}

	CheckIntegerObject(t, result.Elements[0], 1)
	CheckIntegerObject(t, result.Elements[1], 4)
	CheckIntegerObject(t, result.Elements[2], 6)
}

// Test Builtin functions
func TestBuiltinFunctions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`len("")`, 0},
		{`len("four")`, 4},
		{`len("hello world")`, 11},
		{`len(1)`, "argument to `len` not supported. got INTEGER"},
		{`len("one", "two")`, "wrong number of arguments. got=2, expected=1"},
	}

	for _, tt := range tests {
		evaluated := CheckEval(tt.input)

		switch expected := tt.expected.(type) {
		case int:
			CheckIntegerObject(t, evaluated, float64(expected))
		case string:
			errObj, ok := evaluated.(*object.Error)
			if !ok {
				t.Errorf("object is not Error. got=%T (%+v)",
					evaluated, evaluated)
				continue
			}
			if errObj.Message != expected {
				t.Errorf("wrong error message. expecte=%q, got=%q",
					expected, errObj.Message)
			}
		}
	}
}

// Test String +
func TestStringConcatenation(t *testing.T) {
	input := `
'Hello' + ' ' + 'World!'
`

	evaluated := CheckEval(input)
	str, ok := evaluated.(*object.String)
	if !ok {
		t.Fatalf("object is not type *object.String. got=%T (%+v)",
			evaluated, evaluated)
	}
	if str.Value != "Hello World!" {
		t.Errorf("String contains the wrong value. got=%q",
			str.Value)
	}
}

// Test Evaluating of Strings
func TestStringLiteral(t *testing.T) {
	input := `"Hello World!"`

	evaluated := CheckEval(input)
	str, ok := evaluated.(*object.String)
	if !ok {
		t.Fatalf("object is not type *object.String. got=%T (%+v)",
			evaluated, evaluated)
	}

	if str.Value != "Hello World!" {
		t.Errorf("String has wrong value. got=%q",
			str.Value)
	}
}

// Test Closure
func TestClosures(t *testing.T) {
	input := `
let newAdder = fn(x) {
	fn(y) {x + y}
};

let addTwo = newAdder(2);
addTwo(2)`

	CheckIntegerObject(t, CheckEval(input), 4)
}

// Test Calling functions
func TestFunctionCalls(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{"let identity = fn(x) { x } identity(5)", 5},
		{"let identity = fn(x) { return x } identity(5)", 5},
		{"let double = fn(x) { x * 2 } double(5)", 10},
		{"let add = fn(x, y) { x + y } add(5, 5)", 10},
		{"let add = fn(x, y) { x + y } add(5 + 5, add(5, 5))", 20},
		{"fn(x) { x }(5)", 5},
	}

	for _, tt := range tests {
		CheckIntegerObject(t, CheckEval(tt.input), tt.expected)
	}

}

// Test Functions
func TestFunctionObject(t *testing.T) {
	input := "fn(x) { x + 2 }"

	evaluated := CheckEval(input)
	fn, ok := evaluated.(*object.Function)
	if !ok {
		t.Fatalf("object is not type Function. got=%T (%+v)",
			evaluated, evaluated)
	}

	if len(fn.Parameters) != 1 {
		t.Fatalf("function have wrong numbers of parameters. Parameters=%+v",
			fn.Parameters)
	}

	if fn.Parameters[0].ToString() != "(x)" {
		t.Fatalf("parameter is not '(x)'. got=%q", fn.Parameters[0])
	}

	expectedBody := "{((x) + (2))}"

	if fn.Body.ToString() != expectedBody {
		t.Fatalf("body is not %q. got=%q",
			expectedBody, fn.Body.ToString())
	}
}

// Test a let statement
func TestLetStatement(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{"let a = 5; a", 5},
		{"let a = 5 * 5; a", 25},
		{"let a = 5; let b = a; b", 5},
		{"let a = 5; let b = a; let c = a + b + 5; c", 15},
	}

	for _, tt := range tests {
		CheckIntegerObject(t, CheckEval(tt.input), tt.expected)
	}
}

// Test Error handling
func TestErrorHandling(t *testing.T) {
	tests := []struct {
		input           string
		expectedMessage string
	}{
		{
			`"Hello" - 'World'`,
			"unknown operator: STRING - STRING",
		},
		{
			"5 + true",
			"type mismatch: INTEGER + BOOLEAN",
		},
		{
			"5 + true 5",
			"type mismatch: INTEGER + BOOLEAN",
		},
		{
			"-true",
			"unknown operator: -BOOLEAN",
		},
		{
			"true + false",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"5 true + false 5",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"if (10 > 1) { true + false }",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			`if 10 > 1 { if 10 > 1 { return true + false } return 1 }`,
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"foobar",
			"identifier not found: foobar",
		},
	}

	options.FatalErrors = true
	for _, tt := range tests {
		evaluated := CheckEval(tt.input)
		errObj, ok := evaluated.(*object.Error)
		if !ok {
			t.Errorf("no error object returned. got=%T(%+v)",
				evaluated, evaluated)
			continue
		}
		if errObj.Message != tt.expectedMessage {
			t.Errorf("wrong error message. expected=%q, got=%q",
				tt.expectedMessage, errObj.Message)
		}
	}
}

// Test return statements
func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{"5", 5},
		{"return 10", 10},
		{"return 10; 9", 10},
		{"return 2 * 5; 9", 10},
		{"let foo = 2 * 5; 9", 9},
		{"if 10 > 1 { return if 10 > 1 { 10 } 1 }", 10},
	}

	for _, tt := range tests {
		evaluated := CheckEval(tt.input)
		CheckIntegerObject(t, evaluated, tt.expected)
	}

}

// Test If Else Expressions
func TestIfElseExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"if true { 10 }", 10},
		{"if (false) { 10 }", nil},
		{"if 1 { 10 }", 10},
		{"if (1 < 2) { 10 }", 10},
		{"if 1 > 2 { 10 }", nil},
		{"if (1 > 2) { 10 } else { 20 }", 20},
		{"if (1 < 2) { 10 } else { 20 }", 10},
		{"if 1 > 2 { 1 } else if 1 == 2 { 20 } else { 14 }", 14},
		{"if 1 > 2 { 1 } else if 1 != 2 { 20 } else { 14 }", 20},
	}
	for _, tt := range tests {
		evaluated := CheckEval(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			CheckIntegerObject(t, evaluated, float64(integer))
		} else {
			CheckNullObject(t, evaluated)
		}
	}
}

func CheckNullObject(t *testing.T, obj object.Object) bool {
	if obj != NULL {
		t.Errorf("object is not NULL. got=%T (%+v)", obj, obj)
		return false
	}
	return true
}

// Test Eval of Integer
func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{"5", 5},
		{"10", 10},
		{"-10", -10},
		{"-5", -5},
		{"+5", +5},
		{"+10", +10},
		{"-10.5", -10.5},
		{"-0.", 0},
		{"1.5 + 2.5", 4},
		{"3.2 * 3", 9.6},
		{"4 / 2.0", 2},
		{"5.0 % 2.0", 1},
		{" 5 + 5 + 5 + 5 - 10", 10},
		{"2 * 2 * 2 * 2 * 2", 32},
		{"-50 + 100 + -50", 0},
		{"5 * 2 + 10", 20},
		{"5 + 2 * 10", 25},
		{"20 + 2 * -10", 0},
		{"50 / 2 * 2 + 10", 60},
		{"2 * (5 + 10)", 30},
		{"3 * 3 * 3 + 10", 37},
		{"3 * (3 * 3) + 10", 37},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
	}

	for _, tt := range tests {
		evaluated := CheckEval(tt.input)
		CheckIntegerObject(t, evaluated, tt.expected)
	}
}

// Eval an input
func CheckEval(input string) object.Object {
	l := lexer.New(input, "testEval")
	p := parser.New(l)
	program := p.ParseProgram()
	env := object.NewEnvironment()

	return Eval(program, env)
}

// Check if integer object
func CheckIntegerObject(t *testing.T, obj object.Object, expected float64) bool {
	result, ok := obj.(*object.Integer)

	if !ok {
		t.Errorf("object is not type Integer. got=%T (%+v)",
			obj, obj)
		return false
	}
	if !parser.AlmostEqual(result.Value, expected) {
		t.Errorf("object has wrong value. got=%f, expect=%f",
			result.Value, expected)
		return false
	}
	return true

}

// Test Eval of Booleans
func TestEvalBooleanExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
		{"1 < 2", true},
		{"1 > 2", false},
		{"1 < 1", false},
		{"1 > 1", false},
		{"1 == 1", true},
		{"1 != 1", false},
		{"1 == 2", false},
		{"1 != 2", true},
		{"true == true", true},
		{"false == false", true},
		{"true == false", false},
		{"true != false", true},
		{"false != true", true},
		{"1 < 2 == true", true},
		{"1 < 2 == false", false},
		{"1 > 2 == true", false},
		{"1 > 2 == false", true},

		{"1 == false", false},
		{"1 == true", true},
		{"0 == true", false},
		{"0 == false", true},
	}

	for _, tt := range tests {
		evaluated := CheckEval(tt.input)
		CheckBooleanObject(t, evaluated, tt.expected)
	}
}

// Check boolean object
func CheckBooleanObject(t *testing.T, obj object.Object, expected bool) bool {
	result, ok := obj.(*object.Boolean)
	if !ok {
		t.Errorf("object is not type Boolean. got=%T (%+v)",
			obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. got=%t, want=%t",
			result.Value, expected)
		return false
	}
	return true
}

// Test Logic Gates
func TestLogicGates(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true and true", true},
		{"true and false", false},
		{"false and false", false},
		{"true or true", true},
		{"true or false", true},
		{"false or false", false},
		{"true xor true", false},
		{"true xor false", true},
		{"false xor false", false},

		{"true && true", true},
		{"true && false", false},
		{"false && false", false},
		{"true || true", true},
		{"true || false", true},
		{"false || false", false},
		{"true !| true", false},
		{"true !| false", true},
		{"false !| false", false},
	}

	for _, tt := range tests {
		evaluated := CheckEval(tt.input)
		CheckBooleanObject(t, evaluated, tt.expected)
	}
}

// Test ! eval
func TestBangOperator(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"!true", false},
		{"!false", true},
		{"!5", false},
		{"!!true", true},
		{"!!1", true},
	}

	for _, tt := range tests {
		evaluated := CheckEval(tt.input)
		CheckBooleanObject(t, evaluated, tt.expected)
	}
}
