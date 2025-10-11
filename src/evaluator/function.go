package evaluator

import (
	"github.com/duwa-lang/duwa/src/ast"
	"github.com/duwa-lang/duwa/src/object"
	"github.com/duwa-lang/duwa/src/token"
)

func evaluateFunctionCall(node *ast.CallExpression, env *object.Environment) object.Object {
	function := Eval(node.Function, env)
	if isError(function) {
		return function
	}
	args := evalExpressions(node.Arguments, env)
	if len(args) == 1 && isError(args[0]) {
		return args[0]
	}
	return applyFunction(node.Token, function, args, env)
}

func applyFunction(tok token.Token, fn object.Object, args []object.Object, env *object.Environment) object.Object {
	switch fn := fn.(type) {
	case *object.LibraryFunction:
		if result := fn.Function(env, tok, args...); result != nil {
			return result
		}
		return nil
	case *object.Function:
		return fn.Evaluate(env, args)
	case *object.Class:
		if tok.Literal != fn.Name.TokenLiteral() {
			return newError("class name mismatch: expected %s, got %s", fn.Name.TokenLiteral(), tok.Literal)
		}
		return fn.CreateInstance(tok.Literal, args)
	default:
		return newError("not a function: %s", fn.Type())
	}
}

func extendFunctionEnv(
	fn *object.Function,
	args []object.Object,
) *object.Environment {
	env := object.NewEnclosedEnvironment(fn.Env)
	for paramIdx, param := range fn.Parameters {
		env.Set(param.Value, args[paramIdx])
	}
	return env
}
