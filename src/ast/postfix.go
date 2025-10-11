package ast

import (
	"github.com/duwa-lang/duwa/src/token"
)

type PostfixExpression struct {
	Token    token.Token
	Operator string
}

func (pe *PostfixExpression) TokenLiteral() string { return pe.Token.Literal }

// TODO: Print this one properly
func (pe *PostfixExpression) String() string {
	return ""
}
