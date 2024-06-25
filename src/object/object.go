package object

import (
	"bytes"
	"fmt"
	"hash/fnv"
	"strings"

	"github.com/ZeroBl21/go-monkey/src/ast"
)

type ObjectType string

const (
	INTEGER_OBJ      ObjectType = "INTEGER"
	STRING_OBJ       ObjectType = "STRING"
	BOOLEAN_OBJ      ObjectType = "BOOLEAN"
	RETURN_VALUE_OBJ ObjectType = "RETURN_VALUE"
	NULL_OBJ         ObjectType = "NULL"
	ERROR_OBJ        ObjectType = "ERROR"
	FUNCTION_OBJ     ObjectType = "FUNCTION"
	BUILTIN_OBJ      ObjectType = "BUILTIN"
	ARRAY_OBJ        ObjectType = "ARRAY"
	HASH_OBJ         ObjectType = "HASH"
)

type Object interface {
	Type() ObjectType
	Inspect() string
}

type HashKey struct {
	Type  ObjectType
	Value uint64
}

type Hashable interface {
	HashKey() HashKey
}

type Null struct {
	Value bool
}

func (o *Null) Type() ObjectType { return NULL_OBJ }
func (o *Null) Inspect() string  { return "null" }

type Integer struct {
	Value int64
}

func (o *Integer) Type() ObjectType { return INTEGER_OBJ }
func (o *Integer) Inspect() string  { return fmt.Sprintf("%d", o.Value) }
func (o *Integer) HashKey() HashKey {
	return HashKey{
		Type:  o.Type(),
		Value: uint64(o.Value),
	}
}

type String struct {
	Value string
}

func (o *String) Type() ObjectType { return STRING_OBJ }
func (o *String) Inspect() string  { return o.Value }
func (o *String) HashKey() HashKey {
	h := fnv.New64a()
	h.Write([]byte(o.Value))

	return HashKey{
		Type:  o.Type(),
		Value: h.Sum64(),
	}
}

type Boolean struct {
	Value bool
}

func (o *Boolean) Type() ObjectType { return BOOLEAN_OBJ }
func (o *Boolean) Inspect() string  { return fmt.Sprintf("%t", o.Value) }
func (o *Boolean) HashKey() HashKey {
	var value uint64

	if o.Value {
		value = 1
	} else {
		value = 0
	}

	return HashKey{
		Type:  o.Type(),
		Value: value,
	}
}

type ReturnValue struct {
	Value Object
}

func (o *ReturnValue) Type() ObjectType { return RETURN_VALUE_OBJ }
func (o *ReturnValue) Inspect() string  { return o.Value.Inspect() }

type Error struct {
	Message string
}

func (o *Error) Type() ObjectType { return ERROR_OBJ }
func (o *Error) Inspect() string  { return "ERROR: " + o.Message }

type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
}

func (f *Function) Type() ObjectType { return FUNCTION_OBJ }
func (f *Function) Inspect() string {
	var out bytes.Buffer

	params := []string{}
	for _, p := range f.Parameters {
		params = append(params, p.String())
	}

	out.WriteString("fn")
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(")\n")
	out.WriteString(f.Body.String())
	out.WriteString("\n}")

	return out.String()
}

type BuiltinFunction func(args ...Object) Object

type Builtin struct {
	Fn BuiltinFunction
}

func (o *Builtin) Type() ObjectType { return BUILTIN_OBJ }
func (o *Builtin) Inspect() string  { return "builtin function" }

type Array struct {
	Elements []Object
}

func (o *Array) Type() ObjectType { return ARRAY_OBJ }
func (o *Array) Inspect() string {
	var out bytes.Buffer

	elements := []string{}
	for _, element := range o.Elements {
		elements = append(elements, element.Inspect())
	}

	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")

	return out.String()
}

type HashPair struct {
	Key   Object
	Value Object
}

type Hash struct {
	Pairs map[HashKey]HashPair
}

func (h *Hash) Type() ObjectType { return HASH_OBJ }
func (h *Hash) Inspect() string {
	var out bytes.Buffer

	pairs := []string{}
	for _, pair := range h.Pairs {
		pairs = append(pairs, fmt.Sprintf("%s: %s",
			pair.Key.Inspect(), pair.Value.Inspect()))
	}

	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")

	return out.String()
}
