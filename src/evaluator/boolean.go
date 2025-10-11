package evaluator

import (
	"github.com/duwa-lang/duwa/src/object"
	"github.com/duwa-lang/duwa/src/values"
)

func nativeBoolToBooleanObject(input bool) *object.Boolean {
	if input {
		return values.TRUE
	}
	return values.FALSE
}
