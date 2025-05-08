package lsp

import (
	"context"
	"os"
	"sync"
	"time"

	"go.lsp.dev/protocol"
	"go.uber.org/zap"
)

const lsName = "Duwa"

type DocumentStore struct {
	mu        sync.RWMutex
	documents map[protocol.DocumentURI]*TextDocument
}

func NewDocumentStore() *DocumentStore {
	return &DocumentStore{
		documents: make(map[protocol.DocumentURI]*TextDocument),
	}
}

type Server struct {
	Debug        bool
	LogBaseName  string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration

	docStore *DocumentStore

	Logger *zap.Logger
}

func NewServer(logBaseName string, debug bool) *Server {
	return &Server{
		Debug:       debug,
		LogBaseName: logBaseName,
		Logger:      logger,
		docStore: NewDocumentStore(),
	}
}

func (se *Server) Initialize(ctx context.Context, params *protocol.InitializeParams) (result *protocol.InitializeResult, err error) {
	return &protocol.InitializeResult{
		Capabilities: protocol.ServerCapabilities{
			DefinitionProvider: true,
			TextDocumentSync: &protocol.TextDocumentSyncOptions{
				OpenClose: true,
				Change:    protocol.TextDocumentSyncKindFull,
			},
			CompletionProvider: &protocol.CompletionOptions{
				ResolveProvider:   true,
				TriggerCharacters: []string{"."},
			},
		},
		ServerInfo: &protocol.ServerInfo{
			Name:    lsName,
			Version: "0.0.1",
		},
	}, nil
}

func (se *Server) Initialized(ctx context.Context, params *protocol.InitializedParams) (err error) {
	logger.Sugar().Debug("Initialized")
	return nil
}

func (se *Server) Shutdown(ctx context.Context) (err error) {
	logger.Sugar().Debug("Shutdown")
	return nil
}

func (se *Server) Exit(ctx context.Context) (err error) {
	logger.Sugar().Debug("Exit")
	os.Exit(0)
	return nil
}

func (se *Server) WorkDoneProgressCancel(ctx context.Context, params *protocol.WorkDoneProgressCancelParams) (err error) {
	return nil
}

func (se *Server) LogTrace(ctx context.Context, params *protocol.LogTraceParams) (err error) {
	return nil
}

func (se *Server) SetTrace(ctx context.Context, params *protocol.SetTraceParams) (err error) {
	return nil
}

func (se *Server) CodeAction(ctx context.Context, params *protocol.CodeActionParams) (result []protocol.CodeAction, err error) {
	return nil, nil
}

func (se *Server) CodeLens(ctx context.Context, params *protocol.CodeLensParams) (result []protocol.CodeLens, err error) {
	return nil, nil
}

func (se *Server) CodeLensResolve(ctx context.Context, params *protocol.CodeLens) (result *protocol.CodeLens, err error) {
	return nil, nil
}

func (se *Server) ColorPresentation(ctx context.Context, params *protocol.ColorPresentationParams) (result []protocol.ColorPresentation, err error) {
	return nil, nil
}

func (se *Server) Declaration(ctx context.Context, params *protocol.DeclarationParams) (result []protocol.Location, err error) {
	logger.Sugar().Debug("Declaration", params.TextDocument.URI, params.Position)
	return nil, nil
}

func (se *Server) DidChange(ctx context.Context, params *protocol.DidChangeTextDocumentParams) (err error) {
	// For simplicity, we're assuming full content sync mode
	if len(params.ContentChanges) > 0 {
		// In full sync mode, the first change contains the entire document
		se.UpdateDocument(params.TextDocument.URI, params.TextDocument.Version, params.ContentChanges[0].Text)
	}
	return nil
}

func (se *Server) DidChangeConfiguration(ctx context.Context, params *protocol.DidChangeConfigurationParams) (err error) {
	return nil
}

func (se *Server) DidChangeWatchedFiles(ctx context.Context, params *protocol.DidChangeWatchedFilesParams) (err error) {
	return nil
}

func (se *Server) DidChangeWorkspaceFolders(ctx context.Context, params *protocol.DidChangeWorkspaceFoldersParams) (err error) {
	return nil
}

func (se *Server) DidClose(ctx context.Context, params *protocol.DidCloseTextDocumentParams) (err error) {
	se.RemoveDocument(params.TextDocument.URI)
	return nil
}

func (se *Server) DidOpen(ctx context.Context, params *protocol.DidOpenTextDocumentParams) (err error) {
	se.UpdateDocument(params.TextDocument.URI, params.TextDocument.Version, params.TextDocument.Text)
	return nil
}

func (se *Server) DidSave(ctx context.Context, params *protocol.DidSaveTextDocumentParams) (err error) {
	logger.Sugar().Debug("DidSave", params.TextDocument.URI)
	return nil
}

func (se *Server) DocumentColor(ctx context.Context, params *protocol.DocumentColorParams) (result []protocol.ColorInformation, err error) {
	return nil, nil
}

func (se *Server) DocumentHighlight(ctx context.Context, params *protocol.DocumentHighlightParams) (result []protocol.DocumentHighlight, err error) {
	return nil, nil
}

func (se *Server) DocumentLink(ctx context.Context, params *protocol.DocumentLinkParams) (result []protocol.DocumentLink, err error) {
	return nil, nil
}

func (se *Server) DocumentLinkResolve(ctx context.Context, params *protocol.DocumentLink) (result *protocol.DocumentLink, err error) {
	return nil, nil
}

func (se *Server) DocumentSymbol(ctx context.Context, params *protocol.DocumentSymbolParams) (result []interface{}, err error) {
	logger.Sugar().Debug("DocumentSymbol", params.TextDocument.URI)
	return nil, nil
}

func (se *Server) ExecuteCommand(ctx context.Context, params *protocol.ExecuteCommandParams) (result interface{}, err error) {
	return nil, nil
}

func (se *Server) FoldingRanges(ctx context.Context, params *protocol.FoldingRangeParams) (result []protocol.FoldingRange, err error) {
	return nil, nil
}

func (se *Server) Formatting(ctx context.Context, params *protocol.DocumentFormattingParams) (result []protocol.TextEdit, err error) {
	return nil, nil
}

func (se *Server) Hover(ctx context.Context, params *protocol.HoverParams) (result *protocol.Hover, err error) {
	return nil, nil
}

func (se *Server) Implementation(ctx context.Context, params *protocol.ImplementationParams) (result []protocol.Location, err error) {
	return nil, nil
}

func (se *Server) OnTypeFormatting(ctx context.Context, params *protocol.DocumentOnTypeFormattingParams) (result []protocol.TextEdit, err error) {
	return nil, nil
}

func (se *Server) PrepareRename(ctx context.Context, params *protocol.PrepareRenameParams) (result *protocol.Range, err error) {
	return nil, nil
}

func (se *Server) RangeFormatting(ctx context.Context, params *protocol.DocumentRangeFormattingParams) (result []protocol.TextEdit, err error) {
	return nil, nil
}

func (se *Server) References(ctx context.Context, params *protocol.ReferenceParams) (result []protocol.Location, err error) {
	return nil, nil
}

func (se *Server) Rename(ctx context.Context, params *protocol.RenameParams) (result *protocol.WorkspaceEdit, err error) {
	return nil, nil
}

func (se *Server) SignatureHelp(ctx context.Context, params *protocol.SignatureHelpParams) (result *protocol.SignatureHelp, err error) {
	return nil, nil
}

func (se *Server) Symbols(ctx context.Context, params *protocol.WorkspaceSymbolParams) (result []protocol.SymbolInformation, err error) {
	return nil, nil
}

func (se *Server) TypeDefinition(ctx context.Context, params *protocol.TypeDefinitionParams) (result []protocol.Location, err error) {
	logger.Sugar().Debug("TypeDefinition", params.TextDocument.URI, params.Position)
	return nil, nil
}

func (se *Server) WillSave(ctx context.Context, params *protocol.WillSaveTextDocumentParams) (err error) {
	return nil
}

func (se *Server) WillSaveWaitUntil(ctx context.Context, params *protocol.WillSaveTextDocumentParams) (result []protocol.TextEdit, err error) {
	return nil, nil
}

func (se *Server) ShowDocument(ctx context.Context, params *protocol.ShowDocumentParams) (result *protocol.ShowDocumentResult, err error) {
	return nil, nil
}

func (se *Server) WillCreateFiles(ctx context.Context, params *protocol.CreateFilesParams) (result *protocol.WorkspaceEdit, err error) {
	return nil, nil
}

func (se *Server) DidCreateFiles(ctx context.Context, params *protocol.CreateFilesParams) (err error) {
	return nil
}

func (se *Server) WillRenameFiles(ctx context.Context, params *protocol.RenameFilesParams) (result *protocol.WorkspaceEdit, err error) {
	return nil, nil
}

func (se *Server) DidRenameFiles(ctx context.Context, params *protocol.RenameFilesParams) (err error) {
	return nil
}

func (se *Server) WillDeleteFiles(ctx context.Context, params *protocol.DeleteFilesParams) (result *protocol.WorkspaceEdit, err error) {
	return nil, nil
}

func (se *Server) DidDeleteFiles(ctx context.Context, params *protocol.DeleteFilesParams) (err error) {
	return nil
}

func (se *Server) CodeLensRefresh(ctx context.Context) (err error) {
	return nil
}

func (se *Server) PrepareCallHierarchy(ctx context.Context, params *protocol.CallHierarchyPrepareParams) (result []protocol.CallHierarchyItem, err error) {
	return nil, nil
}

func (se *Server) IncomingCalls(ctx context.Context, params *protocol.CallHierarchyIncomingCallsParams) (result []protocol.CallHierarchyIncomingCall, err error) {
	return nil, nil
}

func (se *Server) OutgoingCalls(ctx context.Context, params *protocol.CallHierarchyOutgoingCallsParams) (result []protocol.CallHierarchyOutgoingCall, err error) {
	return nil, nil
}

func (se *Server) SemanticTokensFull(ctx context.Context, params *protocol.SemanticTokensParams) (result *protocol.SemanticTokens, err error) {
	return nil, nil
}

func (se *Server) SemanticTokensFullDelta(ctx context.Context, params *protocol.SemanticTokensDeltaParams) (result interface{}, err error) {
	return nil, nil
}

func (se *Server) SemanticTokensRange(ctx context.Context, params *protocol.SemanticTokensRangeParams) (result *protocol.SemanticTokens, err error) {
	return nil, nil
}

func (se *Server) SemanticTokensRefresh(ctx context.Context) (err error) {
	return nil
}

func (se *Server) LinkedEditingRange(ctx context.Context, params *protocol.LinkedEditingRangeParams) (result *protocol.LinkedEditingRanges, err error) {
	return nil, nil
}

func (se *Server) Moniker(ctx context.Context, params *protocol.MonikerParams) (result []protocol.Moniker, err error) {
	return nil, nil
}

func (se *Server) Request(ctx context.Context, method string, params interface{}) (result interface{}, err error) {
	return nil, nil
}
