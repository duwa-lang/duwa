package lsp

import (
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

func WorkspaceDidRenameFiles(context *glsp.Context, params *protocol.RenameFilesParams) error {
	return nil
}
