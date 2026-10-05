package evaluator

import (
	"github.com/duwa-lang/duwa/src/ast"
	"github.com/duwa-lang/duwa/src/modules"
	"github.com/duwa-lang/duwa/src/object"
)

func evalIdentifier(
	node *ast.Identifier,
	env *object.Environment,
) object.Object {
	if libraryFunction, ok := modules.Functions[node.Value]; ok {
		return libraryFunction
	}

	val, ok := env.Get(node.Value)
	if !ok {
		return newErrorNode(node.Token, "identifier not found: %s", node.Value)
	}

	return val
}
