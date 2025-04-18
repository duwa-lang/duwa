package lsp

import (
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

const lsName = "Duwa"

var (
	version string = "0.0.1"
	Handler protocol.Handler
)

func init() {
	// General Messages
	Handler.Initialize = initialize
	Handler.Initialized = initialized
	Handler.Shutdown = shutdown
	Handler.LogTrace = logTrace
	Handler.SetTrace = setTrace

	// Workspace
	Handler.WorkspaceDidRenameFiles = WorkspaceDidRenameFiles

	// Text Document Synchronization
	Handler.TextDocumentDidOpen = TextDocumentDidOpen
	Handler.TextDocumentDidChange = TextDocumentDidChange
	Handler.TextDocumentDidSave = TextDocumentDidSave
	Handler.TextDocumentDidClose = TextDocumentDidClose

	// Language Features
	Handler.TextDocumentCompletion = TextDocumentCompletion
	Handler.TextDocumentDocumentSymbol = TextDocumentDocumentSymbol
}

func initialize(context *glsp.Context, params *protocol.InitializeParams) (any, error) {
	capabilities := Handler.CreateServerCapabilities()
	log.Info("initialize")
	return protocol.InitializeResult{
		Capabilities: capabilities,
		ServerInfo: &protocol.InitializeResultServerInfo{
			Name:    lsName,
			Version: &version,
		},
	}, nil
}

func initialized(context *glsp.Context, params *protocol.InitializedParams) error {
	log.Info("initialized")
	return nil
}

func shutdown(context *glsp.Context) error {
	log.Info("shutdown")
	protocol.SetTraceValue(protocol.TraceValueOff)
	return nil
}

func setTrace(context *glsp.Context, params *protocol.SetTraceParams) error {
	log.Info("setTrace")
	protocol.SetTraceValue(params.Value)
	return nil
}

func logTrace(context *glsp.Context, params *protocol.LogTraceParams) error {
	log.Info("logTrace")
	return nil
}
