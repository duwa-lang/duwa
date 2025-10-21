package evaluator

import (
	"github.com/duwa-lang/duwa/src/ast"
	"github.com/duwa-lang/duwa/src/object"
)

func evalExpressions(
	exps []ast.Expression,
	env *object.Environment,
) []object.Object {
	result := make([]object.Object, 0, len(exps))
	for _, e := range exps {
		evaluated := Eval(e, env)
		if isError(evaluated) {
			return []object.Object{evaluated}
		}
		result = append(result, evaluated)
	}
	return result
}
