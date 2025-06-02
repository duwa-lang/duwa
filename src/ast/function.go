package ast

import (
	"bytes"
	"strings"

	"github.com/sevenreup/duwa/src/token"
)

type FunctionDeclStatement struct {
	Token      token.Token // The 'ndondomeko' token
	Name       *Identifier // The name of the function
	Parameters []*Identifier
	Body       *BlockStatement
}

func (fl *FunctionDeclStatement) expressionNode()      {}
func (fl *FunctionDeclStatement) TokenLiteral() string { return fl.Token.Literal }
func (fl *FunctionDeclStatement) String() string {
	var out bytes.Buffer
	params := []string{}
	for _, p := range fl.Parameters {
		params = append(params, p.String())
	}
	out.WriteString(fl.TokenLiteral())
	out.WriteString(fl.Name.TokenLiteral())
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") ")
	out.WriteString(fl.Body.String())
	return out.String()
}
