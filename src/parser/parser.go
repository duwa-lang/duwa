package parser

import (
	"fmt"

	"github.com/sevenreup/duwa/src/ast"
	"github.com/sevenreup/duwa/src/lexer"
	"github.com/sevenreup/duwa/src/token"
)

type (
	prefixParseFn   func() ast.Expression
	infixParseFn    func(ast.Expression) ast.Expression
	postfixParserFn func() ast.Expression
)

const (
	_ int = iota
	LOWEST
	LOGICAL_AND // &&
	LOGICAL_OR  // ||
	EQUALS
	// ==
	LESSGREATER // > or <
	SUM
	// +
	PRODUCT
	PREFIX
	CALL
	INDEX // array[index]
)

var precedences = map[token.TokenType]int{
	token.EQUAL_TO:                 EQUALS,
	token.NOT_EQUAL_TO:             EQUALS,
	token.LESS_THAN:                LESSGREATER,
	token.LESS_THAN_OR_EQUAL_TO:    LESSGREATER,
	token.GREATER_THAN:             LESSGREATER,
	token.GREATER_THAN_OR_EQUAL_TO: LESSGREATER,
	token.PLUS:                     SUM,
	token.MINUS:                    SUM,
	token.ASTERISK:                 PRODUCT,
	token.SLASH:                    PRODUCT,
	token.AND_AND:                  LOGICAL_AND,
	token.OR_OR:                    LOGICAL_OR,
	token.OPENING_PAREN:            CALL,
	token.OPENING_BRACKET:          INDEX,
	token.FULL_STOP:                INDEX,
	token.PLUS_EQUAL:               SUM,
	token.MINUS_EQUAL:              SUM,
	token.STAR_EQUAL:               PRODUCT,
	token.SLASH_EQUAL:              PRODUCT,
}

// Pratt parser
type Parser struct {
	l *lexer.Lexer

	errors []string

	previousToken token.Token
	curToken      token.Token
	peekToken     token.Token

	previousIndex *ast.IndexExpression

	prefixParseFns   map[token.TokenType]prefixParseFn
	infixParseFns    map[token.TokenType]infixParseFn
	postfixParserFns map[token.TokenType]postfixParserFn

	imports []*ast.ImportStatement

	Context *Context
}

func NewParser() *Parser {
	p := &Parser{l: lexer.NewLexel(), errors: []string{}, imports: []*ast.ImportStatement{}}

	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	p.registerPrefix(token.INT, p.parseIntegerLiteral)
	p.registerPrefix(token.STR, p.parseStringLiteral)
	p.registerPrefix(token.NULL, p.nullLiteral)

	p.registerPrefix(token.IDENT, p.parseIdentifier)
	p.registerPrefix(token.BANG, p.parsePrefixExpression)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)
	p.registerPrefix(token.TRUE, p.parseBoolean)
	p.registerPrefix(token.FALSE, p.parseBoolean)
	p.registerPrefix(token.IF, p.parseIfExpression)
	p.registerPrefix(token.FOR, p.parseForExpression)
	p.registerPrefix(token.WHILE, p.parseWhileExpression)
	p.registerPrefix(token.FUNCTION, p.parseFunctionLiteral)
	p.registerPrefix(token.CLASS, p.classStatement)
	p.registerPrefix(token.BREAK, p.breakStatement)
	p.registerPrefix(token.CONTINUE, p.continueStatement)
	p.registerPrefix(token.IMPORT, p.importStatement)

	p.registerPrefix(token.OPENING_BRACE, p.mapLiteral)
	p.registerPrefix(token.OPENING_PAREN, p.parseGroupedExpression)
	p.registerPrefix(token.OPENING_BRACKET, p.parseArrayLiteral)

	p.infixParseFns = make(map[token.TokenType]infixParseFn)
	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.SLASH, p.parseInfixExpression)
	p.registerInfix(token.ASTERISK, p.parseInfixExpression)
	p.registerInfix(token.FULL_STOP, p.dotExpression)
	p.registerInfix(token.EQUAL_TO, p.parseInfixExpression)
	p.registerInfix(token.NOT_EQUAL_TO, p.parseInfixExpression)
	p.registerInfix(token.GREATER_THAN, p.parseInfixExpression)
	p.registerInfix(token.GREATER_THAN_OR_EQUAL_TO, p.parseInfixExpression)
	p.registerInfix(token.LESS_THAN, p.parseInfixExpression)
	p.registerInfix(token.LESS_THAN_OR_EQUAL_TO, p.parseInfixExpression)
	p.registerInfix(token.OPENING_PAREN, p.parseCallExpression)
	p.registerInfix(token.OPENING_BRACKET, p.parseIndexExpression)
	p.registerInfix(token.PLUS_EQUAL, p.parseCompoundExpression)
	p.registerInfix(token.MINUS_EQUAL, p.parseCompoundExpression)
	p.registerInfix(token.STAR_EQUAL, p.parseCompoundExpression)
	p.registerInfix(token.SLASH_EQUAL, p.parseCompoundExpression)
	p.registerInfix(token.AND_AND, p.parseInfixExpression)
	p.registerInfix(token.OR_OR, p.parseInfixExpression)

	// Register all of our postfix parse functions
	p.postfixParserFns = make(map[token.TokenType]postfixParserFn)
	p.registerPostfix(token.PLUS_PLUS, p.parsePostfixExpression)
	p.registerPostfix(token.MINUS_MINUS, p.parsePostfixExpression)

	return p
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("%d:%d: syntax error: expected next token to be %s, got %s instead", p.peekToken.Pos.Line, p.peekToken.Pos.Column, t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

func (p *Parser) nextToken() {
	p.previousToken = p.curToken
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) ParseFile(value []byte) *ast.File {
	p.l.Handle(value)

	p.nextToken()
	p.nextToken()

	file := &ast.File{}
	file.Statements = []ast.Statement{}
	file.Imports = []ast.ImportStatement{}

	for p.curToken.Type != token.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			file.Statements = append(file.Statements, stmt)
			importStmt, ok := stmt.(*ast.ImportStatement)
			if ok {
				file.Imports = append(file.Imports, *importStmt)
			}
		}
		p.nextToken()

	}

	return file
}

func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		p.peekError(t)
		return false
	}
}

func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

func (parser *Parser) registerPostfix(tokenType token.TokenType, fn postfixParserFn) {
	parser.postfixParserFns[tokenType] = fn
}

func (p *Parser) noPrefixParseFnError(t token.TokenType) {
	msg := fmt.Sprintf("no prefix parse function for %s found", t)
	p.errors = append(p.errors, msg)
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}
	p.nextToken()
	expression.Right = p.parseExpression(PREFIX)
	return expression
}

func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.curToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left,
	}
	precedence := p.curPrecedence()
	p.nextToken()
	expression.Right = p.parseExpression(precedence)
	return expression
}
