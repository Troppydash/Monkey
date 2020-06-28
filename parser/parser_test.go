package parser

import (
	"Monkey/ast"
	"Monkey/lexer"
	"Monkey/options"
	"fmt"
	"strconv"
	"strings"
	"testing"
)

func FormatFloat(t float64) string {
	return strconv.FormatFloat(t, 'f', -1, 64)
}

func TestFloatingPointNumbers(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{"1.5", 1.5},
		{"300.0", 300.0},
		{"0.", 0},
		{"0.1", 0.1},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input, "testFloat")
		p := New(l)
		program := p.ParseProgram()
		p.CheckParserErrors(t)

		if len(program.Statements) != 1 {
			t.Fatalf("program does not contain enough statements. got=%d",
				len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] not type *ast.ExpressionStatement. got=%T",
				program.Statements[0])
		}

		CheckIntegerLiteral(t, stmt.Expression, tt.expected)
	}
}

func TestPrintExpStmtParsing(t *testing.T) {
	input := `12;`

	l := lexer.New(input, "testCall")
	p := New(l)
	program := p.ParseProgram()
	p.CheckParserErrors(t)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d",
			1, len(program.Statements))
	}

	prntStmt, ok := program.Statements[0].(*ast.PrintExpressionStatement)
	if !ok {
		t.Fatalf("stmt is not type *ast.PrintExpressionStatement. got=%T",
			program.Statements[0])
	}

	if !CheckLiteralExpression(t, prntStmt.Expression, 12) {
		return
	}

}

// Test Function Call Parsing
func TestCallExpressionParsing(t *testing.T) {
	input := `add(1, 2 * 3, 4 + 5)`

	l := lexer.New(input, "testCall")
	p := New(l)
	program := p.ParseProgram()
	p.CheckParserErrors(t)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d",
			1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("stmt is not type *ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.CallExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not type *ast.CallExpression. got=%T",
			stmt.Expression)
	}

	if !CheckIdentifier(t, exp.Function, "add") {
		return
	}

	if len(exp.Arguments) != 3 {
		t.Fatalf("wrong length of arugments. got=%d",
			len(exp.Arguments))
	}

	CheckLiteralExpression(t, exp.Arguments[0], 1)
	CheckInfixExpression(t, exp.Arguments[1], 2, "*", 3)
	CheckInfixExpression(t, exp.Arguments[2], 4, "+", 5)
}

// Test Function Definition Parameters Parsing
func TestFunctionParameterParsing(t *testing.T) {
	tests := []struct {
		input          string
		expectedParams []string
	}{
		{input: "fn() {}", expectedParams: []string{}},
		{input: "fn(x,) {}", expectedParams: []string{"x"}},
		{input: "fn(x, y, z,) {}", expectedParams: []string{"x", "y", "z"}},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input, "testFunctionParam")
		p := New(l)
		program := p.ParseProgram()
		p.CheckParserErrors(t)

		stmt := program.Statements[0].(*ast.ExpressionStatement)
		function := stmt.Expression.(*ast.FunctionLiteral)

		if len(function.Parameters) != len(tt.expectedParams) {
			t.Errorf("parameters length wrong. expected=%d, got=%d",
				len(tt.expectedParams), len(function.Parameters))
		}

		for i, ident := range tt.expectedParams {
			CheckLiteralExpression(t, function.Parameters[i], ident)
		}
	}
}

// Test Parsing a function
func TestFunctionLiteralParsing(t *testing.T) {
	input := `fn(x, y) { x + y }`

	l := lexer.New(input, "testFunction")
	p := New(l)
	program := p.ParseProgram()
	p.CheckParserErrors(t)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain enough statements. got=%d",
			len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not type *ast.ExpressinStatement. got=%T",
			program.Statements[0])
	}

	function, ok := stmt.Expression.(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("stmt.Expression is not type *ast.FunctionLiteral. got=%T",
			stmt.Expression)
	}

	if len(function.Parameters) != 2 {
		t.Fatalf("function literal paramters wrong. expected 2, got=%d",
			len(function.Parameters))
	}

	CheckLiteralExpression(t, function.Parameters[0], "x")
	CheckLiteralExpression(t, function.Parameters[1], "y")

	if len(function.Body.Statements) != 1 {
		t.Fatalf("function.Body.Statements has not 1 statement. got=%d",
			len(function.Body.Statements))
	}

	bodyStmt, ok := function.Body.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("function.Body.Statements[0] is not type *ast.ExpressionStatement. got=%T",
			function.Body.Statements[0])
	}

	CheckInfixExpression(t, bodyStmt.Expression, "x", "+", "y")
}

// Test If else if else Expression
func TestIfElseIfElseExpression(t *testing.T) {
	input := `
if x < y {
    x
} else if y == z {
	y
} else {
	14
}
`

	l := lexer.New(input, "testIfElseIfElse")
	p := New(l)
	program := p.ParseProgram()
	p.CheckParserErrors(t)

	// Check more
	if len(program.Statements) != 1 {
		t.Errorf("program.Statements does not contain enough statements. got=%d",
			len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not type *ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not type *ast.IfExpression. got=%T",
			stmt.Expression)
	}

	if !CheckInfixExpression(t, exp.Condition, "x", "<", "y") {
		return
	}

	if len(exp.Consequence.Statements) != 1 {
		t.Errorf("consequence is not 1 statements. got=%d",
			len(exp.Consequence.Statements))
	}

	consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statements[0] is not type *ast.ExpressionStatement. got=%T",
			exp.Consequence.Statements[0])
	}

	if !CheckIdentifier(t, consequence.Expression, "x") {
		return
	}

	// Else Statement
	alt, ok := exp.Alternative.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("alternative statement is not type *ast.ExpressionStatement. got=%T",
			exp.Alternative.Statements[0])
	}

	elseStmt, ok := alt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("alternative statement is not type *ast.IfExpression. got=%T",
			alt.Expression)
	}

	if !CheckInfixExpression(t, elseStmt.Condition, "y", "==", "z") {
		return
	}

	consequence, ok = elseStmt.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("elseStatements[0] is not type *ast.ExpressionStatement. got=%T",
			elseStmt.Consequence.Statements[0])
	}

	if !CheckIdentifier(t, consequence.Expression, "y") {
		return
	}

	// Else Statement
	alt, ok = elseStmt.Alternative.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("alternative statement is not type *ast.ExpressionStatement. got=%T",
			exp.Alternative.Statements[0])
	}

	if !CheckIntegerLiteral(t, alt.Expression, 14) {
		return
	}

}

// Test Else If
func TestElseIfExpression(t *testing.T) {
	input := `
if x < y {
    x
} else if x > y {
    y
}
`

	l := lexer.New(input, "testElseIf")
	p := New(l)
	program := p.ParseProgram()
	p.CheckParserErrors(t)

	if len(program.Statements) != 1 {
		t.Fatalf("program does not contain enough statements. got=%d",
			len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not type *ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not type *ast.IfExpression. got=%T",
			stmt.Expression)
	}

	if !CheckInfixExpression(t, exp.Condition, "x", "<", "y") {
		return
	}

	if len(exp.Consequence.Statements) != 1 {
		t.Errorf("consequence is not 1 statements. got=%d",
			len(exp.Consequence.Statements))
	}

	consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statements[0] is not type *ast.ExpressionStatement. got=%T",
			exp.Consequence.Statements[0])
	}

	if !CheckIdentifier(t, consequence.Expression, "x") {
		return
	}

	// Else Statement
	alt, ok := exp.Alternative.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("alternative statement is not type *ast.ExpressionStatement. got=%T",
			exp.Alternative.Statements[0])
	}

	elseStmt, ok := alt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("alternative statement is not type *ast.IfExpression. got=%T",
			alt.Expression)
	}

	if !CheckInfixExpression(t, elseStmt.Condition, "x", ">", "y") {
		return
	}

	consequence, ok = elseStmt.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("elseStatements[0] is not type *ast.ExpressionStatement. got=%T",
			elseStmt.Consequence.Statements[0])
	}

	if !CheckIdentifier(t, consequence.Expression, "y") {
		return
	}

}

// Test the parsing of if else expressions
func TestIfElseExpression(t *testing.T) {
	input := `if x < y { x } else { y }`

	l := lexer.New(input, "testIfElse")
	p := New(l)
	program := p.ParseProgram()
	p.CheckParserErrors(t)

	if len(program.Statements) != 1 {
		t.Fatalf("program does not contain enough statements. got=%d",
			len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not type *ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not type *ast.IfExpression. got=%T",
			stmt.Expression)
	}

	if !CheckInfixExpression(t, exp.Condition, "x", "<", "y") {
		return
	}

	if len(exp.Consequence.Statements) != 1 {
		t.Errorf("consequence is not 1 statements. got=%d",
			len(exp.Consequence.Statements))
	}

	consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statements[0] is not type *ast.ExpressionStatement. got=%T",
			exp.Consequence.Statements[0])
	}

	if !CheckIdentifier(t, consequence.Expression, "x") {
		return
	}

	altern, ok := exp.Alternative.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("alternative statement is not type *ast.ExpressionStatement. got=%T",
			exp.Alternative.Statements[0])
	}

	if !CheckIdentifier(t, altern.Expression, "y") {
		return
	}
}

// Test the parsing of if expressions
func TestIfExpression(t *testing.T) {
	input := `if x < y { x }`

	l := lexer.New(input, "testIf")
	p := New(l)
	program := p.ParseProgram()
	p.CheckParserErrors(t)

	if len(program.Statements) != 1 {
		t.Fatalf("program does not contain enough statements. got=%d",
			len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not type *ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not type *ast.IfExpression. got=%T",
			stmt.Expression)
	}

	if !CheckInfixExpression(t, exp.Condition, "x", "<", "y") {
		return
	}

	if len(exp.Consequence.Statements) != 1 {
		t.Errorf("consequence is not 1 statements. got=%d",
			len(exp.Consequence.Statements))
	}

	consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statements[0] is not type *ast.ExpressionStatement. got=%T",
			exp.Consequence.Statements[0])
	}

	if !CheckIdentifier(t, consequence.Expression, "x") {
		return
	}

	if exp.Alternative != nil {
		t.Errorf("exp.Alternative.Statements was not nil. got=%+v", exp.Alternative)
	}
}

// Check if an identifier has a value
func CheckIdentifier(t *testing.T, exp ast.Expression, value string) bool {
	ident, ok := exp.(*ast.Identifier)
	if !ok {
		t.Errorf("exp not type *ast.Identifier. got=%T",
			exp)
		return false
	}

	if ident.Value != value {
		t.Errorf("ident.Value not %s. got=%s",
			value, ident.Value)
		return false
	}

	if ident.TokenLiteral() != value {
		t.Errorf("ident.TokenLiteral not %s. got=%s",
			value, ident.TokenLiteral())
		return false
	}
	return true
}

// Helper function to check a boolean expression
func CheckBooleanLiteral(t *testing.T, exp ast.Expression, value bool) bool {
	bo, ok := exp.(*ast.Boolean)
	if !ok {
		t.Errorf("exp not *ast.Boolean. got=%T",
			exp)
		return false
	}

	if bo.Value != value {
		t.Errorf("bo.Value not %t. got=%t", value, bo.Value)
		return false
	}

	if bo.TokenLiteral() != fmt.Sprintf("%t", value) {
		t.Errorf("bo.TokenLiteral not %t. got=%s",
			value, bo.TokenLiteral())
		return false
	}

	return false
}

// Check a generic expression for a value
func CheckLiteralExpression(
	t *testing.T,
	exp ast.Expression,
	expected interface{},
) bool {
	switch v := expected.(type) {
	case int:
		return CheckIntegerLiteral(t, exp, float64(v))
	case float64:
		return CheckIntegerLiteral(t, exp, v)
	case string:
		return CheckIdentifier(t, exp, v)
	case bool:
		return CheckBooleanLiteral(t, exp, v)
	}
	t.Errorf("type of exp not handled. got=%T", exp)
	return false
}

// Testing parsing precedence expression
func TestOperatorPrecedenceParsing(t *testing.T) {
	if options.NicerToString {
		return
	}

	// Test Cases
	tests := []struct {
		input    string
		expected string
	}{
		{"1 + (2 + 3) + 4", "[(((1) + ((2) + (3))) + (4))]"},
		{"-a * b", "[((-(a)) * (b))]"},
		{"!+a", "[(!(+(a)))]"},
		{"1 + 1 * 1", "[((1) + ((1) * (1)))]"},
		{
			"a + b * c + d / e - f",
			"[((((a) + ((b) * (c))) + ((d) / (e))) - (f))]",
		},
		{
			"3 + 4; -5 * 5",
			"[((3) + (4)) ((-(5)) * (5))]",
		},
		{
			"5 > 4 == 3 < 4",
			"[(((5) > (4)) == ((3) < (4)))]",
		},
		{
			"5 < 4 != 3 > 4",
			"[(((5) < (4)) != ((3) > (4)))]",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"[(((3) + ((4) * (5))) == (((3) * (1)) + ((4) * (5))))]",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"[(((3) + ((4) * (5))) == (((3) * (1)) + ((4) * (5))))]",
		},
		{
			"a + add(b * c) + d",
			"[(((a) + ((add)(((b) * (c))))) + (d))]",
		},
		{
			"add(a, b, 1, 2 * 3, 4 + 5, add(6, 7 * 8))",
			"[((add)((a), (b), (1), ((2) * (3)), ((4) + (5)), ((add)((6), ((7) * (8))))))]",
		},
		{
			"add(a + b + c * d / f + g)",
			"[((add)(((((a) + (b)) + (((c) * (d)) / (f))) + (g))))]",
		},
		{
			"1 == 2 and 5 > 6",
			"[(((1) == (2)) and ((5) > (6)))]",
		},
		{
			"1 or 5 and 3 > 8",
			"[(((1) or (5)) and ((3) > (8)))]",
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

func CheckInfixExpression(t *testing.T, exp ast.Expression, left interface{},
	operator string, right interface{}) bool {
	opExp, ok := exp.(*ast.InfixExpression)
	if !ok {
		t.Errorf("exp is not type *ast.OperatorExpression. got=%T(%s)",
			exp, exp)
		return false
	}

	if !CheckLiteralExpression(t, opExp.Left, left) {
		return false
	}

	if opExp.Operator != operator {
		t.Errorf("exp.Operator is not %q. got=%q",
			operator, opExp.Operator)
		return false
	}

	if !CheckLiteralExpression(t, opExp.Right, right) {
		return false
	}

	return true
}

// Test parsing of infix expressions
func TestParsingInfixExpressions(t *testing.T) {
	infixTests := []struct {
		input      string
		leftValue  interface{}
		operator   string
		rightValue interface{}
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
		{"true == true", true, "==", true},
		{"false != true", false, "!=", true},
		{"false == false", true, "==", false},
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
			t.Errorf("exp is not type *ast.InfixExpression. got=%T(%v)",
				exp, exp)
		}

		if !CheckLiteralExpression(t, exp.Left, tt.leftValue) {
			return
		}

		if !CheckLiteralExpression(t, exp.Right, tt.rightValue) {
			return
		}

		CheckInfixExpression(t, stmt.Expression, tt.leftValue, tt.operator, tt.rightValue)
	}
}

// Test parsing of prefix
func TestParsingPrefixExpressions(t *testing.T) {
	prefixTests := []struct {
		input        string
		operator     string
		integerValue float64
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

func CheckIntegerLiteral(t *testing.T, il ast.Expression, value float64) bool {
	integ, ok := il.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("expression not type *ast.IntegerLiteral. got=%T", il)
		return false
	}

	if !AlmostEqual(integ.Value, value) {
		t.Errorf("integ.Value not %f. got=%f",
			value, integ.Value)
		return false
	}

	if !strings.Contains(integ.TokenLiteral(), FormatFloat(value)) {
		t.Errorf("integ.TokenLiteral not %s. got=%s",
			FormatFloat(value), integ.TokenLiteral())
		return false
	}

	return true
}

// Test parsing of booleans
func TestBooleanExpression(t *testing.T) {
	input := "true"

	l := lexer.New(input, "testBoolean")
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

	literal, ok := stmt.Expression.(*ast.Boolean)
	if !ok {
		t.Fatalf("expression not type *ast.Boolean. got=%T",
			stmt.Expression)
	}
	if literal.Value != true {
		t.Errorf("literal.Value not %v. got=%v",
			true, literal.Value)
	}
	if literal.TokenLiteral() != "true" {
		t.Errorf("literal.TokenLiteral not %s. got=%s",
			"true", literal.TokenLiteral())
	}
}

// Test parsing of literals
func TestIntegerLiteralExpression(t *testing.T) {
	input := "5"

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
		t.Errorf("literal.Value not %d. got=%f",
			5, literal.Value)
	}
	if literal.TokenLiteral() != "5" {
		t.Errorf("literal.TokenLiteral not %s. got=%s",
			"5", literal.TokenLiteral())
	}
}

// Test parsing of the identifier
func TestIdentifierExpression(t *testing.T) {
	input := "foobar"

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

func TestLetStatementsFully(t *testing.T) {
	tests := []struct {
		input              string
		expectedIdentifier string
		expectedValue      interface{}
	}{
		{"let x = 5;", "x", 5},
		{"let y = true;", "y", true},
		{"let foobar = y;", "foobar", "y"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input, "testLet")
		p := New(l)
		program := p.ParseProgram()
		p.CheckParserErrors(t)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statement. got=%d",
				len(program.Statements))
		}

		stmt := program.Statements[0]
		if !CheckLetStatement(t, stmt, tt.expectedIdentifier) {
			return
		}

		val := stmt.(*ast.LetStatement).Value
		if !CheckLiteralExpression(t, val, tt.expectedValue) {
			return
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
