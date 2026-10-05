package evaluator

import (
	"github.com/duwa-lang/duwa/src/ast"
	"github.com/duwa-lang/duwa/src/object"
	"github.com/duwa-lang/duwa/src/values"
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

		return newErrorNode(node.Token, "undefined library %s for Library %s", method.Value, receiver.Name)
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
			return newErrorNode(node.Token, "undefined method %s for %s", node.Method.(*ast.Identifier).Value, node.Left.String())
		}
	}

	return result
}

func evaluateInstanceMethod(node *ast.MethodExpression, receiverInstance *object.Instance, name string, arguments []object.Object) object.Object {
	method, ok := receiverInstance.Class.Env.Get(name)

	if !ok {
		return newErrorNode(node.Token, "undefined instance method %s for class %s", name, receiverInstance.Class.Name.Value)
	}

	if method, ok := method.(*object.Function); ok {
		// Validate argument count
		if len(arguments) != len(method.Parameters) {
			return newErrorNode(node.Token, "argument count mismatch for method %s: expected %d arguments, got %d",
				name, len(method.Parameters), len(arguments))
		}
		// Create environment for method execution that extends the INSTANCE environment
		// This ensures that property assignments in the method affect this specific instance
		extendedEnv := object.NewEnclosedEnvironment(receiverInstance.Env)
		for paramIdx, param := range method.Parameters {
			extendedEnv.Set(param.Value, arguments[paramIdx])
		}
		return Eval(method.Body, extendedEnv)
	} else {
		return newErrorNode(node.Token, "not a method: %s", name)
	}
}

func evaluatePackageMethod(node *ast.MethodExpression, receiver *object.Package, name string, arguments []object.Object) object.Object {
	function, err := receiver.GetPackageFunction(name)
	if err != nil {
		return newErrorNode(node.Token, "%s", err.Error())
	}

	// Validate argument count
	if len(arguments) != len(function.Parameters) {
		return newErrorNode(node.Token, "argument count mismatch for function %s: expected %d arguments, got %d",
			name, len(function.Parameters), len(arguments))
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
