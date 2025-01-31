package evaluator

import (
	"github.com/sevenreup/duwa/src/ast"
	"github.com/sevenreup/duwa/src/object"
)

func evalFile(file *ast.File, env *object.Environment) object.Object {
	var result object.Object
	for _, statement := range file.Statements {
		result = Eval(statement, env)
		switch result := result.(type) {
		case *object.ReturnValue:
			return result.Value
		case *object.Error:
			return result
		}
	}
	return result
}
