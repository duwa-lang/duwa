package lsp

import (
	"context"

	"go.lsp.dev/protocol"
)

func (se *Server) Definition(ctx context.Context, params *protocol.DefinitionParams) (result []protocol.Location, err error) {
	logger.Sugar().Debug("Definition", params.TextDocument.URI, params.Position)

	return nil, nil
}
