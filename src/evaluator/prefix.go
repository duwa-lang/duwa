package evaluator

import (
	"github.com/duwa-lang/duwa/src/object"
	"github.com/duwa-lang/duwa/src/values"
)

func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(right)
	default:
		return newError("unknown operator: %s%s", operator, right.Type())
	}
}

func evalBangOperatorExpression(right object.Object) object.Object {
	switch right {
	case values.TRUE:
		return values.FALSE
	case values.FALSE:
		return values.TRUE
	case values.NULL:
		return values.TRUE
	default:
		return values.FALSE
	}
}

func evalMinusPrefixOperatorExpression(right object.Object) object.Object {
	if right.Type() != object.INTEGER_OBJ {
		return newError("unknown operator:-%s", right.Type())
	}
	value := right.(*object.Integer).Value.Neg()
	return &object.Integer{Value: value}
}
