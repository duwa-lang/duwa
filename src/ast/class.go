package ast

import (
	"bytes"

	"github.com/sevenreup/duwa/src/token"
)

type ClassDeclStatement struct {
	Expression
	Token token.Token
	Name  *Identifier
	Super *Identifier
	Body  *BlockStatement
}

func (class *ClassDeclStatement) expressionNode()      {}
func (class *ClassDeclStatement) TokenLiteral() string { return class.Token.Literal }
func (class *ClassDeclStatement) String() string {
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
