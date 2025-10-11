package evaluator

import (
	"github.com/duwa-lang/duwa/src/ast"
	"github.com/duwa-lang/duwa/src/object"
)

func evaluateCompound(node *ast.Compound, env *object.Environment) object.Object {
	infix := &ast.InfixExpression{
		Token:    node.Token,
		Left:     node.Left,
		Operator: node.Operator[:len(node.Operator)-1],
		Right:    node.Right,
	}

	value := evalInfixExpression(infix, env)

	env.Set(node.Left.(*ast.Identifier).Value, value)

	return nil
}
