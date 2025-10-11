package parser

import "github.com/duwa-lang/duwa/src/ast"

func (p *Parser) nullLiteral() ast.Expression {
	return &ast.NullLiteral{Token: p.curToken}
}
