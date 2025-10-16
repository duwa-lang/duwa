package evaluator

import (
	"github.com/duwa-lang/duwa/src/ast"
	"github.com/duwa-lang/duwa/src/object"
)

// TODO: Handle type
func evaluateDeclaration(node *ast.VariableDeclarationStatement, env *object.Environment) object.Object {
	val := Eval(node.Value, env)
	if isError(val) {
		return val
	}
	// Use SetLocal to ensure the variable is created in the current scope
	// This is critical for class properties - they should be in the class environment,
	// not delegated to outer scopes
	env.SetLocal(node.Identifier.Value, val)
	return nil
}
