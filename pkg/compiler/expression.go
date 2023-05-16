package compiler

import (
	"fmt"

	"github.com/tjarjoura/cc/pkg/ast"
)

func (f *Function) compileInfixExpression(inf *ast.InfixExpression) Operand {
	leftE := f.compileExpression(inf.Left)
	rightE := f.compileExpression(inf.Right)

	if leftE == nil || rightE == nil {
		return nil
	}

	op, ok := f.infixOperations[inf.Operator]
	if !ok {
		f.err(fmt.Sprintf(
			"Can not handle %s operators!", inf.Operator))
		return nil
	}

	if isImmediate(leftE) && isImmediate(rightE) {
		if op.CompileImmediate == nil {
			f.err(fmt.Sprintf(
				"Cannot handle %s operator for immediate operands",
				inf.Operator))
			return nil
		}
		return op.CompileImmediate(inf.Operator,
			leftE.(Immediate), rightE.(Immediate))
	}

	if op.CompileRuntime == nil {
		f.err(fmt.Sprintf("Cannot handle %s operators at runtime!",
			inf.Operator))
		return nil
	}

	return op.CompileRuntime(inf.Operator, leftE, rightE)
}

func (f *Function) compilePrefixExpression(p *ast.PrefixExpression) Operand {
	rightOp := f.compileExpression(p.Right)
	if rightOp == nil {
		return nil
	}

	op, ok := f.prefixOperations[p.Operator]
	if !ok {
		f.err(fmt.Sprintf(
			"can not handle prefix operator '%s'", p.Operator))
		return nil
	}

	if rightOp.OperandType() == OP_TYPE_IMMEDIATE {
		if op.CompileImmediate == nil {
			f.err(fmt.Sprintf(
				"cannot handle prefix operator '%s' at compile time",
				p.Operator))
			return nil
		}

		return op.CompileImmediate(p.Operator, rightOp.(Immediate))
	}

	if op.CompileRuntime == nil {
		f.err(fmt.Sprintf(
			"cannot handle prefix operator '%s' at runtime",
			p.Operator))
		return nil
	}
	return op.CompileRuntime(p.Operator, rightOp)
}

/*
*
Generate machine instructions for the expression and return an operand
representing either the result or where the result will be stored
*
*/
func (f *Function) compileExpression(expr ast.Expression) Operand {
	switch e := expr.(type) {
	case *ast.InfixExpression:
		return f.compileInfixExpression(e)
	case *ast.PrefixExpression:
		return f.compilePrefixExpression(e)
	case *ast.Identifier:
		return f.variables[e.Value]
	case *ast.IntegerLiteral:
		return &ImmediateInt{Value: e.Value}
	}

	return nil
}
