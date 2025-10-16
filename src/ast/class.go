package ast

import (
	"bytes"

	"github.com/duwa-lang/duwa/src/token"
)

type ClassStatement struct {
	Expression
	Token token.Token
	Name  *Identifier
	Super *Identifier
	Body  *BlockStatement
}

func (class *ClassStatement) expressionNode()      {}
func (class *ClassStatement) TokenLiteral() string { return class.Token.Literal }
func (class *ClassStatement) String() string {
	var out bytes.Buffer
	out.WriteString("ndondomeko ")
	out.WriteString(class.Name.String())
	if class.Super != nil {
		out.WriteString(" ndi ")
		out.WriteString(class.Super.String())
	}
	out.WriteString(" {\n")
	out.WriteString(class.Body.String())
	out.WriteString("\n}")
	return out.String()
}
