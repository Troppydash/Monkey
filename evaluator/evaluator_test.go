package evaluator

import (
	"Monkey/lexer"
	"Monkey/object"
	"Monkey/parser"
	"testing"
)

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

	return Eval(program)
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
