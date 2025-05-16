package evaluator

import (
	"github.com/sevenreup/duwa/src/ast"
	"github.com/sevenreup/duwa/src/modules/all"
	"github.com/sevenreup/duwa/src/object"
)

func evalIdentifier(
	node *ast.Identifier,
	env *object.Environment,
) object.Object {
	if libraryFunction, ok := all.Functions[node.Value]; ok {
		return libraryFunction
	}

	val, ok := env.Get(node.Value)
	if !ok {
		return newError("identifier not found: " + node.Value)
	}

	return val
}
