package runtime

import (
	"github.com/duwa-lang/duwa/src/ast"
)

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

type Event struct {
	Type     EventType
	Node     ast.Node
	Data     map[string]interface{}
	Location Location
}

type Location struct {
	File   string
	Line   int
	Column int
}

type CallFrame struct {
	FunctionName string
	Location     Location
	Locals       map[string]string
}

type RuntimeObserver interface {
	OnEvent(event Event)
	Name() string
	Enabled() bool
}

type ObserverManager struct {
	observers      []RuntimeObserver
	hasObservers   bool
	needsRecompute bool
}

func NewObserverManager() *ObserverManager {
	return &ObserverManager{
		observers:      make([]RuntimeObserver, 0),
		hasObservers:   false,
		needsRecompute: false,
	}
}

func (om *ObserverManager) Register(observer RuntimeObserver) {
	om.observers = append(om.observers, observer)
	om.needsRecompute = true
}

func (om *ObserverManager) Unregister(name string) {
	for i, obs := range om.observers {
		if obs.Name() == name {
			om.observers = append(om.observers[:i], om.observers[i+1:]...)
			om.needsRecompute = true
			return
		}
	}
}

func (om *ObserverManager) Notify(event Event) {
	for _, observer := range om.observers {
		if observer.Enabled() {
			observer.OnEvent(event)
		}
	}
}

func (om *ObserverManager) HasObservers() bool {
	if len(om.observers) == 0 {
		om.hasObservers = false
		return false
	}

	if om.needsRecompute {
		om.hasObservers = false
		for _, observer := range om.observers {
			if observer.Enabled() {
				om.hasObservers = true
				break
			}
		}
		om.needsRecompute = false
	}

	return om.hasObservers
}
