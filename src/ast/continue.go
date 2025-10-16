package ast

import "github.com/duwa-lang/duwa/src/token"

type ContinueStatement struct {
	Token token.Token
	Expression
}

func (node *ContinueStatement) TokenLiteral() string { return node.Token.Literal }
func (node *ContinueStatement) String() string {
	return "pitirizani"
}
