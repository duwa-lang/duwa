package parser

import (
	"github.com/duwa-lang/duwa/src/ast"
)

func (parser *Parser) continueStatement() ast.Expression {
	return &ast.ContinueStatement{Token: parser.curToken}
}
