package parser

import (
	"Monkey/ast"
	"Monkey/lexer"
	"Monkey/token"
	"fmt"
	"strconv"
)

// Expression Parsing Precedences
const (
	_ int = iota
	LOWEST
	COMPARE // == or > or < or <= or >= or !=
	SUM     // + or -
	PRODUCT // * or /
	PREFIX  // !X or -X
	CALL    // foobar()
)

// An error struct
type ParseError struct {
	Message string

	Filename     string
	RowNumber    int64
	ColumnNumber int64
}

// The Parser struct
type Parser struct {
	lexer *lexer.Lexer

	currentToken token.Token
	peekToken    token.Token

	// Pratt Maps
	prefixParseFns map[token.TokenType]PrefixParseFn
	infixParseFns  map[token.TokenType]InfixParseFn

	// This is an array of pointers
	errors []*ParseError
}

// Construct a new Parser
func New(l *lexer.Lexer) *Parser {
	p := &Parser{lexer: l, errors: []*ParseError{}}

	// Setup both currentToken and peekToken
	p.NextToken()
	p.NextToken()

	// Setup Pratt Parsing Functions
	p.prefixParseFns = make(map[token.TokenType]PrefixParseFn)
	p.RegisterPrefix(token.IDENT, p.ParseIdentifier)
	p.RegisterPrefix(token.INT, p.ParseIntegerLiteral)

	return p
}

// Register a prefix fn by adding it to the hashmap
func (p *Parser) RegisterPrefix(tokenType token.TokenType, fn PrefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

// Register an infix fn by adding it to the hashmap
func (p *Parser) RegisterInfix(tokenType token.TokenType, fn InfixParseFn) {
	p.infixParseFns[tokenType] = fn
}

// Advance the pointer by reading the next token from the lexer
func (p *Parser) NextToken() {
	p.currentToken = p.peekToken
	p.peekToken = p.lexer.NextToken()
}

// Parse the Whole Program and return an ast tree
func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	// While currentToken is not EOF
	for p.currentToken.Type != token.EOF {
		// Parse a statement
		stmt := p.ParseStatement()
		// Add it to the program.Statements if not nil
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		// Advance pointer
		p.NextToken()
	}

	return program
}

// Parse a sentence
func (p *Parser) ParseStatement() ast.Statement {
	// Switch on the token type
	switch p.currentToken.Type {
	case token.LET:
		// Hand it over to parse let
		return p.ParseLetStatement()
	case token.RETURN:
		// Hand it over to parse return
		return p.ParseReturnStatement()
	default:
		// Hand it over to parse expression
		return p.ParseExpressionStatement()
	}
}

// Parse a let sentence
func (p *Parser) ParseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{
		Token: p.currentToken,
	}

	// If next token is NOT an identifier
	if !p.ExpectPeek(token.IDENT) {
		return nil
	}

	stmt.Name = &ast.Identifier{
		Token: p.currentToken,
		Value: p.currentToken.Literal,
	}

	if !p.ExpectPeek(token.ASSIGN) {
		return nil
	}

	// TODO: we are skipping the expression parsing for now
	for !p.CurrentTokenIs(token.SEMICOLON) {
		p.NextToken()
	}

	return stmt
}

// Check if current token is a certain type
func (p *Parser) CurrentTokenIs(t token.TokenType) bool {
	return p.currentToken.Type == t
}

// Check if the next token is a certain type
func (p *Parser) PeekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

// Peek and advance if type matches
func (p *Parser) ExpectPeek(t token.TokenType) bool {
	if p.PeekTokenIs(t) {
		p.NextToken()
		return true
	} else {
		p.PeekError(t)
		return false
	}
}

// Returns all the errors in the parser
func (p *Parser) Errors() []*ParseError {
	return p.errors
}

// Add a peek error to the parser
func (p *Parser) PeekError(t token.TokenType) {
	message := fmt.Sprintf("expected next token to be %s, got %s instead",
		t, p.peekToken.Type)

	p.GenerateErrorForToken(message, &p.peekToken)
}

// Parse a return statement
func (p *Parser) ParseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.currentToken}

	p.NextToken()

	// TODO: Skipping the expressions for now
	for !p.CurrentTokenIs(token.SEMICOLON) {
		p.NextToken()
	}

	return stmt
}

// Parse a expressionStatement
func (p *Parser) ParseExpressionStatement() *ast.ExpressionStatement {
	// Allocate Memory
	stmt := &ast.ExpressionStatement{Token: p.currentToken}

	// Parse Expression
	stmt.Expression = p.ParseExpression(LOWEST)

	// Advance through ; if exists
	if p.PeekTokenIs(token.SEMICOLON) {
		// TODO: Print if ;
		p.NextToken()
	}

	return stmt
}

// Parse an expression
func (p *Parser) ParseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.currentToken.Type]
	if prefix == nil {
		return nil
	}
	leftExpression := prefix()
	return leftExpression
}

// Generate Error for a token
func (p *Parser) GenerateErrorForToken(message string, token *token.Token) {
	err := &ParseError{
		Message:      message,
		Filename:     token.Filename,
		RowNumber:    token.RowNumber,
		ColumnNumber: token.RowNumber,
	}
	p.errors = append(p.errors, err)
}

// Pratt Parser Function Types
type (
	PrefixParseFn func() ast.Expression
	InfixParseFn  func(ast.Expression) ast.Expression
)

// Pratt Parser Functions //

// Parse Identifier
func (p *Parser) ParseIdentifier() ast.Expression {
	return &ast.Identifier{
		Token: p.currentToken,
		Value: p.currentToken.Literal,
	}
}

// Parse Literal
func (p *Parser) ParseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{Token: p.currentToken}

	value, err := strconv.ParseInt(p.currentToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.currentToken.Literal)
		p.GenerateErrorForToken(msg, &p.currentToken)
		return nil
	}

	lit.Value = value
	return lit
}
