package ast

import (
	"fmt"
	"strings"
)

type Statement interface {
	Node
	statementNode()
}

type BlockStatement struct {
	Statements []Statement
}

func (b *BlockStatement) statementNode() {}
func (b *BlockStatement) String() string {
	stmtStrs := []string{}

	for _, stmt := range b.Statements {
		stmtStrs = append(stmtStrs, stmt.String())
	}

	return strings.Join(stmtStrs, "\n")
}

type DeclarationStatement struct {
	Declarations []Declaration
}

func (d *DeclarationStatement) statementNode() {}
func (d *DeclarationStatement) String() string {
	declStrs := []string{}
	for _, decl := range d.Declarations {
		declStrs = append(declStrs, decl.String())
	}

	return fmt.Sprintf("%s;", strings.Join(declStrs, ", "))
}

type ExpressionStatement struct {
	Expression Expression
}

func (e *ExpressionStatement) statementNode() {}
func (e *ExpressionStatement) String() string {
	return fmt.Sprintf("%s;", e.Expression.String())
}

type ReturnStatement struct {
	ReturnValue Expression
}

func (r *ReturnStatement) statementNode() {}
func (r *ReturnStatement) String() string {
	return fmt.Sprintf("return %s;", r.ReturnValue.String())
}

type IfStatement struct{}
type WhileStatement struct{}
type DoWhileStatement struct{}
type ForStatement struct{}
type SwitchStatement struct{}
