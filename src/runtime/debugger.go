package runtime

import (
	"fmt"
	"log/slog"
	"strings"
)

// DebuggerObserver is a built-in observer that logs execution details
type DebuggerObserver struct {
	enabled bool
	logger  *slog.Logger
	verbose bool
}

// NewDebuggerObserver creates a new debugger observer
func NewDebuggerObserver(logger *slog.Logger, verbose bool) *DebuggerObserver {
	return &DebuggerObserver{
		enabled: true,
		logger:  logger,
		verbose: verbose,
	}
}

// Name returns the observer name
func (d *DebuggerObserver) Name() string {
	return "debugger"
}

// Enabled returns whether the observer is active
func (d *DebuggerObserver) Enabled() bool {
	return d.enabled
}

// SetEnabled enables or disables the observer
func (d *DebuggerObserver) SetEnabled(enabled bool) {
	d.enabled = enabled
}

// OnEvent handles runtime events
func (d *DebuggerObserver) OnEvent(event Event) {
	switch event.Type {
	case EventBeforeEval:
		d.onBeforeEval(event)
	case EventAfterEval:
		d.onAfterEval(event)
	case EventVariableSet:
		d.onVariableSet(event)
	case EventVariableGet:
		if d.verbose {
			d.onVariableGet(event)
		}
	case EventFunctionCall:
		d.onFunctionCall(event)
	case EventFunctionReturn:
		d.onFunctionReturn(event)
	case EventError:
		d.onError(event)
	}
}

func (d *DebuggerObserver) onBeforeEval(event Event) {
	if !d.verbose {
		return
	}

	nodeType, _ := event.Data["node_type"].(string)
	line, _ := event.Data["line"].(int)
	column, _ := event.Data["column"].(int)

	if line > 0 {
		d.logger.Info(fmt.Sprintf("[DEBUG] Evaluating %s at %d:%d", nodeType, line, column))
	}
}

func (d *DebuggerObserver) onAfterEval(event Event) {
	if !d.verbose {
		return
	}

	nodeType, _ := event.Data["node_type"].(string)
	result, _ := event.Data["result"].(string)

	d.logger.Info(fmt.Sprintf("[DEBUG] Evaluated %s => %s", nodeType, result))
}

func (d *DebuggerObserver) onVariableSet(event Event) {
	name, _ := event.Data["name"].(string)
	value, _ := event.Data["value"].(string)

	d.logger.Info(fmt.Sprintf("[DEBUG] Variable set: %s = %s", name, value))
}

func (d *DebuggerObserver) onVariableGet(event Event) {
	name, _ := event.Data["name"].(string)
	value, _ := event.Data["value"].(string)

	d.logger.Info(fmt.Sprintf("[DEBUG] Variable get: %s => %s", name, value))
}

func (d *DebuggerObserver) onFunctionCall(event Event) {
	functionName, _ := event.Data["function"].(string)
	args, _ := event.Data["args"].(string)

	d.logger.Info(fmt.Sprintf("[DEBUG] Calling function: %s(%s)", functionName, args))
}

func (d *DebuggerObserver) onFunctionReturn(event Event) {
	functionName, _ := event.Data["function"].(string)
	result, _ := event.Data["result"].(string)

	d.logger.Info(fmt.Sprintf("[DEBUG] Function %s returned: %s", functionName, result))
}

func (d *DebuggerObserver) onError(event Event) {
	message, _ := event.Data["message"].(string)

	d.logger.Error(fmt.Sprintf("[DEBUG] Error: %s", message))
}

// TraceObserver is a simpler observer that prints a trace of execution
type TraceObserver struct {
	enabled bool
	logger  *slog.Logger
	depth   int
}

// NewTraceObserver creates a new trace observer
func NewTraceObserver(logger *slog.Logger) *TraceObserver {
	return &TraceObserver{
		enabled: true,
		logger:  logger,
		depth:   0,
	}
}

// Name returns the observer name
func (t *TraceObserver) Name() string {
	return "trace"
}

// Enabled returns whether the observer is active
func (t *TraceObserver) Enabled() bool {
	return t.enabled
}

// SetEnabled enables or disables the observer
func (t *TraceObserver) SetEnabled(enabled bool) {
	t.enabled = enabled
}

// OnEvent handles runtime events
func (t *TraceObserver) OnEvent(event Event) {
	indent := strings.Repeat("  ", t.depth)

	switch event.Type {
	case EventFunctionCall:
		functionName, _ := event.Data["function"].(string)
		args, _ := event.Data["args"].(string)
		t.logger.Info(fmt.Sprintf("%s→ %s(%s)", indent, functionName, args))
		t.depth++
	case EventFunctionReturn:
		t.depth--
		functionName, _ := event.Data["function"].(string)
		result, _ := event.Data["result"].(string)
		t.logger.Info(fmt.Sprintf("%s← %s => %s", indent, functionName, result))
	case EventError:
		message, _ := event.Data["message"].(string)
		t.logger.Error(fmt.Sprintf("%s✗ Error: %s", indent, message))
	}
}
