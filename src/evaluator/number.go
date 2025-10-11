package evaluator

import "github.com/duwa-lang/duwa/src/object"

func evalIntegerInfixExpression(
	operator string,
	left, right object.Object,
) object.Object {
	leftVal := left.(*object.Integer).Value
	rightVal := right.(*object.Integer).Value
	switch operator {
	case "+":
		return &object.Integer{Value: leftVal.Add(rightVal)}
	case "-":
		return &object.Integer{Value: leftVal.Sub(rightVal)}
	case "*":
		return &object.Integer{Value: leftVal.Mul(rightVal)}
	case "/":
		return &object.Integer{Value: leftVal.Div(rightVal)}
	case "<":
		return nativeBoolToBooleanObject(leftVal.LessThan(rightVal))
	case "<=":
		return nativeBoolToBooleanObject(leftVal.LessThanOrEqual(rightVal))
	case ">":
		return nativeBoolToBooleanObject(leftVal.GreaterThan(rightVal))
	case ">=":
		return nativeBoolToBooleanObject(leftVal.GreaterThanOrEqual(rightVal))
	case "==":
		return nativeBoolToBooleanObject(leftVal.Equal(rightVal))
	case "!=":
		return nativeBoolToBooleanObject(!leftVal.Equal(rightVal))
	default:
		return newError("unknown operator: %s %s %s",
			left.Type(), operator, right.Type())
	}
}
