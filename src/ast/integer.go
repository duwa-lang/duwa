package ast

import (
	"github.com/duwa-lang/duwa/src/token"
	"github.com/shopspring/decimal"
)

type IntegerLiteral struct {
	Token token.Token
	Value decimal.Decimal
}

func (il *IntegerLiteral) expressionNode()      {}
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Literal }
func (il *IntegerLiteral) String() string       { return il.Token.Literal }
