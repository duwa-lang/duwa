package ast

import (
	"bytes"
	"strings"

	"github.com/duwa-lang/duwa/src/token"
)

type MethodExpression struct {
	Token     token.Token
	Left      Expression
	Method    Expression
	Arguments []Expression
}

func (ie *MethodExpression) expressionNode() {}

func (ie *MethodExpression) TokenLiteral() string { return ie.Token.Literal }

// TODO: Print this one properly
func (ie *MethodExpression) String() string {
	var out bytes.Buffer
	params := []string{}
	for _, p := range ie.Arguments {
		params = append(params, p.String())
	}
	out.WriteString(ie.TokenLiteral())
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") ")
	return out.String()
}
