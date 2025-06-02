package ast

import (
	"bytes"

	"github.com/sevenreup/duwa/src/types"
)

type VariableDeclStatement struct {
	Statement
	Type       types.Type // the token.Nambala token
	Identifier *Identifier
	Value      Expression
}

func (ls *VariableDeclStatement) TokenLiteral() string { return ls.Identifier.TokenLiteral() }

func (ls *VariableDeclStatement) String() string {
	var out bytes.Buffer
	out.WriteString(ls.Type.String() + " ")
	out.WriteString(ls.Identifier.Value)
	out.WriteString(" = ")
	if ls.Value != nil {
		out.WriteString(ls.Value.String())
	}
	out.WriteString(";")
	return out.String()
}
