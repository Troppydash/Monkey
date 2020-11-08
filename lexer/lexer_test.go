package lexer

import (
	"Monkey/token"
	"testing"
)

// Test for Token Parsing
func TestNextToken(t *testing.T) {
	input := `let five = 5
let ten = 10

let add = fn(x, y) {
    x + y;
}

let result = add(five, ten)
!-/*5;
5 < 10 > 5;

if (5 < 10) {
    return true
} else {
    return false
}

10 == 10;
10 != 9;
10 >= 1;
1 <= 10;
"foobar"
"foo bar"
'foobar'
'foo bar'
'foo\n\t\"\':)'
"hello \"world\""
[1, 2]
`

	// What the Parser/Lexer should return
	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
		expectedRow     int64
		expectedColumn  int64
	}{
		{token.LET, "let", 1, 1},
		{token.IDENT, "five", 1, 5},
		{token.ASSIGN, "=", 1, 10},
		{token.INT, "5", 1, 12},
		//{token.SEMICOLON, ";", 1, 13},
		{token.LET, "let", 2, 1},
		{token.IDENT, "ten", 2, 5},
		{token.ASSIGN, "=", 2, 9},
		{token.INT, "10", 2, 11},
		//{token.SEMICOLON, ";", 2, 13},
		{token.LET, "let", 4, 1},
		{token.IDENT, "add", 4, 5},
		{token.ASSIGN, "=", 4, 9},
		{token.FUNCTION, "fn", 4, 11},
		{token.LPAREN, "(", 4, 13},
		{token.IDENT, "x", 4, 14},
		{token.COMMA, ",", 4, 15},
		{token.IDENT, "y", 4, 17},
		{token.RPAREN, ")", 4, 18},
		{token.LBRACE, "{", 4, 20},
		{token.IDENT, "x", 5, 5},
		{token.PLUS, "+", 5, 7},
		{token.IDENT, "y", 5, 9},
		{token.SEMICOLON, ";", 5, 10},
		{token.RBRACE, "}", 6, 1},
		//{token.SEMICOLON, ";", 6, 2},
		{token.LET, "let", 8, 1},
		{token.IDENT, "result", 8, 5},
		{token.ASSIGN, "=", 8, 12},
		{token.IDENT, "add", 8, 14},
		{token.LPAREN, "(", 8, 17},
		{token.IDENT, "five", 8, 18},
		{token.COMMA, ",", 8, 22},
		{token.IDENT, "ten", 8, 24},
		{token.RPAREN, ")", 8, 27},
		//{token.SEMICOLON, ";", 8, 28},
		{token.BANG, "!", 9, 1},
		{token.MINUS, "-", 9, 2},
		{token.SLASH, "/", 9, 3},
		{token.ASTERISK, "*", 9, 4},
		{token.INT, "5", 9, 5},
		{token.SEMICOLON, ";", 9, 6},
		{token.INT, "5", 10, 1},
		{token.LT, "<", 10, 3},
		{token.INT, "10", 10, 5},
		{token.GT, ">", 10, 8},
		{token.INT, "5", 10, 10},
		{token.SEMICOLON, ";", 10, 11},

		{token.IF, "if", 10, 11},
		{token.LPAREN, "(", 10, 11},
		{token.INT, "5", 10, 11},
		{token.LT, "<", 10, 11},
		{token.INT, "10", 10, 11},
		{token.RPAREN, ")", 10, 11},
		{token.LBRACE, "{", 10, 11},

		{token.RETURN, "return", 10, 11},
		{token.TRUE, "true", 10, 11},
		//{token.SEMICOLON, ";", 10, 11},

		{token.RBRACE, "}", 10, 11},
		{token.ELSE, "else", 10, 11},
		{token.LBRACE, "{", 10, 11},
		{token.RETURN, "return", 10, 11},
		{token.FALSE, "false", 10, 11},
		//{token.SEMICOLON, ";", 10, 11},
		{token.RBRACE, "}", 10, 11},

		{token.INT, "10", 10, 11},
		{token.EQ, "==", 10, 11},
		{token.INT, "10", 10, 11},
		{token.SEMICOLON, ";", 10, 11},
		{token.INT, "10", 10, 11},
		{token.NOT_EQ, "!=", 10, 11},
		{token.INT, "9", 10, 11},
		{token.SEMICOLON, ";", 10, 11},

		{token.INT, "10", 10, 11},
		{token.GE, ">=", 10, 11},
		{token.INT, "1", 10, 11},
		{token.SEMICOLON, ";", 10, 11},

		{token.INT, "1", 10, 11},
		{token.LE, "<=", 10, 11},
		{token.INT, "10", 10, 11},
		{token.SEMICOLON, ";", 10, 11},

		{token.STRING, "foobar", 10, 11},
		{token.STRING, "foo bar", 10, 11},

		{token.STRING, "foobar", 10, 11},
		{token.STRING, "foo bar", 10, 11},
		{token.STRING, "foo\n\t\"':)", 10, 11},
		{token.STRING, `hello "world"`, 10, 11},

		{token.LBRACKET, "[", 0, 0},
		{token.INT, "1", 0, 0},
		{token.COMMA, ",", 0, 0},
		{token.INT, "2", 0, 0},
		{token.RBRACKET, "]", 0, 0},

		{token.EOF, "\x00", 10, 12},
	}

	l := New(input, "TestFile")
	for i, tt := range tests {
		tok := l.NextToken()
		for tok.Type == token.NEWLINE {
			tok = l.NextToken()
		}

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}

		if TestLineNumbers {
			if tok.RowNumber != tt.expectedRow {
				t.Fatalf("tests[%d] - rowNumber wrong. expected=%d, got=%d",
					i, tt.expectedRow, tok.RowNumber)
			}

			if tok.ColumnNumber != tt.expectedColumn {
				t.Fatalf("tests[%d] - columnNumber wrong. expected=%d, got=%d",
					i, tt.expectedColumn, tok.ColumnNumber)
			}
		}

	}
}
