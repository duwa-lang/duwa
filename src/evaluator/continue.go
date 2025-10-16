package evaluator

import (
	"github.com/duwa-lang/duwa/src/ast"
	"github.com/duwa-lang/duwa/src/object"
)

func evaluateContinue(node *ast.ContinueStatement, env *object.Environment) object.Object {
	return &object.Continue{}
}
