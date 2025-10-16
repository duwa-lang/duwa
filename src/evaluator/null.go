package evaluator

import (
	"github.com/duwa-lang/duwa/src/ast"
	"github.com/duwa-lang/duwa/src/object"
	"github.com/duwa-lang/duwa/src/values"
)

func evaluateNull(node *ast.NullLiteral, env *object.Environment) object.Object {
	return values.NULL
}
