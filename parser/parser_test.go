package parser

import (
	"Monkey/ast"
	"Monkey/lexer"
	"testing"
)

// Test the parsing of the let statements
func TestLetStatements(t *testing.T) {
	input := `
let x = 5;
let y = 10;
let foobar = 848484;`

	l := lexer.New(input, "test")
	p := New(l)

	program := p.ParseProgram()
	CheckParserErrors(t, p)
	if program == nil {
		t.Fatalf("ParseProgram returned nil")
	}

	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements, got=%d",
			len(program.Statements))
	}

	tests := []struct {
		expectedIdentifier string
	}{
		{"x"},
		{"y"},
		{"foobar"},
	}

	// Test if each of the let statement is correct
	for i, tt := range tests {
		stmt := program.Statements[i]
		if !CheckLetStatement(t, stmt, tt.expectedIdentifier) {
			return
		}
	}
}

// Test a single let statement
func CheckLetStatement(t *testing.T, s ast.Statement, name string) bool {
	// Check if keyword is correct
	if s.TokenLiteral() != "let" {
		t.Errorf("s.TokenLiteral not 'let'. got=%q", s.TokenLiteral())
		return false
	}

	// Check if casting is successful
	letStmt, ok := s.(*ast.LetStatement)
	if !ok {
		t.Errorf("s not type *ast.LetStatment. got=%T", s)
		return false
	}

	// Check if variable name matches
	if letStmt.Name.Value != name {
		t.Errorf("letStmt.Name.Value not '%s'. got=%s",
			name, letStmt.Name.Value)
		return false
	}

	// Check if debug variable name value matches
	if letStmt.Name.TokenLiteral() != name {
		t.Errorf("s.Name not '%s'. got=%s",
			name, letStmt.Name.TokenLiteral())
		return false
	}

	//TODO: Check value

	return true
}

// Verify that there are no errors
func CheckParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()

	// Return if errors is empty
	if len(errors) == 0 {
		return
	}

	// Else print out the parseErrors
	t.Errorf("parser has %d errors", len(errors))
	for _, err := range errors {
		// Pointer accessing here
		t.Errorf("parser error: %q, at %d:%d, in %q",
			err.Message, err.RowNumber, err.ColumnNumber, err.Filename)
	}
	// Fail the test
	t.FailNow()
}
