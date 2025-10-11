package parser

import (
	"github.com/duwa-lang/duwa/src/ast"
	"github.com/duwa-lang/duwa/src/token"
)

func (parser *Parser) dotExpression(left ast.Expression) ast.Expression {
	currentToken := parser.curToken
	currentPrecedence := parser.curPrecedence()

	parser.nextToken()

	if parser.peekTokenIs(token.OPENING_PAREN) {
		// Method
		expression := &ast.MethodExpression{Token: currentToken, Left: left}
		expression.Method = parser.parseExpression(currentPrecedence)

		parser.nextToken()

		expression.Arguments = parser.parseExpressionList(token.CLOSING_PAREN)

		return expression
	}

	// Property
	return &ast.PropertyExpression{Token: parser.curToken, Left: left, Property: parser.parseExpression(currentPrecedence)}
}
