package lsp

import (
	"sync"

	"github.com/sevenreup/duwa/src/parser"
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
	"github.com/tliron/kutil/problems"
)

var documentStates sync.Map // protocol.DocumentUri to DocumentState

type DocumentState struct {
	Content       string
	ParserContext *parser.Context
	Problems      *problems.Problems

	DocumentURI protocol.DocumentUri
	Symbols     []protocol.SymbolInformation
	Diagnostics []protocol.Diagnostic
}

func getDocumentState(documentUri protocol.DocumentUri) *DocumentState {
	if documentState, ok := documentStates.Load(documentUri); ok {
		return documentState.(*DocumentState)
	} else {
		return nil
	}
}

func validateDocumentState(documentUri protocol.DocumentUri, notify glsp.NotifyFunc) *DocumentState {
	documentState, created := _getOrCreateDocumentState(documentUri)

	if created {
		go notify(protocol.ServerTextDocumentPublishDiagnostics, &protocol.PublishDiagnosticsParams{
			URI:         documentUri,
			Diagnostics: documentState.Diagnostics,
		})
	}

	return documentState
}

func _getOrCreateDocumentState(documentUri protocol.DocumentUri) (*DocumentState, bool) {
	if documentState, ok := documentStates.Load(documentUri); ok {
		return documentState.(*DocumentState), false
	} else {
		documentState := NewDocumentState(documentUri)
		if existing, loaded := documentStates.LoadOrStore(documentUri, documentState); loaded {
			return existing.(*DocumentState), false
		} else {
			return documentState, true
		}
	}
}

func NewDocumentState(documentUri protocol.DocumentUri) *DocumentState {
	self := DocumentState{DocumentURI: documentUri}
	content, err := getDocument(documentUri)
	if err != nil {
		self.Fill()
		return &self
	}
	_ = parserDW.ParseFile([]byte(content))
	self.ParserContext = parserDW.Context
	self.Fill()
	return &self
}

func (docState *DocumentState) Fill() {
	docState.Diagnostics = createDiagnostics(docState.Problems, docState.Content)
	if docState.ParserContext != nil {
		docState.Symbols = createSymbols(docState.ParserContext, docState.Content, docState.DocumentURI)
	}
}
