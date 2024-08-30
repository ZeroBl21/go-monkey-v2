package evaluator

import (
	"fmt"
	"unicode/utf8"

	"github.com/ZeroBl21/go-monkey/src/object"
)

// The built-in functions / standard-library methods are stored here.
var builtins = map[string]*object.Builtin{}

// RegisterBuiltin registers a built-in function. This is used to register
// our "standard library" functions.
func RegisterBuiltin(name string, fn object.BuiltinFunction) {
	builtins[name] = &object.Builtin{Fn: fn}
}

// length of item in runes
func _lenFn(args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError("wrong number of arguments. got=%d, want=1",
			len(args))
	}
	switch arg := args[0].(type) {
	case *object.String:
		return &object.Integer{Value: int64(utf8.RuneCountInString(arg.Value))}
	case *object.Array:
		return &object.Integer{Value: int64(len(arg.Elements))}
	default:
		return newError("argument to `len` not supported, got=%s",
			args[0].Type())
	}
}

// length of item but counting bytes individually
func _unicodeLenFn(args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError("wrong number of arguments. got=%d, want=1",
			len(args))
	}
	switch arg := args[0].(type) {
	case *object.String:
		return &object.Integer{Value: int64(len(arg.Value))}
	default:
		return newError("argument to `unicodeLen` not supported, got=%s",
			args[0].Type())
	}
}

// Return the last element of the given array.
func _lastFn(args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError("wrong number of arguments. got=%d, want=1",
			len(args))
	}
	if args[0].Type() != object.ARRAY_OBJ {
		return newError("argument to `last` must be ARRAY, got=%s", args[0].Type())
	}

	arr := args[0].(*object.Array)
	length := len(arr.Elements)
	if length > 0 {
		return arr.Elements[length-1]
	}

	return NULL
}

// Return the first element of the given array.
func _firstFn(args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError("wrong number of arguments. got=%d, want=1",
			len(args))
	}
	if args[0].Type() != object.ARRAY_OBJ {
		return newError("argument to `first` must be ARRAY, got=%s", args[0].Type())
	}

	arr := args[0].(*object.Array)
	if len(arr.Elements) > 0 {
		return arr.Elements[0]
	}

	return NULL
}

func _restFn(args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError("wrong number of arguments. got=%d, want=1",
			len(args))
	}
	if args[0].Type() != object.ARRAY_OBJ {
		return newError("argument to `rest` must be ARRAY, got=%s", args[0].Type())
	}

	arr := args[0].(*object.Array)
	length := len(arr.Elements)
	if length > 0 {
		newElements := make([]object.Object, length-1)
		copy(newElements, arr.Elements[1:length])

		return &object.Array{
			Elements: newElements,
		}
	}

	return NULL
}

func _pushFn(args ...object.Object) object.Object {
	if len(args) != 2 {
		return newError("wrong number of arguments. got=%d, want=2",
			len(args))
	}
	if args[0].Type() != object.ARRAY_OBJ {
		return newError("argument to `push` must be ARRAY, got %s",
			args[0].Type())
	}

	arr := args[0].(*object.Array)
	length := len(arr.Elements)

	newElements := make([]object.Object, length+1)
	copy(newElements, arr.Elements)
	newElements[length] = args[1]

	return &object.Array{Elements: newElements}
}

func _putsFn(args ...object.Object) object.Object {
	for _, arg := range args {
		fmt.Println(arg.Inspect())
	}

	return NULL
}

func init() {
	RegisterBuiltin("len", _lenFn)
	RegisterBuiltin("unicodeLen", _unicodeLenFn)
	RegisterBuiltin("first", _firstFn)
	RegisterBuiltin("last", _lastFn)
	RegisterBuiltin("rest", _restFn)
	RegisterBuiltin("push", _pushFn)
	RegisterBuiltin("puts", _putsFn)
}
