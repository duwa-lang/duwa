package object

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/duwa-lang/duwa/src/ast"
)

const FUNCTION_OBJ = "FUNCTION"

type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
}

func (f *Function) Type() ObjectType { return FUNCTION_OBJ }

func (f *Function) String() string {
	var out bytes.Buffer
	params := []string{}
	for _, p := range f.Parameters {
		params = append(params, p.String())
	}
	out.WriteString("ndondomeko")
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") {\n")
	out.WriteString(f.Body.String())
	out.WriteString("\n}")
	return out.String()
}

func (i *Function) Method(method string, args []Object) (Object, bool) {
	return nil, false
}

func (function *Function) Evaluate(env *Environment, args []Object) Object {
	// Check argument count matches parameter count
	if len(args) != len(function.Parameters) {
		return &Error{
			Message: fmt.Sprintf("runtime error: argument count mismatch: expected %d arguments, got %d",
				len(function.Parameters), len(args)),
		}
	}

	extendedEnv := extendFunctionEnv(env, function, args)
	evaluated := evaluator(function.Body, extendedEnv)
	return unwrapReturnValue(evaluated)
}

func extendFunctionEnv(
	env *Environment,
	fn *Function,
	args []Object,
) *Environment {
	// Arguments are already validated in Evaluate, so safe to access
	for paramIdx, param := range fn.Parameters {
		env.Set(param.Value, args[paramIdx])
	}
	return env
}

func unwrapReturnValue(obj Object) Object {
	if returnValue, ok := obj.(*ReturnValue); ok {
		return returnValue.Value
	}
	return obj
}
