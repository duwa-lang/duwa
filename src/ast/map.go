package ast

import (
	"github.com/duwa-lang/duwa/src/token"
)

type MapExpression struct {
	Expression
	Token token.Token
	Pairs map[Expression]Expression
}

func (ie *MapExpression) TokenLiteral() string { return ie.Token.Literal }

// TODO: Print this one properly
func (ie *MapExpression) String() string {
	return ""
}
