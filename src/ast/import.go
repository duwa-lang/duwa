package ast

import "github.com/sevenreup/duwa/src/token"

type ImportType string

const (
	DefaultImport = "default"
	NamedImport   = "named"
)

type ImportExpression struct {
	Expression
	Token token.Token

	Type ImportType

	Module       *StringLiteral
	DefaultAlias *StringLiteral
	Exports      map[string]Identifier
}
