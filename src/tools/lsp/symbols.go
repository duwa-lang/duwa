package lsp

import (
	"github.com/sevenreup/duwa/src/parser"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

func createSymbols(context *parser.Context, content string, documentUri protocol.DocumentUri) []protocol.SymbolInformation {
	var symbols []protocol.SymbolInformation

	return symbols
}
