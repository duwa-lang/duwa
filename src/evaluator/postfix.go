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
			return newErrorNode(node.Token, "identifier not found: %s", node.Token.Literal)
		}

		if value.Type() != object.INTEGER_OBJ {
			return newErrorNode(node.Token, "identifier is not a number: %s", node.Token.Literal)
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
			return newErrorNode(node.Token, "identifier not found: %s", node.Token.Literal)
		}

		if value.Type() != object.INTEGER_OBJ {
			return newErrorNode(node.Token, "identifier is not a number: %s", node.Token.Literal)
		}

		one := decimal.NewFromInt(1)

		newValue := &object.Integer{
			Value: value.(*object.Integer).Value.Sub(one),
		}

		env.Set(node.Token.Literal, newValue)

		return newValue
	default:
		return newErrorNode(node.Token, "unknown operator: %s", node.Operator)
	}
}
