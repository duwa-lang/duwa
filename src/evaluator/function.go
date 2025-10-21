package evaluator

import (
	"strings"

	"github.com/duwa-lang/duwa/src/ast"
	"github.com/duwa-lang/duwa/src/object"
	"github.com/duwa-lang/duwa/src/runtime"
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
	hasObservers := env.ObserverManager != nil && env.ObserverManager.HasObservers()

	if hasObservers {
		functionName := tok.Literal
		argStrs := make([]string, len(args))
		for i, arg := range args {
			argStrs[i] = arg.String()
		}
		argsStr := strings.Join(argStrs, ", ")

		env.NotifyObservers(runtime.EventFunctionCall, map[string]any{
			"function": functionName,
			"args":     argsStr,
		})

		env.PushCallFrame(runtime.CallFrame{
			FunctionName: functionName,
			Location: runtime.Location{
				File:   tok.File,
				Line:   tok.Pos.Line,
				Column: tok.Pos.Column,
			},
			Locals: map[string]string{},
		})
	}

	var result object.Object

	switch fn := fn.(type) {
	case *object.LibraryFunction:
		if res := fn.Function(env, tok, args...); res != nil {
			result = res
		} else {
			result = nil
		}
	case *object.Function:
		result = fn.Evaluate(env, args)
	case *object.Class:
		if tok.Literal != fn.Name.TokenLiteral() {
			result = newError("class name mismatch: expected %s, got %s", fn.Name.TokenLiteral(), tok.Literal)
		} else {
			result = fn.CreateInstance(tok.Literal, args)
		}
	default:
		result = newError("not a function: %s", fn.Type())
	}

	if hasObservers {
		env.PopCallFrame()

		resultStr := "null"
		if result != nil {
			resultStr = result.String()
		}
		env.NotifyObservers(runtime.EventFunctionReturn, map[string]interface{}{
			"function": tok.Literal,
			"result":   resultStr,
		})
	}

	return result
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
