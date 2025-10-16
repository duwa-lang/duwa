package parser

import (
	"github.com/duwa-lang/duwa/src/ast"
	"github.com/duwa-lang/duwa/src/token"
)

func (p *Parser) parseIndexExpression(left ast.Expression) ast.Expression {
	exp := &ast.IndexExpression{Token: p.curToken, Left: left}

	p.nextToken()

	exp.Index = p.parseExpression(LOWEST)

	if !p.expectPeek(token.CLOSING_BRACKET) {
		return nil
	}

	if p.peekTokenIs(token.ASSIGN) {
		return p.handleIndexAssigment(exp)
	}

	p.previousIndex = exp

	return exp
}
