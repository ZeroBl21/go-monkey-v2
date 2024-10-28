package object

import (
	"fmt"
	"unicode/utf8"
)

type BuiltinsFns struct {
	Name    string
	Builtin *Builtin
}

// The built-in functions / standard-library methods are stored here.
var builtins = []BuiltinsFns{}

// init registers built-in functions to the "standard library" map.
func init() {
	RegisterBuiltin("len", _lenFn)
	RegisterBuiltin("unicodeLen", _unicodeLenFn)
	RegisterBuiltin("first", _firstFn)
	RegisterBuiltin("last", _lastFn)
	RegisterBuiltin("rest", _restFn)
	RegisterBuiltin("push", _pushFn)
	RegisterBuiltin("puts", _putsFn)
}

// Utils

// RegisterBuiltin registers a built-in function.
func RegisterBuiltin(name string, fn BuiltinFunction) {
	builtins = append(builtins, BuiltinsFns{
		Name:    name,
		Builtin: &Builtin{Fn: fn},
	})
}

func GetBuiltinByName(name string) *Builtin {
	for _, def := range builtins {
		if def.Name == name {
			return def.Builtin
		}
	}

	return nil
}

// Utility function to create a new error object.
func newError(format string, a ...any) *Error {
	return &Error{
		Message: fmt.Sprintf(format, a...),
	}
}

// Utility function to check if an object is an error.
func isError(obj Object) bool {
	if obj != nil {
		return obj.Type() == ERROR_OBJ
	}
	return false
}

// Internal Functions

// length of item in runes
func _lenFn(args ...Object) Object {
	if len(args) != 1 {
		return newError("wrong number of arguments. got=%d, want=1", len(args))
	}
	switch arg := args[0].(type) {
	case *String:
		return &Integer{Value: int64(utf8.RuneCountInString(arg.Value))}
	case *Array:
		return &Integer{Value: int64(len(arg.Elements))}
	default:
		return newError("argument to `len` not supported, got=%s", args[0].Type())
	}
}

// length of item but counting bytes individually
func _unicodeLenFn(args ...Object) Object {
	if len(args) != 1 {
		return newError("wrong number of arguments. got=%d, want=1", len(args))
	}
	switch arg := args[0].(type) {
	case *String:
		return &Integer{Value: int64(len(arg.Value))}
	default:
		return newError("argument to `unicodeLen` not supported, got=%s", args[0].Type())
	}
}

// Return the first element of the given array.
func _firstFn(args ...Object) Object {
	if len(args) != 1 {
		return newError("wrong number of arguments. got=%d, want=1",
			len(args))
	}
	if args[0].Type() != ARRAY_OBJ {
		return newError(
			"argument to `first` must be ARRAY, got=%s",
			args[0].Type(),
		)
	}

	arr := args[0].(*Array)
	if len(arr.Elements) > 0 {
		return arr.Elements[0]
	}

	return nil
}

// Return the last element of the given array.
func _lastFn(args ...Object) Object {
	if len(args) != 1 {
		return newError("wrong number of arguments. got=%d, want=1",
			len(args))
	}
	if args[0].Type() != ARRAY_OBJ {
		return newError(
			"argument to `last` must be ARRAY, got=%s",
			args[0].Type(),
		)
	}

	arr := args[0].(*Array)
	length := len(arr.Elements)
	if length > 0 {
		return arr.Elements[length-1]
	}

	return nil
}

func _restFn(args ...Object) Object {
	if len(args) != 1 {
		return newError("wrong number of arguments. got=%d, want=1",
			len(args))
	}
	if args[0].Type() != ARRAY_OBJ {
		return newError(
			"argument to `rest` must be ARRAY, got=%s",
			args[0].Type(),
		)
	}

	arr := args[0].(*Array)
	length := len(arr.Elements)
	if length > 0 {
		newElements := make([]Object, length-1)
		copy(newElements, arr.Elements[1:length])

		return &Array{
			Elements: newElements,
		}
	}

	return nil
}

func _pushFn(args ...Object) Object {
	if len(args) != 2 {
		return newError("wrong number of arguments. got=%d, want=2",
			len(args))
	}
	if args[0].Type() != ARRAY_OBJ {
		return newError("argument to `push` must be ARRAY, got %s",
			args[0].Type())
	}

	arr := args[0].(*Array)
	length := len(arr.Elements)

	newElements := make([]Object, length+1)
	copy(newElements, arr.Elements)
	newElements[length] = args[1]

	return &Array{Elements: newElements}
}

func _putsFn(args ...Object) Object {
	for _, arg := range args {
		fmt.Println(arg.Inspect())
	}

	return nil
}
