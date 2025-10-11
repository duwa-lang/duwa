package parser

import (
	"fmt"

	"github.com/duwa-lang/duwa/src/ast"
	"github.com/shopspring/decimal"
)

func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{Token: p.curToken}
	value, err := decimal.NewFromString(p.curToken.Literal)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}
	lit.Value = value
	return lit
}
