package ast

import "github.com/duwa-lang/duwa/src/token"

type PropertyExpression struct {
	Token    token.Token
	Left     Expression
	Property Expression
}

func (pe *PropertyExpression) expressionNode() {}

func (pe *PropertyExpression) TokenLiteral() string { return pe.Token.Literal }

func (pe *PropertyExpression) String() string {
	return pe.Left.String() + "." + pe.Property.String()
}
