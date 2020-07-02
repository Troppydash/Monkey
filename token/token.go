package token

// TODO: Implement power and sqrt as functions

// All Token Types
const (
	// Others
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	// Identifiers + literals
	IDENT = "IDENT"
	INT   = "INT"

	// Operators
	ASSIGN   = "="
	PLUS     = "+"
	MINUS    = "-"
	BANG     = "!"
	ASTERISK = "*"
	SLASH    = "/"
	PERCENT  = "%"

	// Compare
	LT     = "<"
	GT     = ">"
	EQ     = "=="
	NOT_EQ = "!="
	LE     = "<="
	GE     = ">="

	// Delimiters
	COMMA     = ","
	SEMICOLON = ";"
	COLON     = ":"

	// Specials
	LPAREN = "("
	RPAREN = ")"
	LBRACE = "{"
	RBRACE = "}"

	// Arrays
	LBRACKET = "["
	RBRACKET = "]"

	// Keywords
	FUNCTION = "FUNCTION"
	LET      = "LET"
	TRUE     = "TRUE"
	FALSE    = "FALSE"
	IF       = "IF"
	ELSE     = "ELSE"
	RETURN   = "RETURN"

	AND = "and"
	OR  = "or"
	XOR = "xor"

	STRING = "STRING"
)

// The Type of a Token
type TokenType string

// A Token
type Token struct {
	Type    TokenType
	Literal string

	RowNumber    int64
	ColumnNumber int64
	Filename     string
}

// Error reporting data
type TokenData struct {
	Filename     string
	RowNumber    int64
	ColumnNumber int64
}

func (t *Token) ToTokenData() *TokenData {
	return &TokenData{
		Filename:     t.Filename,
		RowNumber:    t.RowNumber,
		ColumnNumber: t.ColumnNumber,
	}
}

// Hashmap to store string-TokenType values
var keywords = map[string]TokenType{
	"fn":     FUNCTION,
	"let":    LET,
	"true":   TRUE,
	"false":  FALSE,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,

	"and": AND,
	"or":  OR,
	"xor": XOR,
}

// Return a TokenType from a plain string
func LookupIdent(ident string) TokenType {
	// Hashmap lookup
	if tok, ok := keywords[ident]; ok {
		// If found return it
		return tok
	}
	// Else return ident
	return IDENT
}
