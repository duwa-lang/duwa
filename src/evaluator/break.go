package evaluator

import (
	"github.com/duwa-lang/duwa/src/ast"
	"github.com/duwa-lang/duwa/src/object"
)

func evaluateBreak(node *ast.BreakStatement, env *object.Environment) object.Object {
	return &object.Break{}
}
