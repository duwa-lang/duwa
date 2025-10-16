package evaluator

import (
	"github.com/duwa-lang/duwa/src/ast"
	"github.com/duwa-lang/duwa/src/object"
	"github.com/duwa-lang/duwa/src/values"
)

func evalIndexExpression(node *ast.IndexExpression, env *object.Environment) object.Object {
	left := Eval(node.Left, env)
	if isError(left) {
		return left
	}
	index := Eval(node.Index, env)
	if isError(index) {
		return index
	}

	switch {
	case left.Type() == object.ARRAY_OBJ && index.Type() == object.INTEGER_OBJ:
		return evalArrayIndexExpression(left, index)
	case left.Type() == object.MAP_OBJ:
		return evaluateMapIndex(node, left, index)
	default:
		return newError("index operator not supported: %s", left.Type())
	}
}

func evalArrayIndexExpression(array, index object.Object) object.Object {
	arrayObject := array.(*object.Array)
	idx := index.(*object.Integer).Value.IntPart()
	maxValue := int64(len(arrayObject.Elements) - 1)

	if idx < 0 || idx > maxValue {
		return values.NULL
	}

	return arrayObject.Elements[idx]
}

func evaluateMapIndex(node *ast.IndexExpression, left, index object.Object) object.Object {
	mapObject := left.(*object.Map)

	key, ok := index.(object.Mappable)

	if !ok {
		return newError("%d:%d:%s: runtime error: unusable as map key: %s", node.Token.Pos.Line, node.Token.Pos.Column, node.Token.File, index.Type())
	}

	pair, ok := mapObject.Pairs[key.MapKey()]

	if !ok {
		return values.NULL
	}

	return pair.Value
}
