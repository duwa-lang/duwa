package evaluator

import (
	"github.com/duwa-lang/duwa/src/ast"
	"github.com/duwa-lang/duwa/src/object"
)

func evalMapExpression(node *ast.MapExpression,
	env *object.Environment) object.Object {
	pairs := make(map[object.MapKey]object.MapPair)

	for keyNode, valueNode := range node.Pairs {
		identifier, ok := keyNode.(*ast.Identifier)

		if ok {
			keyNode = &ast.StringLiteral{
				Token: identifier.Token,
				Value: identifier.Value,
			}
		}

		key := Eval(keyNode, env)

		if isError(key) {
			return key
		}

		mapKey, ok := key.(object.Mappable)

		if !ok {
			return newError("%d:%d:%s: runtime error: unusable as map key: %s", node.Token.Pos.Line, node.Token.Pos.Column, node.Token.File, key.Type())
		}

		value := Eval(valueNode, env)

		if isError(value) {
			return value
		}

		hashed := mapKey.MapKey()

		pairs[hashed] = object.MapPair{Key: key, Value: value}
	}

	return &object.Map{Pairs: pairs}
}
