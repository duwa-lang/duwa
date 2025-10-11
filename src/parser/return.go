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
		// Don't advance past the semicolon - leave curToken at the last token of the statement
		// The caller (parseBlockStatement) will advance to the next statement
	} else {
		stmt.ReturnValue = p.parseExpression(LOWEST)
		if p.peekTokenIs(token.SEMICOLON) {
			p.nextToken()
		}
	}
	return stmt
}
