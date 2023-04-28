package ast

type Node interface {
	String() string
}

type TranslationUnit struct {
	Node
	DeclarationStatements []*DeclarationStatement
}

func (t *TranslationUnit) String() string { return "" } //TODO
