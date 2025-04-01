package parser

import (
	"fmt"

	"github.com/sevenreup/duwa/src/ast"
	"github.com/sevenreup/duwa/src/token"
)

func (p *Parser) importStatement() ast.Expression {
	statement := &ast.ImportStatement{Token: p.curToken, Exports: make(map[string]ast.Identifier)}

	p.nextToken()

	if p.curToken.Type == token.IDENT {
		statement.Type = ast.DefaultImport
		statement.DefaultAlias = &ast.StringLiteral{Token: p.curToken, Value: p.curToken.Literal}
		p.nextToken()
	} else if p.curTokenIs(token.OPENING_BRACE) {
		p.nextToken()
		statement.Type = ast.NamedImport
		for p.curTokenIs(token.IDENT) {
			identifier := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
			statement.Exports[identifier.Value] = *identifier
			if p.peekTokenIs(token.CLOSING_BRACE) {
				p.nextToken()
				p.nextToken()
				break
			} else if p.peekTokenIs(token.COMMA) {
				p.nextToken()
				p.nextToken()
			} else {
				p.errors = append(p.errors, fmt.Sprintf("expected kuchokera after tenga but got %s", p.curToken.Literal))
				return nil
			}
		}
	} else {

		return nil
	}

	if p.curToken.Type == token.FROM && p.expectPeek(token.STR) {
		statement.Module = &ast.StringLiteral{Token: p.curToken, Value: p.curToken.Literal}
	} else {
		p.errors = append(p.errors, fmt.Sprintf("expected kuchokera after tenga but got %s", p.curToken.Literal))
		return nil
	}

	return statement
}
