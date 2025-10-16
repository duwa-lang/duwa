package parser

import (
	"github.com/duwa-lang/duwa/src/ast"
)

func (parser *Parser) parseCompoundExpression(left ast.Expression) ast.Expression {
	compound := &ast.Compound{
		Token:    parser.curToken,
		Operator: parser.curToken.Literal,
		Left:     left,
	}

	precedence := parser.curPrecedence()

	parser.nextToken()

	compound.Right = parser.parseExpression(precedence)

	return compound
}
