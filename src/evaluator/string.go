package evaluator

import "github.com/duwa-lang/duwa/src/object"

func evalStringInfixExpression(
	operator string,
	left, right object.Object,
) object.Object {
	leftVal := _getString(left)
	rightVal := _getString(right)
	switch operator {
	case "+":
		return &object.String{Value: leftVal + rightVal}
	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)
	default:
		return newError("unknown operator: %s %s %s",
			left.Type(), operator, right.Type())
	}
}

func _getString(item object.Object) string {
	if item.Type() == object.STRING_OBJ {
		return item.(*object.String).Value
	}
	return item.String()
}
