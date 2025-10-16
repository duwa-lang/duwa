package debug

import (
	"fmt"
	"strings"

	"github.com/duwa-lang/duwa/src/object"
	"github.com/duwa-lang/duwa/src/token"
)

// method=lembaVars args=[] return={null}
// This method prints all variables in the current scope
//
// `Example`
// ```
// debug.lembaVars() # prints all variables
// ```
func methodPrintVars(env *object.Environment, tok token.Token, args ...object.Object) object.Object {
	vars := env.All()

	if len(vars) == 0 {
		env.Logger.Info("[DEBUG] No variables in current scope")
		return nil
	}

	env.Logger.Info("[DEBUG] Variables in current scope:")
	for name, value := range vars {
		env.Logger.Info(fmt.Sprintf("  %s = %s", name, value.String()))
	}

	return nil
}

// method=stackTrace args=[] return={null}
// This method prints the current call stack
//
// `Example`
// ```
// debug.stackTrace() # prints call stack
// ```
func methodStackTrace(env *object.Environment, tok token.Token, args ...object.Object) object.Object {
	stack := env.GetCallStack()

	if len(stack) == 0 {
		env.Logger.Info("[DEBUG] Call stack is empty")
		return nil
	}

	env.Logger.Info("[DEBUG] Call stack:")
	for i := len(stack) - 1; i >= 0; i-- {
		frame := stack[i]
		env.Logger.Info(fmt.Sprintf("  [%d] %s at %s:%d:%d",
			i,
			frame.FunctionName,
			frame.Location.File,
			frame.Location.Line,
			frame.Location.Column,
		))
	}

	return nil
}

// method=funso args=[any{expression}] return={mawu}
// This method inspects a value and returns detailed information about it
//
// `Example`
// ```
// nambala x = 10;
// debug.funso(x) # returns "INTEGER: 10"
// ```
func methodInspect(env *object.Environment, tok token.Token, args ...object.Object) object.Object {
	if len(args) != 1 {
		return &object.Error{Message: "debug.funso requires exactly one argument"}
	}

	obj := args[0]
	info := fmt.Sprintf("%s: %s", obj.Type(), obj.String())

	// Add additional type-specific information
	switch v := obj.(type) {
	case *object.Array:
		info += fmt.Sprintf(" (length: %d)", len(v.Elements))
	case *object.Map:
		info += fmt.Sprintf(" (size: %d)", len(v.Pairs))
	case *object.String:
		info += fmt.Sprintf(" (length: %d)", len(v.Value))
	case *object.Function:
		params := make([]string, len(v.Parameters))
		for i, p := range v.Parameters {
			params[i] = p.Value
		}
		info += fmt.Sprintf(" (parameters: %s)", strings.Join(params, ", "))
	}

	return &object.String{Value: info}
}

// method=zimitsidwa args=[] return={zoona}
// This method returns whether debugging is enabled
//
// `Example`
// ```
// ngati (debug.zimitsidwa()) { lemba("Debugging is on") }
// ```
func methodIsEnabled(env *object.Environment, tok token.Token, args ...object.Object) object.Object {
	if env.ObserverManager != nil && env.ObserverManager.HasObservers() {
		return &object.Boolean{Value: true}
	}
	return &object.Boolean{Value: false}
}

// library=debug
// This is the debug module
// It contains functions for debugging and inspecting runtime state
// It is used to help with development and troubleshooting
func Module() *object.LibraryModule {
	methods := map[string]*object.LibraryFunction{
		"lembaVars":  object.NewBuiltin("lembaVars", methodPrintVars),
		"stackTrace": object.NewBuiltin("stackTrace", methodStackTrace),
		"funso":      object.NewBuiltin("funso", methodInspect),
		"zimitsidwa": object.NewBuiltin("zimitsidwa", methodIsEnabled),
	}

	return object.NewBuiltInLibraryModule("debug", methods)
}
