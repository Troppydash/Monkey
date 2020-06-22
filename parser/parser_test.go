package parser

import (
	"Monkey/ast"
	"Monkey/lexer"
	"fmt"
	"testing"
)

// Testing parsing precedence expression
func TestOperatorPrecedenceParsing(t *testing.T) {
	// Test Cases
	tests := []struct {
		input    string
		expected string
	}{
		{"-a * b", "((-a) * b)"},
		{"!+a", "(!(+a))"},
		{"1 + 1 * 1", "(1 + (1 * 1))"},
		{
			"a + b * c + d / e - f",
			"(((a + (b * c)) + (d / e)) - f)",
		},
		{
			"3 + 4; -5 * 5",
			"(3 + 4)((-5) * 5)",
		},
		{
			"5 > 4 == 3 < 4",
			"((5 > 4) == (3 < 4))",
		},
		{
			"5 < 4 != 3 > 4",
			"((5 < 4) != (3 > 4))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input, "testPrecedence")
		p := New(l)
		program := p.ParseProgram()
		p.CheckParserErrors(t)

		actual := program.ToString()

		if actual != tt.expected {
			t.Errorf("expected=%q. got=%q",
				tt.expected, actual)
		}
	}
}

// Test parsing of infix expressions
func TestParsingInfixExpressions(t *testing.T) {
	infixTests := []struct {
		input      string
		leftValue  int64
		operator   string
		rightValue int64
	}{
		{"5 + 5", 5, "+", 5},
		{"5 - 5", 5, "-", 5},
		{"5 * 5", 5, "*", 5},
		{"5 / 5", 5, "/", 5},
		{"5 > 5", 5, ">", 5},
		{"5 < 5", 5, "<", 5},
		{"5 >= 5", 5, ">=", 5},
		{"5 <= 5", 5, "<=", 5},
		{"5 == 5", 5, "==", 5},
		{"5 != 5", 5, "!=", 5},
	}

	for _, tt := range infixTests {
		l := lexer.New(tt.input, "testInfix")
		p := New(l)
		program := p.ParseProgram()
		p.CheckParserErrors(t)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain %d statements. got=%d",
				1, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not type *ast.ExpressionStatement. got=%T",
				program.Statements[0])
		}

		exp, ok := stmt.Expression.(*ast.InfixExpression)
		if !ok {
			t.Fatalf("exp is not ast.InfixExpression. got=%T",
				stmt.Expression)
		}

		if !CheckIntegerLiteral(t, exp.Left, tt.leftValue) {
			return
		}

		if exp.Operator != tt.operator {
			t.Fatalf("exp.Operator is not %q. got=%s",
				tt.operator, exp.Operator)
		}

		if !CheckIntegerLiteral(t, exp.Right, tt.rightValue) {
			return
		}
	}
}

// Test parsing of prefix
func TestParsingPrefixExpressions(t *testing.T) {
	prefixTests := []struct {
		input        string
		operator     string
		integerValue int64
	}{
		{"!5", "!", 5},
		{"-15", "-", 15},
		{"+15", "+", 15},
	}

	for _, tt := range prefixTests {
		l := lexer.New(tt.input, "testPrefix")
		p := New(l)
		program := p.ParseProgram()
		p.CheckParserErrors(t)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain %d statements. got=%d",
				1, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not type ast.ExpressionStatement. got=%T",
				program.Statements[0])
		}

		exp, ok := stmt.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("stmt is not type ast.PrefixExpression. got=%T",
				stmt.Expression)
		}
		if exp.Operator != tt.operator {
			t.Fatalf("exp.Operator is not %q. got=%q",
				tt.operator, exp.Operator)
		}

		if !CheckIntegerLiteral(t, exp.Right, tt.integerValue) {
			return
		}
	}
}

func CheckIntegerLiteral(t *testing.T, il ast.Expression, value int64) bool {
	integ, ok := il.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("il not type *ast.IntegerLiteral. got=%T", il)
		return false
	}

	if integ.Value != value {
		t.Errorf("integ.Value not %d. got=%d",
			value, integ.Value)
		return false
	}

	if integ.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Errorf("integ.TokenLiteral not %d. got=%s",
			value, integ.TokenLiteral())
		return false
	}

	return true
}

// Test parsing of literals
func TestIntegerLiteralExpression(t *testing.T) {
	input := "5;"

	l := lexer.New(input, "testLiteral")
	p := New(l)
	program := p.ParseProgram()
	p.CheckParserErrors(t)

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statments. got=%d",
			len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] not type ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	literal, ok := stmt.Expression.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("expression not type *ast.IntergerLiterla. got=%T",
			stmt.Expression)
	}
	if literal.Value != 5 {
		t.Errorf("literal.Value not %d. got=%d",
			5, literal.Value)
	}
	if literal.TokenLiteral() != "5" {
		t.Errorf("literal.TokenLiteral not %s. got=%s",
			"5", literal.TokenLiteral())
	}
}

// Test parsing of the identifier
func TestIdentifierExpression(t *testing.T) {
	input := "foobar;"

	l := lexer.New(input, "testIdentifier")
	p := New(l)
	program := p.ParseProgram()
	p.CheckParserErrors(t)

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d",
			len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not type ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	ident, ok := stmt.Expression.(*ast.Identifier)
	if !ok {
		t.Fatalf("stmt.Expression not type *ast.Identifier. got=%T",
			stmt.Expression)
	}
	if ident.Value != "foobar" {
		t.Errorf("ident.Value not foobar. got=%s",
			ident.Value)
	}
	if ident.TokenLiteral() != "foobar" {
		t.Errorf("ident.TokenLiteral not foobar. got=%s",
			ident.TokenLiteral())
	}

}

// Test the parsing of the return statements
func TestReturnStatements(t *testing.T) {
	input := `
return 5;
return 10;
return 921391;`

	l := lexer.New(input, "testReturn")
	p := New(l)

	program := p.ParseProgram()
	p.CheckParserErrors(t)

	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contian 3 statements. got=%d",
			len(program.Statements))
	}

	for _, stmt := range program.Statements {
		returnStmt, ok := stmt.(*ast.ReturnStatement)
		if !ok {
			t.Errorf("stmt not *ast.returnStatment. got=%T", stmt)
			continue
		}
		if returnStmt.TokenLiteral() != "return" {
			t.Errorf("returnStmt.TokenLiteral not 'return', got=%q",
				returnStmt.TokenLiteral())
		}
	}
}

// Test the parsing of the let statements
func TestLetStatements(t *testing.T) {
	input := `
let x = 5;
let y = 10;
let foobar = 848484;`

	l := lexer.New(input, "testLet")
	p := New(l)

	program := p.ParseProgram()
	p.CheckParserErrors(t)
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
func (p *Parser) CheckParserErrors(t *testing.T) {
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
