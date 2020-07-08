package object

import (
	"Monkey/ast"
	"Monkey/token"
	"fmt"
	"hash/fnv"
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
	STRING_OBJ       = "STRING"       // ""
	BUILTIN_OBJ      = "BUILTIN"      // Builtin Functions
	ARRAY_OBJ        = "ARRAY"        // Arrays
	HASH_OBJ         = "HASH"         // Hashmaps
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

	Hash *HashKey
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

	Hash *HashKey
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

// Error type
type Error struct {
	Message string
	*token.TokenData
}

func (e *Error) Type() ObjectType {
	return ERROR_OBJ
}
func (e *Error) Inspect() string {
	return fmt.Sprintf("Runtime Error: %s, at %d:%d, in file %s\n",
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

// String object
type String struct {
	Value string

	Hash *HashKey
}

func (s *String) Type() ObjectType {
	return STRING_OBJ
}
func (s *String) Inspect() string {
	return s.Value
}

// Builtin function type
type BuiltinFunction func(token token.Token, args ...Object) Object

// Builtin function wrapper
type Builtin struct {
	Fn BuiltinFunction
}

func (b *Builtin) Type() ObjectType {
	return BUILTIN_OBJ
}
func (b *Builtin) Inspect() string {
	return "builtin function"
}

// An Array
type Array struct {
	Elements []Object
}

func (ao *Array) Type() ObjectType {
	return ARRAY_OBJ
}
func (ao *Array) Inspect() string {
	var out strings.Builder

	var elements []string
	for _, e := range ao.Elements {
		elements = append(elements, e.Inspect())
	}

	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")

	return out.String()
}

type Hashable interface {
	HashKey() HashKey
}

// Hash of item
type HashKey struct {
	Type  ObjectType
	Value uint64
}

// Boolean hash
func (b *Boolean) HashKey() HashKey {
	if b.Hash != nil {
		return *b.Hash
	}

	var value uint64

	if b.Value {
		value = 1
	} else {
		value = 2
	}

	hash := HashKey{Type: b.Type(), Value: value}
	b.Hash = &hash

	return hash
}

// Integer hash
func (i *Integer) HashKey() HashKey {
	if i.Hash != nil {
		return *i.Hash
	}

	hash := HashKey{Type: i.Type(), Value: uint64(i.Value)}
	i.Hash = &hash
	return hash
}

// String hash
func (s *String) HashKey() HashKey {
	// Todo: make this into a hashmap caching
	if s.Hash != nil {
		return *s.Hash
	}

	h := fnv.New64a()
	h.Write([]byte(s.Value))

	hash := HashKey{Type: s.Type(), Value: h.Sum64()}
	s.Hash = &hash
	return hash
}

// Store key as well for the Inspect() method
type HashPair struct {
	Key   Object
	Value Object
}

// Hashmap
type Hash struct {
	Pairs map[HashKey]HashPair
}

func (h *Hash) Type() ObjectType {
	return HASH_OBJ
}
func (h *Hash) Inspect() string {
	var out strings.Builder

	var pairs []string
	for _, pair := range h.Pairs {
		pairs = append(pairs, fmt.Sprintf("%s: %s",
			pair.Key.Inspect(), pair.Value.Inspect()))
	}

	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")

	return out.String()
}
