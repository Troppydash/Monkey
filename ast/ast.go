package ast

import (
	"Monkey/options"
	"Monkey/token"
	"strings"
)

// Root Node
type Node interface {
	// Return the underlying string for debug purposes
	TokenLiteral() string

	// ToString Method
	ToString() string
}

// A statement
type Statement interface {
	Node
	// Statement Identifier
	StatementNode()
}

// An expression
type Expression interface {
	Node
	// Expression Identifier
	ExpressionNode()
}

// A program
type Program struct {
	Statements []Statement // List of statements
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}
	return ""
}
func (p *Program) ToString() string {
	var out strings.Builder

	AddOptionalString(&out, "[")
	for i, s := range p.Statements {
		out.WriteString(s.ToString())
		if i != len(p.Statements)-1 {
			out.WriteString(" ")
		}
	}
	AddOptionalString(&out, "]")

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
	var out strings.Builder

	AddOpeningBrace(&out)
	out.WriteString(ls.TokenLiteral() + " ")
	out.WriteString(ls.Name.ToString())
	out.WriteString(" = ")

	if ls.Value != nil {
		out.WriteString(ls.Value.ToString())
	} else {
		out.WriteString("nil")
	}

	out.WriteString(";")
	AddClosingBrace(&out)

	return out.String()
}

// A return statement
type ReturnStatement struct {
	Token       token.Token // RETURN Token
	ReturnValue Expression  // The return value expression
}

func (rs *ReturnStatement) StatementNode() {}
func (rs *ReturnStatement) TokenLiteral() string {
	return rs.Token.Literal
}
func (rs *ReturnStatement) ToString() string {
	var out strings.Builder

	AddOpeningBrace(&out)
	out.WriteString(rs.TokenLiteral())

	if rs.ReturnValue != nil {
		out.WriteString(" " + rs.ReturnValue.ToString())
	}

	out.WriteString(";")
	AddClosingBrace(&out)
	return out.String()
}

// TODO: Remove this in place of infix
type AssignmentExpression struct {
	Token token.Token
	Ident *Identifier
	Value Expression
}

func (a *AssignmentExpression) ExpressionNode() {}
func (a *AssignmentExpression) TokenLiteral() string {
	return a.Token.Literal
}
func (a *AssignmentExpression) ToString() string {
	var out strings.Builder

	AddOpeningBrace(&out)
	out.WriteString(a.Ident.ToString())
	out.WriteString(a.Value.ToString())
	AddClosingBrace(&out)
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
	var out strings.Builder

	AddOpeningBrace(&out)
	out.WriteString(i.Value)
	AddClosingBrace(&out)
	return out.String()
}

// ExpressionStatement Wrapper that tells it to print
type PrintExpressionStatement struct {
	Token      token.Token
	Expression Expression
}

func (pes *PrintExpressionStatement) StatementNode() {}
func (pes *PrintExpressionStatement) TokenLiteral() string {
	return pes.Token.Literal
}
func (pes *PrintExpressionStatement) ToString() string {
	if pes.Expression != nil {
		return pes.Expression.ToString()
	}
	return ""
}

// An expression statement, a wrapper around an expression
type ExpressionStatement struct {
	Token      token.Token // Expression Token
	Expression Expression  // The wrapped Expression
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
	Token token.Token // INT Token
	Value float64     // Number Value
}

func (il *IntegerLiteral) ExpressionNode() {}
func (il *IntegerLiteral) TokenLiteral() string {
	return il.Token.Literal
}
func (il *IntegerLiteral) ToString() string {
	var out strings.Builder

	AddOpeningBrace(&out)
	out.WriteString(il.Token.Literal)
	AddClosingBrace(&out)
	return out.String()
}

// An expression prefix
type PrefixExpression struct {
	Token    token.Token // Operator Prefix Token
	Operator string      // Operator
	Right    Expression  // Right Expression
}

func (pe *PrefixExpression) ExpressionNode() {}
func (pe *PrefixExpression) TokenLiteral() string {
	return pe.Token.Literal
}
func (pe *PrefixExpression) ToString() string {
	var out strings.Builder

	AddOpeningBrace(&out)
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.ToString())
	AddClosingBrace(&out)
	return out.String()
}

type InfixExpression struct {
	Token    token.Token // Operator Infix Token
	Operator string      // The Operator
	Left     Expression  // Left Expression
	Right    Expression  // Right Expression
}

func (ie *InfixExpression) ExpressionNode() {}
func (ie *InfixExpression) TokenLiteral() string {
	return ie.Token.Literal
}
func (ie *InfixExpression) ToString() string {
	var out strings.Builder

	out.WriteString("(")
	out.WriteString(ie.Left.ToString())
	out.WriteString(" " + ie.Operator + " ")
	out.WriteString(ie.Right.ToString())
	out.WriteString(")")

	return out.String()
}

type Null struct {
	Token token.Token
}

func (n *Null) ExpressionNode() {}
func (n *Null) TokenLiteral() string {
	return n.Token.Literal
}
func (n *Null) ToString() string {
	return n.Token.Literal
}

type Break struct {
	Token token.Token
}

func (b *Break) ExpressionNode() {}
func (b *Break) TokenLiteral() string {
	return b.Token.Literal
}
func (b *Break) ToString() string {
	return b.Token.Literal
}

// Boolean Type
type Boolean struct {
	Token token.Token // Boolean Token
	Value bool        // True or False
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
	Token      token.Token // { Token
	Statements []Statement // List of Statements
}

func (bs *BlockStatement) StatementNode() {}
func (bs *BlockStatement) TokenLiteral() string {
	return bs.Token.Literal
}
func (bs *BlockStatement) ToString() string {
	var out strings.Builder

	out.WriteString("{")
	for _, s := range bs.Statements {
		out.WriteString(s.ToString())
		out.WriteString("")
	}
	out.WriteString("}")
	return out.String()
}

// An if(else) expression
type IfExpression struct {
	Token       token.Token     // IF Token
	Condition   Expression      // If Condition
	Consequence *BlockStatement // True Block
	Alternative *BlockStatement // False Block
}

func (ie *IfExpression) ExpressionNode() {}
func (ie *IfExpression) TokenLiteral() string {
	return ie.Token.Literal
}
func (ie *IfExpression) ToString() string {
	var out strings.Builder

	AddOpeningBrace(&out)
	out.WriteString("if ")
	out.WriteString(ie.Condition.ToString())
	out.WriteString(" ")
	AddOpeningBrace(&out)
	out.WriteString(ie.Consequence.ToString())
	AddClosingBrace(&out)
	if ie.Alternative != nil {
		out.WriteString(" else ")
		out.WriteString(ie.Alternative.ToString())
	}
	AddClosingBrace(&out)
	return out.String()
}

// A module expression
type ModuleExpression struct {
	Token token.Token
	Body  *BlockStatement
}

func (ms *ModuleExpression) ExpressionNode() {}
func (ms *ModuleExpression) TokenLiteral() string {
	return ms.Token.Literal
}
func (ms *ModuleExpression) ToString() string {
	return "module " + ms.Body.ToString()
}

// Function Definition Expression
type FunctionLiteral struct {
	Token      token.Token     // Fn Token
	Parameters []*Identifier   // List of Parameters
	Body       *BlockStatement // The Function Body
}

func (fl *FunctionLiteral) ExpressionNode() {}
func (fl *FunctionLiteral) TokenLiteral() string {
	return fl.Token.Literal
}
func (fl *FunctionLiteral) ToString() string {
	var out strings.Builder

	// Build Param list
	var params []string
	for _, p := range fl.Parameters {
		params = append(params, p.ToString())
	}

	AddOpeningBrace(&out)

	out.WriteString(fl.TokenLiteral())
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") ")
	out.WriteString(fl.Body.ToString())

	AddClosingBrace(&out)

	return out.String()
}

// Function Call Expression
type CallExpression struct {
	Token     token.Token  // ( Token
	Function  Expression   // The target Function
	Arguments []Expression // Argument list
}

func (ce *CallExpression) ExpressionNode() {}
func (ce *CallExpression) TokenLiteral() string {
	return ce.Token.Literal
}
func (ce *CallExpression) ToString() string {
	var out strings.Builder

	var args []string
	for _, a := range ce.Arguments {
		args = append(args, a.ToString())
	}

	AddOpeningBrace(&out)

	out.WriteString(ce.Function.ToString())
	out.WriteString("(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")

	AddClosingBrace(&out)
	return out.String()
}

// String
type StringLiteral struct {
	Token token.Token
	Value string
}

func (sl *StringLiteral) ExpressionNode() {}
func (sl *StringLiteral) TokenLiteral() string {
	return sl.Token.Literal
}
func (sl *StringLiteral) ToString() string {
	return sl.Value
}

// Array object
type ArrayLiteral struct {
	Token    token.Token
	Elements []Expression
}

func (al *ArrayLiteral) ExpressionNode() {}
func (al *ArrayLiteral) TokenLiteral() string {
	return al.Token.Literal
}
func (al *ArrayLiteral) ToString() string {
	var out strings.Builder

	var elements []string
	for _, el := range al.Elements {
		elements = append(elements, el.ToString())
	}

	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")

	return out.String()
}

// Indexing Array
// [1:2]
// [1:]
// [:2]
// [:]
// [1]
type IndexExpression struct {
	Token token.Token
	Left  Expression

	Start    Expression
	End      Expression
	HasRange bool
}

func (ie *IndexExpression) ExpressionNode() {}
func (ie *IndexExpression) TokenLiteral() string {
	return ie.Token.Literal
}
func (ie *IndexExpression) ToString() string {
	var out strings.Builder

	AddOpeningBrace(&out)
	out.WriteString(ie.Left.ToString())
	out.WriteString("[")
	if ie.Start != nil {
		out.WriteString(ie.Start.ToString())
	}
	if ie.HasRange {
		out.WriteString(":")
	}
	if ie.End != nil {
		out.WriteString(ie.End.ToString())
	}
	out.WriteString("]")
	AddClosingBrace(&out)

	return out.String()
}

type HashLiteral struct {
	Token token.Token
	Pairs map[Expression]Expression
}

func (hl *HashLiteral) ExpressionNode() {}
func (hl *HashLiteral) TokenLiteral() string {
	return hl.Token.Literal
}
func (hl *HashLiteral) ToString() string {
	var out strings.Builder

	var pairs []string
	for key, value := range hl.Pairs {
		pairs = append(pairs, key.ToString()+":"+value.ToString())
	}

	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")

	return out.String()
}

// MacroLiteral represent a macro def
type MacroLiteral struct {
	Token      token.Token
	Parameters []*Identifier
	Body       *BlockStatement
}

func (ml *MacroLiteral) ExpressionNode() {}
func (ml *MacroLiteral) TokenLiteral() string {
	return ml.Token.Literal
}
func (ml *MacroLiteral) ToString() string {
	var out strings.Builder

	var params []string
	for _, p := range ml.Parameters {
		params = append(params, p.ToString())
	}

	out.WriteString(ml.TokenLiteral())
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(")")
	out.WriteString(ml.Body.ToString())

	return out.String()
}

func AddOptionalString(out *strings.Builder, str string) {
	if !options.NicerToString {
		out.WriteString(str)
	}
}

func AddOpeningBrace(out *strings.Builder) {
	AddOptionalString(out, "(")
}

func AddClosingBrace(out *strings.Builder) {
	AddOptionalString(out, ")")
}
