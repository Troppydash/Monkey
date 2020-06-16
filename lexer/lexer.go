package lexer

import "Monkey/token"

// Lexer Struct
type Lexer struct {
	input        string // UTF-8 Input
	position     int    // Current Pointer Position
	readPosition int    // One after Current Pointer Position
	ch           rune   // Current character

	currentRow    int64
	currentColumn int64

	currentFile string
}

// Create a new Lexer Struct
func New(input string, filename string) *Lexer {
	l := &Lexer{
		input:         input,
		currentColumn: 0,
		currentRow:    1,
		currentFile:   filename,
	}
	// Set up pointers
	l.ReadChar()
	return l
}

// Read next Character and advance pointer
func (l *Lexer) ReadChar() {
	// If overflows
	if l.readPosition >= len(l.input) {
		// Set nil
		l.ch = 0
	} else {
		// Else set the character on the current readPosition
		l.ch = []rune(l.input)[l.readPosition]
	}

	// Increase Pointers
	l.position = l.readPosition
	l.readPosition += 1

	// Set Rows and Columns
	if l.ch == '\n' {
		l.currentRow += 1
		l.currentColumn = 1
	} else {
		l.currentColumn += 1
	}
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	switch l.ch {
	case '=':
		tok = newToken(token.ASSIGN, l.ch, l.currentRow, l.currentColumn, l.currentFile)
	case ';':
		tok = newToken(token.SEMICOLON, l.ch, l.currentRow, l.currentColumn, l.currentFile)
	case '(':
		tok = newToken(token.LPAREN, l.ch, l.currentRow, l.currentColumn, l.currentFile)
	case ')':
		tok = newToken(token.RPAREN, l.ch, l.currentRow, l.currentColumn, l.currentFile)
	case ',':
		tok = newToken(token.COMMA, l.ch, l.currentRow, l.currentColumn, l.currentFile)
	case '+':
		tok = newToken(token.PLUS, l.ch, l.currentRow, l.currentColumn, l.currentFile)
	case '{':
		tok = newToken(token.LBRACE, l.ch, l.currentRow, l.currentColumn, l.currentFile)
	case '}':
		tok = newToken(token.RBRACE, l.ch, l.currentRow, l.currentColumn, l.currentFile)
	case 0:
		tok = newToken(token.EOF, 0, l.currentRow, l.currentColumn, l.currentFile)
	}
	l.ReadChar()
	return tok
}

// Create a new Token
func newToken(tokenType token.TokenType, ch rune, rowNumber int64, columnNumber int64, filename string) token.Token {
	return token.Token{
		Type:         tokenType,
		Literal:      string(ch),
		RowNumber:    rowNumber,
		ColumnNumber: columnNumber,
		Filename:     filename,
	}
}
