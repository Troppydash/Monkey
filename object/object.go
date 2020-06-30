package object

import (
	"Monkey/ast"
	"Monkey/token"
	"fmt"
	"strconv"
	"strings"
)

// Object Types
const (
	INTEGER_OBJ      = "INTEGER"      // Int
	BOOLEAN_OBJ      = "BOOLEAN"      // Bool
	NULL_OBJ         = "NULL"         // Disgusting
	RETURN_VALUE_OBJ = "RETURN_VALUE" // return
	ERROR_OBJ        = "ERROR"        // error
	FUNCTION_OBJ     = "FUNCTION"     // fn
)

// The type of the object
type ObjectType string

// The struct where every object inherent from
type Object interface {
	Type() ObjectType // The type of the object
	Inspect() string  // The ToString method
}

// The integer wrapper
type Integer struct {
	Value float64
}

func (i *Integer) Type() ObjectType {
	return INTEGER_OBJ
}
func (i *Integer) Inspect() string {
	return strconv.FormatFloat(i.Value, 'f', -1, 64)
}

// The boolean wrapper
type Boolean struct {
	Value bool
}

func (b *Boolean) Type() ObjectType {
	return BOOLEAN_OBJ
}
func (b *Boolean) Inspect() string {
	return fmt.Sprintf("%t", b.Value)
}

// Disgusting
type Null struct{}

func (n *Null) Type() ObjectType {
	return NULL_OBJ
}
func (n *Null) Inspect() string {
	return "null"
}

type ReturnValue struct {
	Value Object
}

func (rv *ReturnValue) Type() ObjectType {
	return RETURN_VALUE_OBJ
}
func (rv *ReturnValue) Inspect() string {
	return rv.Value.Inspect()
}

type Error struct {
	Message string
	*token.TokenData
}

func (e *Error) Type() ObjectType {
	return ERROR_OBJ
}
func (e *Error) Inspect() string {
	return fmt.Sprintf("ERROR: %s, at %d:%d, in file %s",
		e.Message, e.RowNumber, e.ColumnNumber, e.Filename)
}

type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
}

func (f *Function) Type() ObjectType {
	return FUNCTION_OBJ
}
func (f *Function) Inspect() string {
	var out strings.Builder

	var params []string
	for _, p := range f.Parameters {
		params = append(params, p.ToString())
	}

	out.WriteString("fn")
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") { ")
	out.WriteString(f.Body.ToString())
	out.WriteString(" }")

	return out.String()
}
