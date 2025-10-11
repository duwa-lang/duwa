package ast

import "github.com/duwa-lang/duwa/src/token"

type NullLiteral struct {
	Token token.Token
}

func (il *NullLiteral) TokenLiteral() string { return il.Token.Literal }
func (il *NullLiteral) String() string       { return il.Token.Literal }
