package ast

type Node interface {
	String() string
}

type Statement interface {
	Node
	statementNode()
}

type TranslationUnit struct {
	Node
	Declarations []Declaration
}

func (t *TranslationUnit) String() string { return "" } //TODO
