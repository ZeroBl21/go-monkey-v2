package evaluator

import (
	"github.com/ZeroBl21/go-monkey/src/object"
)

// The built-in functions / standard-library methods are stored here.
var builtins = map[string]*object.Builtin{
	"len":   object.GetBuiltinByName("len"),
	"puts":  object.GetBuiltinByName("puts"),
	"first": object.GetBuiltinByName("first"),
	"last":  object.GetBuiltinByName("last"),
	"rest":  object.GetBuiltinByName("rest"),
	"push":  object.GetBuiltinByName("push"),

	"unicodeLen": object.GetBuiltinByName("unicodeLen"),
}
