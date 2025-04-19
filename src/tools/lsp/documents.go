package lsp

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/tliron/exturl"
	"go.lsp.dev/protocol"
)

const INTERNAL_PATH_PREFIX = "language-server:"

func documentUriToInternalPath(uri protocol.DocumentURI) string {
	return INTERNAL_PATH_PREFIX + string(uri)
}

func getDocument(documentUri protocol.DocumentURI) ([]byte, error) {
	urlContext := exturl.NewContext()
	defer urlContext.Release()

	url := string(documentUri)
	logger.Sugar().Info(url)
	url_, err := urlContext.NewURL(url)
	if err != nil {
		return nil, err
	}

	if reader, err := url_.Open(context.TODO()); err == nil {
		defer reader.Close()
		return io.ReadAll(reader)
	} else {
		return nil, err
	}
}

type TextDocument struct {
	URI     protocol.DocumentURI
	Version int32
	Content string
	Lines   []string
}

// NewTextDocument creates a new text document from content
func NewTextDocument(uri protocol.DocumentURI, version int32, content string) *TextDocument {
	lines := strings.Split(content, "\n")
	return &TextDocument{
		URI:     uri,
		Version: version,
		Content: content,
		Lines:   lines,
	}
}

// LineAt returns the content of the line at the given index
func (d *TextDocument) LineAt(line uint32) string {
	if line < uint32(len(d.Lines)) {
		return d.Lines[line]
	}
	return ""
}

// getDocument retrieves a text document from the server's document storage
func (se *Server) getDocument(uri protocol.DocumentURI) (*TextDocument, error) {
	se.docStore.mu.RLock()
	defer se.docStore.mu.RUnlock()

	doc, exists := se.docStore.documents[uri]
	if !exists {
		return nil, fmt.Errorf("document not found: %s", uri)
	}
	return doc, nil
}

// UpdateDocument adds or updates a document in the store
func (se *Server) UpdateDocument(uri protocol.DocumentURI, version int32, content string) {
	se.docStore.mu.Lock()
	defer se.docStore.mu.Unlock()

	se.docStore.documents[uri] = NewTextDocument(uri, version, content)
}

// RemoveDocument removes a document from the store
func (se *Server) RemoveDocument(uri protocol.DocumentURI) {
	se.docStore.mu.Lock()
	defer se.docStore.mu.Unlock()

	delete(se.docStore.documents, uri)
}
