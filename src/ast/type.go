package ast

import "github.com/sevenreup/duwa/src/types"

type TypeDeclStmt struct {
	Name string
	Type types.Type
}

func (t *TypeDeclStmt) Accept(visitor StmtVisitor) {
	visitor.VisitTypeDecl(t)
}
