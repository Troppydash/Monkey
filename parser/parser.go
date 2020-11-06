package parser

import (
	"Monkey/ast"
	"Monkey/lexer"
	"Monkey/token"
	"fmt"
	"math"
	"strconv"
)

// TODO: Implement Namespaces
const float64EqualityThreshold = 1e-9

func AlmostEqual(left float64, right float64) bool {
	return math.Abs(left-right) <= float64EqualityThreshold

}

// Expression Parsing Precedences
const (
	_ int = iota
	LOWEST
	ASSIGN  // =
	GATE    // and, or, xor
	EQUAL   // == or !=
	COMPARE // > or < or <= or >=
	SUM     // + or -
	PRODUCT // * or / or %
	PREFIX  // !X or -X
	CALL    // foobar()
	INDEX   // [x]
	DOT     // x.something
)

// A Map Contains a Token to Precedences key value pair
var precedences = map[token.TokenType]int{
	token.XOR:      GATE,
	token.AND:      GATE,
	token.OR:       GATE,
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
	token.PERCENT:  PRODUCT,
	token.LPAREN:   CALL,
	token.LBRACKET: INDEX,
	token.DOT:      DOT,
	token.ASSIGN:   ASSIGN,
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

	done bool
}

// Check if parser have any error
func (p *Parser) HasError() bool {
	return len(p.errors) != 0
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
	p.RegisterPrefix(token.AND, p.ParsePrefixExpression)
	p.RegisterPrefix(token.OR, p.ParsePrefixExpression)
	p.RegisterPrefix(token.XOR, p.ParsePrefixExpression)
	p.RegisterPrefix(token.PERCENT, p.ParsePrefixExpression)
	p.RegisterPrefix(token.ASSIGN, p.ParsePrefixExpression)

	p.RegisterPrefix(token.TRUE, p.ParseBoolean)
	p.RegisterPrefix(token.FALSE, p.ParseBoolean)
	p.RegisterPrefix(token.BREAK, p.ParseBreak)
	p.RegisterPrefix(token.NULL, p.ParseNull)

	p.RegisterPrefix(token.LPAREN, p.parseGroupedExpression)

	p.RegisterPrefix(token.IF, p.ParseIfExpression)
	p.RegisterPrefix(token.FUNCTION, p.ParseFunctionLiteral)
	p.RegisterPrefix(token.HASH, p.ParseHashFunctionLiteral)
	p.RegisterPrefix(token.MODULE, p.ParseModuleExpression)

	p.RegisterPrefix(token.STRING, p.ParseStringLiteral)

	p.RegisterPrefix(token.LBRACKET, p.ParseArrayLiteral)

	p.RegisterPrefix(token.LBRACE, p.ParseHashLiteral)
	//p.RegisterPrefix(token.NEWLINE, p.ParseNewLine)
	//p.RegisterPrefix(token.ASSIGN, p.ParsePrefixExpression)
	p.RegisterPrefix(token.MACRO, p.ParseMacroLiteral)

	// Setup Infix Functions
	p.infixParseFns = make(map[token.TokenType]InfixParseFn)
	p.RegisterInfix(token.ASSIGN, p.ParseInfixExpression)
	p.RegisterInfix(token.DOT, p.ParseInfixExpression)
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
	p.RegisterInfix(token.PERCENT, p.ParseInfixExpression)

	p.RegisterInfix(token.LPAREN, p.ParseCallExpression)

	p.RegisterInfix(token.AND, p.ParseInfixExpression)
	p.RegisterInfix(token.OR, p.ParseInfixExpression)
	p.RegisterInfix(token.XOR, p.ParseInfixExpression)

	p.RegisterInfix(token.LBRACKET, p.ParseIndexExpression)

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
	p.peekToken = <-p.lexer.Tokens
}

// Deprecated
func (p *Parser) PeekPeekToken() token.Token {
	return p.lexer.PeekToken()
}

func (p *Parser) IsPeekEndOfLine() bool {
	return p.PeekTokenIs(token.NEWLINE) || p.PeekTokenIs(token.EOF)
}

func (p *Parser) RemoveNewLines() {
	for p.PeekTokenIs(token.NEWLINE) {
		p.NextToken()
	}
}

// Parse the Whole Program and return an ast tree
func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	// While currentToken is not EOF
	for p.currentToken.Type != token.EOF {
		if p.CurrentTokenIs(token.NEWLINE) {
			p.NextToken()
			continue
		}
		// Parse a statement
		stmt := p.ParseStatement()
		if !p.IsPeekEndOfLine() {
			if !p.CurrentTokenIs(token.NEWLINE) && !p.HasError() {
				// TODO: Warnings
				p.GenerateErrorForToken("No newline after parsing statement", &p.peekToken)
			}
		} else {
			p.NextToken()
		}

		//if p.PeekTokenIs(token.NEWLINE) {
		//	p.NextToken()
		//	fmt.Println("Error")
		//}
		// Add it to the program.Statements if not nil
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		// Advance pointer
		p.NextToken()

	}

	if p.HasError() {
		for _, err := range p.Errors() {
			PrintParserError(err)
			//fmt.Printf("Parser Error: %s, at %d:%d, in file %s\n",
			//	err.Message, err.RowNumber, err.ColumnNumber, err.Filename)
		}
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

	if !p.PeekTokenIs(token.ASSIGN) {
		stmt.Value = &ast.Null{
			Token: p.currentToken,
		}
		return stmt
	}

	p.NextToken()
	//p.RemoveNewLines()
	p.NextToken()

	stmt.Value = p.ParseExpression(LOWEST)

	//if p.PeekTokenIs(token.SEMICOLON) {
	//	p.NextToken()
	//}

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

	stmt.ReturnValue = p.ParseExpression(LOWEST)

	//if p.PeekTokenIs(token.SEMICOLON) {
	//	p.NextToken()
	//}

	return stmt
}

// Parse an expressionStatement
func (p *Parser) ParseExpressionStatement() interface {
	ast.Statement
} {
	// Allocate Memory
	stmt := &ast.ExpressionStatement{Token: p.currentToken}

	// Parse Expression
	exp := p.ParseExpression(LOWEST)

	// Advance through ; if exists
	if p.PeekTokenIs(token.SEMICOLON) {
		pStmt := &ast.PrintExpressionStatement{
			Token:      stmt.Token,
			Expression: exp,
		}
		p.NextToken()
		return pStmt
	}

	stmt.Expression = exp

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

	for {
		if precedence >= p.PeekPrecedence() || p.IsPeekEndOfLine() {
			break
		}
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExpression
		}
		p.NextToken()
		leftExpression = infix(leftExpression)
	}
	//for !p.IsPeekEndOfLine() && precedence < p.PeekPrecedence() {
	//	infix := p.infixParseFns[p.peekToken.Type]
	//	if infix == nil {
	//		return leftExpression
	//	}
	//	p.NextToken()
	//	leftExpression = infix(leftExpression)
	//}
	return leftExpression
}

// Generate Error for a token
func (p *Parser) GenerateErrorForToken(message string, token *token.Token) {
	err := &ParseError{
		Message:      message,
		Filename:     token.Filename,
		RowNumber:    token.RowNumber,
		ColumnNumber: token.ColumnNumber,
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
	ident := &ast.Identifier{
		Token: p.currentToken,
		Value: p.currentToken.Literal,
	}

	//if p.PeekTokenIs(token.ASSIGN) {
	//	p.NextToken()
	//	p.NextToken()
	//	return &ast.AssignmentExpression{
	//		Token: p.currentToken,
	//		Ident: ident,
	//		Value: p.ParseExpression(LOWEST),
	//	}
	//}
	return ident
}

// Parse Literal
func (p *Parser) ParseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{Token: p.currentToken}

	value, err := strconv.ParseFloat(p.currentToken.Literal, 64)
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
	p.RemoveNewLines()
	p.NextToken()
	exp := p.ParseExpression(LOWEST)

	p.RemoveNewLines()
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

	p.RemoveNewLines()

	//if !p.ExpectPeek(token.RPAREN) {
	//	return nil
	//}

	if !p.ExpectPeek(token.LBRACE) {
		return nil
	}

	// Parse then case
	expression.Consequence = p.ParseBlockStatement()

	//if p.PeekPeekToken().Type == token.ELSE {
	//	p.RemoveNewLines()
	//}

	// Parse else case
	if p.PeekTokenIs(token.ELSE) {
		p.NextToken()
		p.RemoveNewLines()

		if p.PeekTokenIs(token.IF) {
			// Parse else if
			p.NextToken()
			p.RemoveNewLines()

			block := &ast.BlockStatement{
				Token: p.peekToken,
				Statements: []ast.Statement{
					&ast.ExpressionStatement{
						Token:      p.peekToken,
						Expression: p.ParseIfExpression(),
					},
				},
			}
			expression.Alternative = block
		} else {
			// Parse else
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
	p.RemoveNewLines()
	p.NextToken()

	for !p.CurrentTokenIs(token.RBRACE) && !p.CurrentTokenIs(token.EOF) {
		stmt := p.ParseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.RemoveNewLines()
		p.NextToken()
	}
	return block
}

func (p *Parser) ParseHashFunctionLiteral() ast.Expression {
	fnLit := &ast.FunctionLiteral{Token: p.currentToken}

	fnLit.Parameters = []*ast.Identifier{}

	p.RemoveNewLines()

	// Check {
	if !p.ExpectPeek(token.LBRACE) {
		return nil
	}

	fnLit.Body = p.ParseBlockStatement()
	return fnLit
}

// ParseModuleExpression parses a module expression
func (p *Parser) ParseModuleExpression() ast.Expression {
	mLit := &ast.ModuleExpression{Token: p.currentToken}

	if !p.ExpectPeek(token.LBRACE) {
		return nil
	}
	mLit.Body = p.ParseBlockStatement()
	return mLit
}

// Parse a function expression
func (p *Parser) ParseFunctionLiteral() ast.Expression {
	fnLit := &ast.FunctionLiteral{Token: p.currentToken}

	// Check (
	if !p.ExpectPeek(token.LPAREN) {
		return nil
	}

	fnLit.Parameters = p.ParseFunctionParameters()

	p.RemoveNewLines()

	// Check {
	if !p.ExpectPeek(token.LBRACE) {
		return nil
	}
	fnLit.Body = p.ParseBlockStatement()
	return fnLit
}

// Parse the parameter list in a function
func (p *Parser) ParseFunctionParameters() []*ast.Identifier {
	var identifiers []*ast.Identifier

	p.RemoveNewLines()
	// If parameter list is empty
	if p.PeekTokenIs(token.RPAREN) {
		p.NextToken()
		return identifiers
	}

	// Advance to the first identifier
	p.NextToken()

	// Parse the first identifier
	ident := &ast.Identifier{
		Token: p.currentToken,
		Value: p.currentToken.Literal,
	}
	identifiers = append(identifiers, ident)

	// Parse the rest of identifiers
	for p.PeekTokenIs(token.COMMA) {
		// Advance to next Identifier
		p.NextToken()
		p.RemoveNewLines()
		// Trailing Comma
		if p.PeekTokenIs(token.RPAREN) {
			break
		}
		p.NextToken()

		ident := &ast.Identifier{
			Token: p.currentToken,
			Value: p.currentToken.Literal,
		}
		identifiers = append(identifiers, ident)
		p.RemoveNewLines()
	}

	// Check )
	if !p.ExpectPeek(token.RPAREN) {
		return nil
	}

	return identifiers
}

// Parse a call expression
func (p *Parser) ParseCallExpression(function ast.Expression) ast.Expression {
	exp := &ast.CallExpression{
		Token:    p.currentToken,
		Function: function,
	}

	if p.CurrentTokenIs(token.LPAREN) {
		exp.Arguments = p.ParseExpressionList(token.RPAREN)
	} else {
		exp.Arguments = []ast.Expression{}
	}

	// So that last parameter can be a function
	for p.PeekTokenIs(token.HASH) || p.PeekTokenIs(token.FUNCTION) {
		if p.PeekTokenIs(token.HASH) {
			p.NextToken()
			exp.Arguments = append(exp.Arguments, p.ParseHashFunctionLiteral())
		} else if p.PeekTokenIs(token.FUNCTION) {
			p.NextToken()
			exp.Arguments = append(exp.Arguments, p.ParseFunctionLiteral())
		}
	}

	return exp
}

// Parse function calling arguments
func (p *Parser) ParseCallArguments() []ast.Expression {
	var args []ast.Expression

	// Empty Arguments
	if p.PeekTokenIs(token.RPAREN) {
		p.NextToken()
		return args
	}

	p.NextToken()
	args = append(args, p.ParseExpression(LOWEST))

	for p.PeekTokenIs(token.COMMA) {
		p.NextToken()
		if p.PeekTokenIs(token.RPAREN) {
			break
		}
		p.NextToken()

		args = append(args, p.ParseExpression(LOWEST))
	}

	if !p.ExpectPeek(token.RPAREN) {
		return nil
	}

	return args
}

// Parse a string expression
func (p *Parser) ParseStringLiteral() ast.Expression {
	return &ast.StringLiteral{
		Token: p.currentToken,
		Value: p.currentToken.Literal,
	}
}

// Parse an array expression
func (p *Parser) ParseArrayLiteral() ast.Expression {
	array := &ast.ArrayLiteral{Token: p.currentToken}

	array.Elements = p.ParseExpressionList(token.RBRACKET)

	return array
}

// Parse a comma separated expression
func (p *Parser) ParseExpressionList(end token.TokenType) []ast.Expression {
	var list []ast.Expression

	p.RemoveNewLines()

	//if p.PeekTokenIs(token.COMMA) {
	//	p.NextToken()
	//}
	if p.PeekTokenIs(end) {
		p.NextToken()
		return list
	}

	p.NextToken()
	list = append(list, p.ParseExpression(LOWEST))

	p.RemoveNewLines()

	for p.PeekTokenIs(token.COMMA) {
		p.NextToken()
		p.RemoveNewLines()
		if p.PeekTokenIs(end) {
			break
		}

		p.NextToken()
		list = append(list, p.ParseExpression(LOWEST))
		p.RemoveNewLines()

	}

	if !p.ExpectPeek(end) {
		return nil
	}

	return list
}

// Parse indexing expression
func (p *Parser) ParseIndexExpression(left ast.Expression) ast.Expression {
	exp := &ast.IndexExpression{Token: p.currentToken, Left: left}

	// Empty First
	if p.PeekTokenIs(token.COLON) {
		exp.Start = nil
	} else {
		p.NextToken()
		exp.Start = p.ParseExpression(LOWEST)
	}

	if p.PeekTokenIs(token.COLON) {
		exp.HasRange = true
		p.NextToken()

		if !p.PeekTokenIs(token.RBRACKET) {
			p.NextToken()
			exp.End = p.ParseExpression(LOWEST)
		}
	}

	if !p.ExpectPeek(token.RBRACKET) {
		return nil
	}

	return exp
}

// Parse HashMaps
func (p *Parser) ParseHashLiteral() ast.Expression {
	hash := &ast.HashLiteral{Token: p.currentToken}
	hash.Pairs = make(map[ast.Expression]ast.Expression)

	for !p.PeekTokenIs(token.RBRACE) {
		p.RemoveNewLines()
		p.NextToken()
		key := p.ParseExpression(LOWEST)

		if !p.ExpectPeek(token.COLON) {
			return nil
		}

		p.NextToken()
		value := p.ParseExpression(LOWEST)

		hash.Pairs[key] = value

		p.RemoveNewLines()
		if !p.PeekTokenIs(token.RBRACE) && !p.ExpectPeek(token.COMMA) {
			return nil
		}
		p.RemoveNewLines()
	}

	if !p.ExpectPeek(token.RBRACE) {
		return nil
	}

	return hash
}

// ParseBreak parse a break token
func (p *Parser) ParseBreak() ast.Expression {
	return &ast.Break{
		Token: p.currentToken,
	}
}

// ParseNull parse a null token
func (p *Parser) ParseNull() ast.Expression {
	return &ast.Null{
		Token: p.currentToken,
	}
}

// ParseMacroLiteral parses a macro definition
func (p *Parser) ParseMacroLiteral() ast.Expression {
	lit := &ast.MacroLiteral{
		Token: p.currentToken,
	}

	if !p.ExpectPeek(token.LPAREN) {
		return nil
	}

	lit.Parameters = p.ParseFunctionParameters()

	if !p.ExpectPeek(token.LBRACE) {
		return nil
	}

	lit.Body = p.ParseBlockStatement()

	return lit
}

// FormatFloat formats floats into a string
func FormatFloat(t float64) string {
	return strconv.FormatFloat(t, 'f', -1, 64)
}
