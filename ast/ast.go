package ast

import (
	"Monkey/token"
	"bytes"
	"strings"
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

// An expression Literal
type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (il *IntegerLiteral) ExpressionNode() {}
func (il *IntegerLiteral) TokenLiteral() string {
	return il.Token.Literal
}
func (il *IntegerLiteral) ToString() string {
	return il.Token.Literal
}

// An expression prefix
type PrefixExpression struct {
	Token    token.Token
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) ExpressionNode() {}
func (pe *PrefixExpression) TokenLiteral() string {
	return pe.Token.Literal
}
func (pe *PrefixExpression) ToString() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.ToString())
	out.WriteString(")")

	return out.String()
}

type InfixExpression struct {
	Token    token.Token
	Operator string
	Left     Expression
	Right    Expression
}

func (ie *InfixExpression) ExpressionNode() {}
func (ie *InfixExpression) TokenLiteral() string {
	return ie.Token.Literal
}
func (ie *InfixExpression) ToString() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(ie.Left.ToString())
	out.WriteString(" " + ie.Operator + " ")
	out.WriteString(ie.Right.ToString())
	out.WriteString(")")

	return out.String()
}

// Boolean Type
type Boolean struct {
	Token token.Token
	Value bool
}

func (b *Boolean) ExpressionNode() {}
func (b *Boolean) TokenLiteral() string {
	return b.Token.Literal
}
func (b *Boolean) ToString() string {
	return b.Token.Literal
}

// A Group/Block of statements
type BlockStatement struct {
	Token      token.Token
	Statements []Statement
}

func (bs *BlockStatement) StatementNode() {}
func (bs *BlockStatement) TokenLiteral() string {
	return bs.Token.Literal
}
func (bs *BlockStatement) ToString() string {
	var out strings.Builder
	for _, s := range bs.Statements {
		out.WriteString(s.ToString())
	}
	return out.String()
}

// An if(else) expression
type IfExpression struct {
	Token       token.Token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (ie *IfExpression) ExpressionNode() {}
func (ie *IfExpression) TokenLiteral() string {
	return ie.Token.Literal
}
func (ie *IfExpression) ToString() string {
	// TODO: Replace all occurence of this
	var out strings.Builder

	out.WriteString("if")
	out.WriteString(ie.Condition.ToString())
	out.WriteString(" ")
	out.WriteString(ie.Consequence.ToString())

	if ie.Alternative != nil {
		out.WriteString("else ")
		out.WriteString(ie.Alternative.ToString())
	}

	return out.String()
}
