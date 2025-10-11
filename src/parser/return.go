package parser

import (
	"github.com/duwa-lang/duwa/src/ast"
	"github.com/duwa-lang/duwa/src/token"
)

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.curToken}
	p.nextToken()
	if p.curTokenIs(token.SEMICOLON) {
		stmt.ReturnValue = &ast.NullLiteral{
			Token: p.curToken,
		}
		p.nextToken()
	} else {
		stmt.ReturnValue = p.parseExpression(LOWEST)
		if p.peekTokenIs(token.SEMICOLON) {
			p.nextToken()
		}
	}
	return stmt
}
