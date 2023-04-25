package ast

import "github.com/tjarjoura/cc/pkg/token"

type Expression interface {
	Node
	expressionNode()
}

type InfixExpression struct {
}

type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (i *IntegerLiteral) expressionNode() {}
func (i *IntegerLiteral) String() string  { return i.Token.Literal }

type FloatLiteral struct {
	Token token.Token
	Value float64
}

func (fp *FloatLiteral) expressionNode() {}
func (fp *FloatLiteral) String() string  { return fp.Token.Literal }
