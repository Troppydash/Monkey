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

	// Set Rows and Columns
	if l.ch == '\n' {
		l.currentRow += 1
		l.currentColumn = 1
	} else {
		l.currentColumn += 1
	}

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
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	// Strip all the whitespace until a valid character is find
	l.SkipWhitespace()

	switch l.ch {
	case '=':
		tok = NewToken(token.ASSIGN, l.ch, l.currentRow, l.currentColumn, l.currentFile)
	case ';':
		tok = NewToken(token.SEMICOLON, l.ch, l.currentRow, l.currentColumn, l.currentFile)
	case '(':
		tok = NewToken(token.LPAREN, l.ch, l.currentRow, l.currentColumn, l.currentFile)
	case ')':
		tok = NewToken(token.RPAREN, l.ch, l.currentRow, l.currentColumn, l.currentFile)
	case ',':
		tok = NewToken(token.COMMA, l.ch, l.currentRow, l.currentColumn, l.currentFile)
	case '+':
		tok = NewToken(token.PLUS, l.ch, l.currentRow, l.currentColumn, l.currentFile)
	case '{':
		tok = NewToken(token.LBRACE, l.ch, l.currentRow, l.currentColumn, l.currentFile)
	case '}':
		tok = NewToken(token.RBRACE, l.ch, l.currentRow, l.currentColumn, l.currentFile)
	case 0:
		tok = NewToken(token.EOF, 0, l.currentRow, l.currentColumn, l.currentFile)
	default:
		// If is Text
		if IsLetter(l.ch) {
			// Assign Col/Row/File values
			tok.ColumnNumber = l.currentColumn
			tok.RowNumber = l.currentRow
			tok.Filename = l.currentFile

			// Read the actual text
			tok.Literal = l.ReadIdentifier()

			// Set the type from the text
			tok.Type = token.LookupIdent(tok.Literal)

			return tok
		} else if IsDigit(l.ch) {
			// Assign Col/Row/File values
			tok.ColumnNumber = l.currentColumn
			tok.RowNumber = l.currentRow
			tok.Filename = l.currentFile

			// Set Integer Type and Value
			tok.Type = token.INT
			tok.Literal = l.ReadNumber()
			return tok
		} else {
			// Else return illegal character
			tok = NewToken(token.ILLEGAL, l.ch, l.currentColumn, l.currentColumn, l.currentFile)
		}

	}

	// Advance Pointer
	l.ReadChar()

	return tok
}

// Create a new Token
func NewToken(tokenType token.TokenType, ch rune, rowNumber int64, columnNumber int64, filename string) token.Token {
	return token.Token{
		Type:         tokenType,
		Literal:      string(ch),
		RowNumber:    rowNumber,
		ColumnNumber: columnNumber,
		Filename:     filename,
	}
}

// Read an entire identifier and advance the pointer
func (l *Lexer) ReadIdentifier() string {
	// Cache Starting position
	position := l.position

	// Consume until the next rune is not a valid letter
	for IsLetter(l.ch) {
		l.ReadChar()
	}

	// Return string slice
	return l.input[position:l.position]
}

// If a rune is a valid letter
func IsLetter(ch rune) bool {
	// [a-zA-Z_\?]
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_' || ch == '?'
}

// Eat up all the whitespace
func (l *Lexer) SkipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.ReadChar()
	}
}

func (l *Lexer) ReadNumber() string {
	// Cache Position
	position := l.position

	// Advance until it is not a number
	for IsDigit(l.ch) {
		l.ReadChar()
	}

	// Return string-number slice
	return l.input[position:l.position]
}

// If a rune is a numeric number
func IsDigit(ch rune) bool {
	return '0' <= ch && ch <= '9'
}
