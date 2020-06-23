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
	EQUAL   // == or !=
	COMPARE // > or < or <= or >=
	SUM     // + or -
	PRODUCT // * or /
	PREFIX  // !X or -X
	CALL    // foobar()
)

// A Map Contains a Token to Precedences key value pair
var precedences = map[token.TokenType]int{
	token.EQ:       EQUAL,
	token.NOT_EQ:   EQUAL,
	token.LE:       COMPARE,
	token.GE:       COMPARE,
	token.GT:       COMPARE,
	token.LT:       COMPARE,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.SLASH:    PRODUCT,
	token.ASTERISK: PRODUCT,
}

// Peek the precedence of the next token in the parser
func (p *Parser) PeekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}
	return LOWEST
}

// Gets the precedences of the current token in the parser
func (p *Parser) CurrentPrecedence() int {
	if p, ok := precedences[p.currentToken.Type]; ok {
		return p
	}
	return LOWEST
}

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

	// Setup Prefix Functions
	p.prefixParseFns = make(map[token.TokenType]PrefixParseFn)
	p.RegisterPrefix(token.IDENT, p.ParseIdentifier)
	p.RegisterPrefix(token.INT, p.ParseIntegerLiteral)
	p.RegisterPrefix(token.BANG, p.ParsePrefixExpression)
	p.RegisterPrefix(token.MINUS, p.ParsePrefixExpression)
	p.RegisterPrefix(token.PLUS, p.ParsePrefixExpression)

	p.RegisterPrefix(token.TRUE, p.ParseBoolean)
	p.RegisterPrefix(token.FALSE, p.ParseBoolean)

	p.RegisterPrefix(token.LPAREN, p.parseGroupedExpression)

	p.RegisterPrefix(token.IF, p.ParseIfExpression)

	// Setup Infix Functions
	p.infixParseFns = make(map[token.TokenType]InfixParseFn)
	p.RegisterInfix(token.PLUS, p.ParseInfixExpression)
	p.RegisterInfix(token.MINUS, p.ParseInfixExpression)
	p.RegisterInfix(token.SLASH, p.ParseInfixExpression)
	p.RegisterInfix(token.ASTERISK, p.ParseInfixExpression)
	p.RegisterInfix(token.EQ, p.ParseInfixExpression)
	p.RegisterInfix(token.NOT_EQ, p.ParseInfixExpression)
	p.RegisterInfix(token.GT, p.ParseInfixExpression)
	p.RegisterInfix(token.LT, p.ParseInfixExpression)
	p.RegisterInfix(token.GE, p.ParseInfixExpression)
	p.RegisterInfix(token.LE, p.ParseInfixExpression)

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
		p.NoPrefixParseFnError(p.currentToken)
		return nil
	}
	leftExpression := prefix()

	for !p.PeekTokenIs(token.SEMICOLON) && precedence < p.PeekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExpression
		}

		p.NextToken()
		leftExpression = infix(leftExpression)
	}
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

// Error when can't find any available prefix parse functions
func (p *Parser) NoPrefixParseFnError(token token.Token) {
	msg := fmt.Sprintf("no prefix parse function for %s found", token.Type)
	p.GenerateErrorForToken(msg, &token)
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

// Parses Prefix
func (p *Parser) ParsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{
		Token:    p.currentToken,
		Operator: p.currentToken.Literal,
	}

	p.NextToken()

	expression.Right = p.ParseExpression(PREFIX)

	return expression
}

// Parses Infix
func (p *Parser) ParseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Token:    p.currentToken,
		Operator: p.currentToken.Literal,
		Left:     left,
	}

	// Gets the Precedence of the operator
	prec := p.CurrentPrecedence()
	// Advance to the right expression
	p.NextToken()
	expression.Right = p.ParseExpression(prec)

	return expression
}

// Parse Boolean
func (p *Parser) ParseBoolean() ast.Expression {
	return &ast.Boolean{
		Token: p.currentToken,
		Value: p.CurrentTokenIs(token.TRUE),
	}
}

// Parse grouped expression
func (p *Parser) parseGroupedExpression() ast.Expression {
	p.NextToken()
	exp := p.ParseExpression(LOWEST)

	if !p.ExpectPeek(token.RPAREN) {
		return nil
	}

	return exp
}

// Parse if expression
func (p *Parser) ParseIfExpression() ast.Expression {
	expression := &ast.IfExpression{Token: p.currentToken}

	//if !p.ExpectPeek(token.LPAREN) {
	//	return nil
	//}

	// Advance pass the if token
	p.NextToken()

	expression.Condition = p.ParseExpression(LOWEST)

	//if !p.ExpectPeek(token.RPAREN) {
	//	return nil
	//}

	if !p.ExpectPeek(token.LBRACE) {
		return nil
	}

	// Parse then case
	expression.Consequence = p.ParseBlockStatement()

	// Parse else case
	if p.PeekTokenIs(token.ELSE) {
		p.NextToken()

		if p.PeekTokenIs(token.IF) {
			p.NextToken()
			block := &ast.BlockStatement{
				Token:      p.peekToken,
				Statements: []ast.Statement{},
			}
			exp := p.ParseIfExpression()
			stmt := &ast.ExpressionStatement{
				Token:      p.peekToken,
				Expression: exp,
			}
			block.Statements = append(block.Statements, stmt)
			expression.Alternative = block
		} else {
			if !p.ExpectPeek(token.LBRACE) {
				return nil
			}
			expression.Alternative = p.ParseBlockStatement()
		}

	}

	return expression
}

// Parse a block statement
func (p *Parser) ParseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{Token: p.currentToken}
	block.Statements = []ast.Statement{}

	// Advance pass the curly brace token
	p.NextToken()

	for !p.CurrentTokenIs(token.RBRACE) && !p.CurrentTokenIs(token.EOF) {
		stmt := p.ParseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.NextToken()
	}
	return block
}
