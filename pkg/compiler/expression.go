package compiler

import "github.com/tjarjoura/cc/pkg/ast"

/*
Generate machine instructions for the expression and return an operand representing either the result or where the result will be stored
*/
func (f *Function) compileExpression(expr ast.Expression) (Operand, ast.Declaration) {
	switch e := expr.(type) {
	case *ast.Identifier:
		return nil, nil
	case *ast.IntegerLiteral:
		return &ImmediateInt{Value: e.Value}, IntType(e.Value)
	}

	return nil, nil
}
