package parser

import (
	"github.com/duwa-lang/duwa/src/ast"
	"github.com/duwa-lang/duwa/src/token"
)

func (p *Parser) parseStatement() ast.Statement {
	if p.curToken.Type == token.RETURN {
		return p.parseReturnStatement()
	}

	if token.LookupVariableType(p.curToken.Type) != "" || (p.curToken.Type == token.IDENT && p.peekTokenIs(token.IDENT)) {
		return p.parseVariableDeclarationStatement()
	}

	statement := p.parseAssignmentStatement()

	if statement != nil {
		return statement
	}

	return p.parseExpressionStatement()
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.curToken}
	stmt.Expression = p.parseExpression(LOWEST)
	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	return stmt
}
