package evaluator

import (
	"github.com/sevenreup/duwa/src/ast"
	"github.com/sevenreup/duwa/src/object"
)

// TODO: Handle type
func evaluateDeclaration(node *ast.VariableDeclStatement, env *object.Environment) object.Object {
	val := Eval(node.Value, env)
	if isError(val) {
		return val
	}
	env.Set(node.Identifier.Value, val)
	return nil
}
