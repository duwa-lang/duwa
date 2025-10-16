package ast

import (
	"bytes"

	"github.com/duwa-lang/duwa/src/token"
)

type VariableDeclarationStatement struct {
	Statement
	Type       token.Token // the token.Nambala token
	Identifier *Identifier
	Value      Expression
}

func (ls *VariableDeclarationStatement) TokenLiteral() string { return ls.Identifier.TokenLiteral() }

func (ls *VariableDeclarationStatement) String() string {
	var out bytes.Buffer
	out.WriteString(ls.Type.Literal + " ")
	out.WriteString(ls.Identifier.Value)
	out.WriteString(" = ")
	if ls.Value != nil {
		out.WriteString(ls.Value.String())
	}
	out.WriteString(";")
	return out.String()
}
