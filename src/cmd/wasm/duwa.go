//go:build js && wasm

package main

import (
	"context"
	"fmt"
	"log/slog"
	"runtime"
	"syscall/js"

	"github.com/duwa-lang/duwa/src/duwa"
	"github.com/duwa-lang/duwa/src/object"
	"github.com/duwa-lang/duwa/src/runtime/wasm"
)

var compiler *duwa.Duwa

type WasmConsoleHandler struct {
}

func main() {
	// Ensure the garbage collector runs more frequently
	runtime.SetFinalizer(compiler, nil)
	logger := slog.New(&WasmConsoleHandler{})
	console := wasm.NewConsole()
	compiler = duwa.New(object.New(logger, console))

	// Register the function in the global scope
	js.Global().Set("duwaRun", js.FuncOf(wasm.WrapFunction(run)))

	fmt.Println("Duwa is ready")

	// Keep the program running
	<-make(chan bool)
}

func run(this js.Value, inputs []js.Value) interface{} {
	if len(inputs) < 1 {
		return wasm.Wrap("Please provide a file to run")
	}

	file := inputs[0].String()
	result := compiler.Run(file)
	if result == nil {
		return wasm.Wrap(nil)
	}
	if object.IsError(result) {
		emitConsoleLogEvent(slog.Record{
			Level:   slog.LevelError,
			Message: fmt.Sprintf("%v", result),
		})
		return wasm.Wrap(result)
	}
	return wasm.Wrap(result.String())
}

func (h *WasmConsoleHandler) Enabled(_ context.Context, level slog.Level) bool {
	return true
}

func emitConsoleLogEvent(r slog.Record) {
	eventInit := js.Global().Get("Object").New()
	eventInit.Set("message", r.Message)
	eventInit.Set("level", slogLevelToConsoleLevel(r.Level))
	logType := "runtime"
	for attr := range r.Attrs {
		if attr.Key == "type" {
			logType = attr.Value.String()
		}
	}
	eventInit.Set("type", logType)
	wasm.DispatchEvent("duwaLogEvent", eventInit)
}

func (h *WasmConsoleHandler) Handle(ctx context.Context, r slog.Record) error {
	emitConsoleLogEvent(r)
	return nil
}

func (h *WasmConsoleHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return h
}

func (h *WasmConsoleHandler) WithGroup(name string) slog.Handler {
	return h
}

func (h *WasmConsoleHandler) Handler() slog.Handler {
	return h
}

func slogLevelToConsoleLevel(level slog.Level) string {
	switch level {
	case slog.LevelDebug:
		return "debug"
	case slog.LevelInfo:
		return "info"
	case slog.LevelWarn:
		return "warn"
	case slog.LevelError:
		return "error"
	default:
		return "info"
	}
}
