package evaluator

import (
	"github.com/sevenreup/duwa/src/ast"
	"github.com/sevenreup/duwa/src/object"
)

func evalInfixExpression(
	node *ast.InfixExpression,
	env *object.Environment,
) object.Object {
	left := Eval(node.Left, env)
	if isError(left) {
		return left
	}
	right := Eval(node.Right, env)
	if isError(right) {
		return right
	}
	operator := node.Operator
	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalIntegerInfixExpression(operator, left, right)
	case left.Type() == object.STRING_OBJ || right.Type() == object.STRING_OBJ:
		return evalStringInfixExpression(operator, left, right)
	case operator == "==":
		return nativeBoolToBooleanObject(left == right)
	case operator == "!=":
		return nativeBoolToBooleanObject(left != right)
	case left.Type() == object.BOOLEAN_OBJ && right.Type() == object.BOOLEAN_OBJ:
		{
			l, _ := left.(*object.Boolean)
			r, _ := right.(*object.Boolean)
			if operator == "&&" {
				return nativeBoolToBooleanObject(l.Value && r.Value)
			} else if operator == "||" {
				return nativeBoolToBooleanObject(l.Value || r.Value)
			}
		}
	case left.Type() != right.Type():
		return newErrorNode(node.Token, "type mismatch: %s %s %s",
			left.Type(), operator, right.Type())
	}
	return newError("unknown operator: %s %s %s",
		left.Type(), operator, right.Type())
}
