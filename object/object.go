package object

import (
	"fmt"
	"strconv"
)

// Object Types
const (
	INTEGER_OBJ = "INTEGER" // Int
	BOOLEAN_OBJ = "BOOLEAN" // Bool
	NULL_OBJ    = "NULL"    // Disgusting
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
