package parser

import (
	"github.com/sevenreup/duwa/src/ast"
	"github.com/sevenreup/duwa/src/token"
)

func (p *Parser) classStatement() ast.Expression {
	class := &ast.ClassDeclStatement{Token: p.curToken}

	p.nextToken()

	class.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	// TODO: Implement inheritance

	if !p.expectPeek(token.OPENING_BRACE) {
		return nil
	}

	class.Body = p.parseBlockStatement()

	return class
}
