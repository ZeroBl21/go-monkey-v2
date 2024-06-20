package evaluator

import (
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

func init() {
	RegisterBuiltin("len", _lenFn)
	RegisterBuiltin("unicodeLen", _unicodeLenFn)
}
