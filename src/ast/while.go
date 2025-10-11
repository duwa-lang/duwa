package ast

import "github.com/duwa-lang/duwa/src/token"

type WhileExpression struct {
	Expression
	Token       token.Token
	Condition   Expression
	Consequence *BlockStatement
}

func (ie *WhileExpression) expressionNode() {}

func (ie *WhileExpression) TokenLiteral() string { return ie.Token.Literal }

// TODO: Print this one properly
func (ie *WhileExpression) String() string {
	return ""
}
