package parser

import (
	"github.com/duwa-lang/duwa/src/ast"
)

func (parser *Parser) breakStatement() ast.Expression {
	return &ast.BreakStatement{Token: parser.curToken}
}
