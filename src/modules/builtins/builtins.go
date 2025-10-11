package builtins

import (
	"github.com/duwa-lang/duwa/src/object"
	"github.com/duwa-lang/duwa/src/token"
	"github.com/duwa-lang/duwa/src/values"
	"github.com/shopspring/decimal"
)

// type=builtin-func method=kuNambala args=[string{value}] return={integer}
// The kuNambala function converts a string to a number.
func BuiltInParseStringToNumber(env *object.Environment, tok token.Token, args ...object.Object) object.Object {
	if len(args) != 1 {
		// TODO: Return error dont panic
		panic("parse requires exactly one argument")
	}

	if args[0].Type() != object.STRING_OBJ {
		return values.NULL
	}

	str := args[0].(*object.String).Value
	number, err := decimal.NewFromString(str)

	if err != nil {
		return values.NULL
	}

	return &object.Integer{Value: number}
}
