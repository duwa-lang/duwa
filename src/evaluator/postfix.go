package evaluator

import (
	"github.com/duwa-lang/duwa/src/ast"
	"github.com/duwa-lang/duwa/src/object"
	"github.com/shopspring/decimal"
)

func evaluatePostfix(node *ast.PostfixExpression, env *object.Environment) object.Object {
	switch node.Operator {
	case "++":
		value, ok := env.Get(node.Token.Literal)

		if !ok {
			return newError("%d:%d:%s: runtime error: identifier not found: %s", node.Token.Pos.Line, node.Token.Pos.Column, node.Token.File, node.Token.Literal)
		}

		if value.Type() != object.INTEGER_OBJ {
			return newError("%d:%d:%s: runtime error: identifier is not a number: %s", node.Token.Pos.Line, node.Token.Pos.Column, node.Token.File, node.Token.Literal)
		}

		one := decimal.NewFromInt(1)

		newValue := &object.Integer{
			Value: value.(*object.Integer).Value.Add(one),
		}

		env.Set(node.Token.Literal, newValue)

		return newValue
	case "--":
		value, ok := env.Get(node.Token.Literal)

		if !ok {
			return newError("%d:%d:%s: runtime error: identifier not found: %s", node.Token.Pos.Line, node.Token.Pos.Column, node.Token.File, node.Token.Literal)
		}

		if value.Type() != object.INTEGER_OBJ {
			return newError("%d:%d:%s: runtime error: identifier is not a number: %s", node.Token.Pos.Line, node.Token.Pos.Column, node.Token.File, node.Token.Literal)
		}

		one := decimal.NewFromInt(1)

		newValue := &object.Integer{
			Value: value.(*object.Integer).Value.Sub(one),
		}

		env.Set(node.Token.Literal, newValue)

		return newValue
	default:
		return newError("%d:%d:%s: runtime error: unknown operator: %s", node.Token.Pos.Line, node.Token.Pos.Column, node.Token.File, node.Operator)
	}
}
