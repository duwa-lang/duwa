package ast

import (
	"testing"

	"github.com/duwa-lang/duwa/src/token"
)

func TestString(t *testing.T) {
	file := &File{
		Statements: []Statement{
			&VariableDeclarationStatement{
				Type: token.Token{Type: token.INTEGER, Literal: "nambala"},
				Identifier: &Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "myVar"},
					Value: "myVar",
				},
				Value: &Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "anotherVar"},
					Value: "anotherVar",
				},
			},
		},
	}
	if file.String() != "nambala myVar = anotherVar;" {
		t.Errorf("file.String() wrong. got=%q", file.String())
	}
}
