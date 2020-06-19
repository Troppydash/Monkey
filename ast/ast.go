package ast

import (
	"Monkey/token"
	"bytes"
)

// Root Node
type Node interface {
	// Return the underlying string for debug purposes
	TokenLiteral() string
	ToString() string
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

func (p *Program) ToString() string {
	var out bytes.Buffer

	for _, s := range p.Statements {
		out.WriteString(s.ToString())
	}

	return out.String()
}

// A let statement
type LetStatement struct {
	Token token.Token // LET Token
	Name  *Identifier // Variable Name
	Value Expression  // Value
}

func (ls *LetStatement) StatementNode() {}
func (ls *LetStatement) TokenLiteral() string {
	return ls.Token.Literal
}
func (ls *LetStatement) ToString() string {
	var out bytes.Buffer

	out.WriteString(ls.TokenLiteral() + " ")
	out.WriteString(ls.Name.ToString())
	out.WriteString(" = ")

	if ls.Value != nil {
		out.WriteString(ls.Value.ToString())
	} else {
		out.WriteString("nil")
	}

	out.WriteString(";")

	return out.String()
}

// A return statement
type ReturnStatement struct {
	Token       token.Token
	ReturnValue Expression
}

func (rs *ReturnStatement) StatementNode() {}
func (rs *ReturnStatement) TokenLiteral() string {
	return rs.Token.Literal
}
func (rs *ReturnStatement) ToString() string {
	var out bytes.Buffer

	out.WriteString(rs.TokenLiteral())

	if rs.ReturnValue != nil {
		out.WriteString(" " + rs.ReturnValue.ToString())
	}

	out.WriteString(";")
	return out.String()
}

// An identifier
type Identifier struct {
	Token token.Token // IDENT Token
	Value string      // Value
}

func (i *Identifier) ExpressionNode() {}
func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}
func (i *Identifier) ToString() string {
	return i.Value
}

// An expression statement, a wrapper around an expression
type ExpressionStatement struct {
	Token      token.Token
	Expression Expression
}

func (es *ExpressionStatement) StatementNode() {}
func (es *ExpressionStatement) TokenLiteral() string {
	return es.Token.Literal
}
func (es *ExpressionStatement) ToString() string {
	if es.Expression != nil {
		return es.Expression.ToString()
	}
	return ""
}
