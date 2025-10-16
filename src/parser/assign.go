package parser

import (
	"github.com/duwa-lang/duwa/src/ast"
	"github.com/duwa-lang/duwa/src/token"
)

func (p *Parser) parseAssignmentStatement() *ast.AssigmentStatement {
	statement := &ast.AssigmentStatement{}

	if p.curTokenIs(token.IDENT) {
		statement.Identifier = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	}

	if !p.peekTokenIs(token.ASSIGN) {
		return nil
	}

	p.nextToken()
	p.nextToken()

	statement.Value = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return statement
}

func (p *Parser) handleIndexAssigment(indexExp *ast.IndexExpression) *ast.AssigmentStatement {
	statement := &ast.AssigmentStatement{
		Identifier: indexExp.Left,
	}

	p.nextToken()

	statement.Identifier = indexExp

	if !p.curTokenIs(token.ASSIGN) {
		return nil
	}

	p.nextToken()

	statement.Value = p.parseExpression(LOWEST)

	p.previousIndex = nil

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return statement
}

func (p *Parser) handlePropertyAssignment(propertyExp *ast.PropertyExpression) *ast.AssigmentStatement {
	statement := &ast.AssigmentStatement{
		Identifier: propertyExp,
	}

	p.nextToken()

	if !p.curTokenIs(token.ASSIGN) {
		return nil
	}

	p.nextToken()

	statement.Value = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return statement
}
