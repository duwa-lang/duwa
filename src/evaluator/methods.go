package evaluator

import (
	"github.com/sevenreup/duwa/src/ast"
	"github.com/sevenreup/duwa/src/object"
	"github.com/sevenreup/duwa/src/values"
)

func evaluateMethod(node *ast.MethodExpression, env *object.Environment) object.Object {
	left := Eval(node.Left, env)

	if isError(left) {
		return left
	}

	arguments := evalExpressions(node.Arguments, env)

	if len(arguments) == 1 && isError(arguments[0]) {
		return arguments[0]
	}

	result, found := left.Method(node.Method.(*ast.Identifier).Value, arguments)

	if isError(result) {
		return result
	}

	switch receiver := left.(type) {
	case *object.LibraryModule:
		method := node.Method.(*ast.Identifier)

		if function, ok := receiver.Methods[method.Value]; ok {
			return applyFunction(node.Token, function, arguments, env)
		}

		return newError("%d:%d:%s: runtime error: undefined library %s for Library %s", node.Token.Pos.Line, node.Token.Pos.Column, node.Token.File, method.Value, receiver.Name)
	case *object.Instance:
		method := node.Method.(*ast.Identifier)
		evaluated := evaluateInstanceMethod(node, receiver, method.Value, arguments)

		if isError(evaluated) {
			return evaluated
		}

		return unwrapReturn(evaluated)
	case *object.Package:
		return evaluatePackageMethod(node, receiver, node.Method.(*ast.Identifier).Value, arguments)
	default:
		if !found {
			return newError("%d:%d:%s: runtime error: undefined method %s for %s", node.Token.Pos.Line, node.Token.Pos.Column, node.Token.File, node.Method.(*ast.Identifier).Value, node.Left.String())
		}
	}

	return result
}

func evaluateInstanceMethod(node *ast.MethodExpression, receiverInstance *object.Instance, name string, arguments []object.Object) object.Object {
	method, ok := receiverInstance.Class.Env.Get(name)

	if !ok {
		return newError("%d:%d:%s: runtime error: undefined instance method %s for class %s", node.Token.Pos.Line, node.Token.Pos.Column, node.Token.File, name, receiverInstance.Class.Name.Value)
	}

	if method, ok := method.(*object.Function); ok {
		extendedEnv := extendFunctionEnv(method, arguments)
		return Eval(method.Body, extendedEnv)
	} else {
		return newError("not a method: %s", name)
	}
}

func evaluatePackageMethod(node *ast.MethodExpression, receiver *object.Package, name string, arguments []object.Object) object.Object {
	function, err := receiver.GetPackageFunction(name)
	if err != nil {
		return newError("%d:%d:%s: runtime error: %s", node.Token.Pos.Line, node.Token.Pos.Column, node.Token.File, err.Error())
	}

	evaluated := Eval(function.Body, extendFunctionEnv(function, arguments))

	return unwrapReturn(evaluated)
}

func unwrapReturn(obj object.Object) object.Object {
	switch value := obj.(type) {
	case *object.Error:
		return obj
	case *object.ReturnValue:
		return value.Value
	}

	return values.NULL
}
