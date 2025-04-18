package lsp

import (
	"context"
	"io"

	"github.com/tliron/exturl"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

const INTERNAL_PATH_PREFIX = "language-server:"

func documentUriToInternalPath(uri protocol.DocumentUri) string {
	return INTERNAL_PATH_PREFIX + string(uri)
}

func getDocument(documentUri protocol.DocumentUri) ([]byte, error) {
	urlContext := exturl.NewContext()
	defer urlContext.Release()

	url := string(documentUri)
	log.Info(url)
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
