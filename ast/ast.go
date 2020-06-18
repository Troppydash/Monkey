package ast

import "Monkey/token"

// Root Node
type Node interface {
	// Return the underlying string for debug purposes
	TokenLiteral() string
}

// A statement
type Statement interface {
	Node
	StatementNode()
}

// An expression
type Expression interface {
	Node
	ExpressionNode()
}

// A program
type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}
	return ""
}

// A let statement
type LetStatement struct {
	Token token.Token // LET Token
	Name  *Identifier // Variable Name
	Value Expression  // Value
}

func (ls *LetStatement) StatementNode() {

}
func (ls *LetStatement) TokenLiteral() string {
	return ls.Token.Literal
}

// An identifier
type Identifier struct {
	Token token.Token // IDENT Token
	Value string      // Value
}

func (i *Identifier) ExpressionNode() {

}
func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}
