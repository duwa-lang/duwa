package evaluator

import (
	"github.com/sevenreup/duwa/src/ast"
	"github.com/sevenreup/duwa/src/object"
)

func evaluateClass(node *ast.ClassDeclStatement, env *object.Environment) object.Object {
	classEnv := object.NewEnclosedEnvironment(env)

	class := &object.Class{
		Name: node.Name,
		Env:  classEnv,
	}

	result := Eval(node.Body, classEnv)

	if isError(result) {
		return result
	}

	env.Set(class.Name.Value, class)

	return class
}
