package evaluator

import (
	"github.com/duwa-lang/duwa/src/ast"
	"github.com/duwa-lang/duwa/src/object"
	"github.com/duwa-lang/duwa/src/values"
)

func evaluateProperty(node *ast.PropertyExpression, env *object.Environment) object.Object {
	left := Eval(node.Left, env)

	if isError(left) {
		return left
	}

	switch receiver := left.(type) {
	case *object.Instance:
		return evaluateInstanceProperty(node, receiver)
	case *object.LibraryModule:
		return evaluateLibraryModuleProperty(node, receiver)
	}

	return nil
}

func evaluateLibraryModuleProperty(node *ast.PropertyExpression, module *object.LibraryModule) object.Object {
	property := node.Property.(*ast.Identifier)

	if value, ok := module.Properties[property.Value]; ok {
		return value
	}

	// Property not found in module
	return newError("%d:%d:%s: runtime error: undefined property %s for library module %s",
		node.Token.Pos.Line, node.Token.Pos.Column, node.Token.File, property.Value, module.Name)
}

func evaluateInstanceProperty(node *ast.PropertyExpression, instance *object.Instance) object.Object {
	property := node.Property.(*ast.Identifier)

	if val, ok := instance.Env.Get(property.Value); ok {
		return val
	}

	if val, ok := instance.Class.Env.Get(property.Value); ok {
		return val
	}

	return values.NULL
}
