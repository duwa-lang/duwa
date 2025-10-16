package object

import (
	"errors"
	"log/slog"

	"github.com/duwa-lang/duwa/src/runtime"
	"github.com/duwa-lang/duwa/src/runtime/native"
)

func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewDefaultEnvironment()
	env.outer = outer
	env.ObserverManager = outer.ObserverManager
	env.callStack = outer.callStack
	return env
}

func NewDefaultEnvironment() *Environment {
	logger := slog.Default()
	store := make(map[string]Object)
	console := native.NewConsole()
	return &Environment{
		store:           store,
		outer:           nil,
		Logger:          logger,
		Console:         console,
		ObserverManager: runtime.NewObserverManager(),
		callStack:       make([]runtime.CallFrame, 0),
	}
}

func CopyEnvironmentDefaults(outer *Environment) *Environment {
	return &Environment{
		store:           make(map[string]Object),
		outer:           nil,
		Logger:          outer.Logger,
		Console:         outer.Console,
		ObserverManager: outer.ObserverManager,
		callStack:       outer.callStack,
	}
}

func New(logger *slog.Logger, console runtime.Console) *Environment {
	s := make(map[string]Object)
	return &Environment{
		store:           s,
		outer:           nil,
		Logger:          logger,
		Console:         console,
		ObserverManager: runtime.NewObserverManager(),
		callStack:       make([]runtime.CallFrame, 0),
	}
}

type Environment struct {
	store     map[string]Object
	outer     *Environment
	directory string

	Logger  *slog.Logger
	Console runtime.Console

	// Debugging and observation
	ObserverManager *runtime.ObserverManager
	callStack       []runtime.CallFrame
}

func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]
	if !ok && e.outer != nil {
		obj, ok = e.outer.Get(name)
	}
	return obj, ok
}

func (e *Environment) Set(name string, val Object) Object {
	// TODO: Make sure we dont accidentally mutate data that is not in the current scope
	_, ok := e.store[name]
	if !ok && e.outer != nil {
		e.outer.Set(name, val)
		return val
	}
	e.store[name] = val

	// Notify observers of variable change
	e.NotifyObservers(runtime.EventVariableSet, map[string]interface{}{
		"name":  name,
		"value": val.String(),
	})

	return val
}

func (e *Environment) SetLocal(name string, val Object) Object {
	if val == nil {
		return NewError("cannot set nil value for variable: %s", name)
	}

	e.store[name] = val

	// Notify observers of variable change
	e.NotifyObservers(runtime.EventVariableSet, map[string]interface{}{
		"name":  name,
		"value": val.String(),
	})

	return val
}

func (e *Environment) Has(name string) bool {
	_, ok := e.store[name]
	if !ok && e.outer != nil {
		return e.outer.Has(name)
	}
	return ok
}

func (e *Environment) All() map[string]Object {
	return e.store
}

func (e *Environment) Delete(name string) {
	delete(e.store, name)
}

func (e *Environment) SetDirectory(directory string) {
	e.directory = directory
}

func (e *Environment) GetDirectory() string {
	directory := e.directory

	if directory == "" && e.outer != nil {
		directory = e.outer.GetDirectory()
	}

	return directory
}

func (e *Environment) Call(function string, args []Object) Object {
	if object, ok := e.Get(function); ok {
		if function, ok := object.(*Function); ok {
			return function.Evaluate(e, args)
		}
	}

	return NewError("function not found: %s", function)
}

func (e *Environment) CallE(function string, args []Object) (Object, error) {
	result := e.Call(function, args)
	if errObj, ok := result.(*Error); ok {
		return nil, errors.New(errObj.Message)
	}
	return result, nil
}

// PushCallFrame adds a new frame to the call stack
func (e *Environment) PushCallFrame(frame runtime.CallFrame) {
	e.callStack = append(e.callStack, frame)
}

// PopCallFrame removes the top frame from the call stack
func (e *Environment) PopCallFrame() {
	if len(e.callStack) > 0 {
		e.callStack = e.callStack[:len(e.callStack)-1]
	}
}

// GetCallStack returns a copy of the current call stack
func (e *Environment) GetCallStack() []runtime.CallFrame {
	stack := make([]runtime.CallFrame, len(e.callStack))
	copy(stack, e.callStack)
	return stack
}

// GetCurrentFrame returns the current call frame (top of stack)
func (e *Environment) GetCurrentFrame() *runtime.CallFrame {
	if len(e.callStack) > 0 {
		return &e.callStack[len(e.callStack)-1]
	}
	return nil
}

// NotifyObservers sends an event to all registered observers
func (e *Environment) NotifyObservers(eventType runtime.EventType, data map[string]interface{}) {
	if e.ObserverManager != nil && e.ObserverManager.HasObservers() {
		event := runtime.Event{
			Type: eventType,
			Data: data,
		}
		e.ObserverManager.Notify(event)
	}
}
