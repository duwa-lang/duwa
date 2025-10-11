package runtime

import (
	"github.com/duwa-lang/duwa/src/ast"
)

// EventType represents different types of runtime events
type EventType string

const (
	EventBeforeEval     EventType = "before_eval"
	EventAfterEval      EventType = "after_eval"
	EventVariableSet    EventType = "variable_set"
	EventVariableGet    EventType = "variable_get"
	EventFunctionCall   EventType = "function_call"
	EventFunctionReturn EventType = "function_return"
	EventError          EventType = "error"
)

// Event represents a runtime event with associated data
type Event struct {
	Type     EventType
	Node     ast.Node
	Data     map[string]interface{}
	Location Location
}

// Location represents a source code location
type Location struct {
	File   string
	Line   int
	Column int
}

// CallFrame represents a single frame in the call stack
type CallFrame struct {
	FunctionName string
	Location     Location
	Locals       map[string]string // variable name -> value representation
}

// RuntimeObserver is an interface for observing runtime events
type RuntimeObserver interface {
	// OnEvent is called when a runtime event occurs
	OnEvent(event Event)

	// Name returns the name of this observer
	Name() string

	// Enabled returns whether this observer is currently active
	Enabled() bool
}

// ObserverManager manages multiple runtime observers
type ObserverManager struct {
	observers []RuntimeObserver
}

// NewObserverManager creates a new observer manager
func NewObserverManager() *ObserverManager {
	return &ObserverManager{
		observers: make([]RuntimeObserver, 0),
	}
}

// Register adds a new observer
func (om *ObserverManager) Register(observer RuntimeObserver) {
	om.observers = append(om.observers, observer)
}

// Unregister removes an observer by name
func (om *ObserverManager) Unregister(name string) {
	for i, obs := range om.observers {
		if obs.Name() == name {
			om.observers = append(om.observers[:i], om.observers[i+1:]...)
			return
		}
	}
}

// Notify sends an event to all enabled observers
func (om *ObserverManager) Notify(event Event) {
	for _, observer := range om.observers {
		if observer.Enabled() {
			observer.OnEvent(event)
		}
	}
}

// HasObservers returns true if there are any enabled observers
func (om *ObserverManager) HasObservers() bool {
	for _, observer := range om.observers {
		if observer.Enabled() {
			return true
		}
	}
	return false
}
