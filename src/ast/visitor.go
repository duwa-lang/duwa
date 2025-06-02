package ast

type StmtVisitor interface {
	VisitVarDecl(*VariableDeclStatement)
	VisitAssignment(*AssigmentStatement)
	VisitFuncDecl(*FunctionDeclStatement)
	VisitReturn(*ReturnStatement)
	VisitTypeDecl(*TypeDeclStmt)
	VisitStructDecl(*ClassDeclStatement)
	// TODO: Implement enums VisitEnumDecl(*EnumDeclStmt)
}
