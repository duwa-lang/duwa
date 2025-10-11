package parser

import (
	"github.com/duwa-lang/duwa/src/ast"
	"github.com/duwa-lang/duwa/src/token"
)

func (parser *Parser) mapLiteral() ast.Expression {
	mapLiteral := &ast.MapExpression{Token: parser.curToken}
	mapLiteral.Pairs = make(map[ast.Expression]ast.Expression)

	for !parser.peekTokenIs(token.OPENING_BRACE) {
		if parser.peekTokenIs(token.CLOSING_BRACE) {
			break
		}
		parser.nextToken()

		key := parser.parseExpression(LOWEST)

		if !parser.expectPeek(token.COLON) {
			return nil
		}

		parser.nextToken()

		value := parser.parseExpression(LOWEST)

		mapLiteral.Pairs[key] = value

		if !parser.peekTokenIs(token.CLOSING_BRACE) && !parser.peekTokenIs(token.COMMA) {
			return nil
		}
		if parser.peekTokenIs(token.CLOSING_BRACE) {
			break
		}

		parser.nextToken()
	}

	if !parser.expectPeek(token.CLOSING_BRACE) {
		return nil
	}

	return mapLiteral
}
