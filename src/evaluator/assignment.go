package evaluator

import (
	"github.com/duwa-lang/duwa/src/ast"
	"github.com/duwa-lang/duwa/src/object"
	"github.com/duwa-lang/duwa/src/values"
)

func evaluateAssigment(node *ast.AssigmentStatement, env *object.Environment) object.Object {
	val := Eval(node.Value, env)
	if isError(val) {
		return val
	}
	switch smt := node.Identifier.(type) {
	case *ast.Identifier:
		return evaluateIdentifierAssignment(smt, val, env)
	case *ast.IndexExpression:
		return evaluateIndexAssignment(smt, val, env)
	}
	return nil
}

func evaluateIdentifierAssignment(node *ast.Identifier, val object.Object, env *object.Environment) object.Object {
	env.Set(node.Value, val)
	return nil
}

func evaluateIndexAssignment(node *ast.IndexExpression, val object.Object, env *object.Environment) object.Object {
	l := Eval(node.Left, env)
	index := Eval(node.Index, env)

	left, ok := l.(*object.Array)
	if !ok {
		return newError("index operator not supported: %s", left.Type())
	}

	idx := int(index.(*object.Integer).Value.IntPart())
	elements := left.Elements

	if idx < 0 {
		return newError("%d:%d:%s: runtime error: index out of range: %d", node.Token.Pos.Line, node.Token.Pos.Column, node.Token.File, idx)
	}

	if idx >= len(elements) {
		for i := len(elements); i <= idx; i++ {
			elements = append(elements, values.NULL)
		}

		left.Elements = elements
	}

	elements[idx] = val

	return nil
}
