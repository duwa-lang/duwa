package lsp

import (
	"sync"

	"go.lsp.dev/protocol"
)

type SymbolDefinition struct {
	Name     string
	Kind     protocol.SymbolKind
	Location protocol.Location
}

type SymbolIndex struct {
	mu        sync.RWMutex
	symbols   map[string][]SymbolDefinition
	docIndex  map[protocol.DocumentURI][]string
}

func NewSymbolIndex() *SymbolIndex {
	return &SymbolIndex{
		symbols:  make(map[string][]SymbolDefinition),
		docIndex: make(map[protocol.DocumentURI][]string),
	}
}