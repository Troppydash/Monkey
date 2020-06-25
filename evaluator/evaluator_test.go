package evaluator

import (
	"Monkey/lexer"
	"Monkey/object"
	"Monkey/parser"
	"testing"
)

// Test Eval of Integer
func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"10", 10},
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
func CheckIntegerObject(t *testing.T, obj object.Object, expected int64) bool {
	result, ok := obj.(*object.Integer)
	if !ok {
		t.Errorf("object is not type Intger. got=%T (%+v)",
			obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. got=%d, expect=%d",
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
