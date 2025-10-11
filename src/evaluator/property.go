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
	}

	return nil
}

func evaluateInstanceProperty(node *ast.PropertyExpression, instance *object.Instance) object.Object {
	property := node.Property.(*ast.Identifier)

	if instance.Env.Has(property.Value) {
		val, _ := instance.Env.Get(property.Value)
		return val
	}

	if instance.Class.Env.Has(property.Value) {
		val, _ := instance.Class.Env.Get(property.Value)
		return val
	}

	// If the property is not found in the instance or the class, return NULL
	instance.Env.Set(property.Value, values.NULL)
	return values.NULL
}
