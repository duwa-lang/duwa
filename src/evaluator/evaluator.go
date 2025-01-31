package evaluator

import (
	"fmt"

	"github.com/sevenreup/duwa/src/ast"
	"github.com/sevenreup/duwa/src/object"
	"github.com/sevenreup/duwa/src/values"
)

type Evaluator func(node ast.Node, env *object.Environment) object.Object

func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {
	case *ast.File:
		return evalFile(node, env)
	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)
	case *ast.PrefixExpression:
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(node.Operator, right)
	case *ast.InfixExpression:
		return evalInfixExpression(node, env)
	case *ast.Compound:
		return evaluateCompound(node, env)
	case *ast.BlockStatement:
		return evalBlockStatement(node, env)
	case *ast.IfExpression:
		return evalIfExpression(node, env)
	case *ast.ReturnStatement:
		val := Eval(node.ReturnValue, env)
		if isError(val) {
			return val
		}
		return &object.ReturnValue{Value: val}
	case *ast.Identifier:
		return evalIdentifier(node, env)
	case *ast.VariableDeclarationStatement:
		return evaluateDeclaration(node, env)
	case *ast.AssigmentStatement:
		return evaluateAssigment(node, env)
	case *ast.FunctionLiteral:
		function := &object.Function{Parameters: node.Parameters, Env: env, Body: node.Body}
		if node.Name != nil {
			env.Set(node.Name.Value, function)
		}
		return function
	case *ast.CallExpression:
		return evaluateFunctionCall(node, env)
	case *ast.PropertyExpression:
		return evaluateProperty(node, env)
	case *ast.ArrayLiteral:
		elements := evalExpressions(node.Elements, env)
		if len(elements) == 1 && isError(elements[0]) {
			return elements[0]
		}
		return &object.Array{Elements: elements}
	case *ast.IndexExpression:
		return evalIndexExpression(node, env)
	case *ast.MethodExpression:
		return evaluateMethod(node, env)
	case *ast.ForExpression:
		return evalForLoop(node, env)
	case *ast.WhileExpression:
		return evaluateWhile(node, env)
	// Expressions
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.MapExpression:
		return evalMapExpression(node, env)
	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Value)
	case *ast.StringLiteral:
		return &object.String{Value: node.Value}
	case *ast.PostfixExpression:
		return evaluatePostfix(node, env)
	case *ast.ClassStatement:
		return evaluateClass(node, env)
	case *ast.BreakStatement:
		return evaluateBreak(node, env)
	case *ast.ContinueStatement:
		return evaluateContinue(node, env)
	case *ast.NullLiteral:
		return evaluateNull(node, env)
	}
	return nil
}

func isTruthy(obj object.Object) bool {
	switch obj {
	case values.NULL:
		return false
	case values.TRUE:
		return true
	case values.FALSE:
		return false
	default:
		return true
	}
}

func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJ
	}
	return false
}

// isTerminator determines if the referenced object is an error, break, or continue.
func isTerminator(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJ || obj.Type() == object.BREAK_OBJ || obj.Type() == object.CONTINUE_OBJ || obj.Type() == object.RETURN_VALUE_OBJ
	}

	return false
}
