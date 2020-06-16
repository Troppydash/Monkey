package lexer

import (
	"Monkey/token"
	"testing"
)

// TODO: Change the test
// Test for Token Parsing
func TestNextToken(t *testing.T) {
	input := `=+(){},;`

	// What the Parser/Lexer should return
	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
		expectedRow     int64
		expectedColumn  int64
	}{
		{token.ASSIGN, "=", 1, 1},
		{token.PLUS, "+", 1, 2},
		{token.LPAREN, "(", 1, 3},
		{token.RPAREN, ")", 1, 4},
		{token.LBRACE, "{", 1, 5},
		{token.RBRACE, "}", 1, 6},
		{token.COMMA, ",", 1, 7},
		{token.SEMICOLON, ";", 1, 8},
	}

	l := New(input, "TestFile")
	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}

		if tok.RowNumber != tt.expectedRow {
			t.Fatalf("tests[%d] - rownumber wrong. expected=%q, got=%q",
				i, tt.expectedRow, tok.RowNumber)
		}

		if tok.ColumnNumber != tt.expectedColumn {
			t.Fatalf("tests[%d] - columnnumber wrong. expected=%q, got=%q",
				i, tt.expectedColumn, tok.ColumnNumber)
		}
	}
}
