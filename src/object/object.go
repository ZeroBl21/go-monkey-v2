package object

import "fmt"

type ObjectType string

const (
	INTEGER_OBJ      ObjectType = "INTEGER"
	BOOLEAN_OBJ      ObjectType = "BOOLEAN"
	RETURN_VALUE_OBJ ObjectType = "RETURN_VALUE"
	NULL_OBJ         ObjectType = "NULL"
)

type Object interface {
	Type() ObjectType
	Inspect() string
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

type Boolean struct {
	Value bool
}

func (o *Boolean) Type() ObjectType { return BOOLEAN_OBJ }
func (o *Boolean) Inspect() string  { return fmt.Sprintf("%t", o.Value) }

type ReturnValue struct {
	Value Object
}

func (o *ReturnValue) Type() ObjectType { return RETURN_VALUE_OBJ }
func (o *ReturnValue) Inspect() string  { return o.Value.Inspect() }
