package ast

import "github.com/duwa-lang/duwa/src/token"

type Identifier struct {
	AssignmentNode
	Token token.Token // the token.IDENT token
	Value string
}

func (i *Identifier) String() string { return i.Value }

func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
