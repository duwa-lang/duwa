package lsp

import (
	parserpkg "github.com/sevenreup/duwa/src/parser"
	"github.com/tliron/commonlog"
)

var log = commonlog.GetLogger("puccini-language-server.tosca")
var parserDW = parserpkg.NewParser()
