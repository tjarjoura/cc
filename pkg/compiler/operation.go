package compiler

import (
	"fmt"

	"github.com/tjarjoura/cc/pkg/token"
)

type (
	immediateCompileFn func(string, Immediate, Immediate) Immediate
	runtimeCompileFn   func(string, Operand, Operand) Operand
)

type Operation struct {
	CompileImmediate immediateCompileFn
	CompileRuntime   runtimeCompileFn
}

func (f *Function) compileArithmetic(op string, a Operand, b Operand) Operand {
	if isPointer(a.Type()) || isPointer(b.Type()) {
		f.err("pointer arithmetic is not supported yet")
		return nil
	} else if isFloat(a.Type()) || isFloat(b.Type()) {
		f.err("foating point arithmetic is not supported yet")
		return nil
	}

	// find the register we are storing the result in
	var resultReg *RegisterOperand
	if regA, ok := a.(*RegisterOperand); ok {
		resultReg = regA
	} else if regB, ok := b.(*RegisterOperand); ok {
		resultReg = regB
		b = a
	} else {
		resultReg = &RegisterOperand{Register: f.allocNextReg(),
			DataType: biggestType(a.Type(), b.Type())}
		f.Instructions = append(f.Instructions, Mov(resultReg, a))
	}

	switch op {
	case token.MINUS:
		f.Instructions = append(f.Instructions, Sub(resultReg, b))
	default:
		f.err(fmt.Sprintf("Can't support operator '%s' yet", op))
		return nil
	}

	if regB, ok := b.(*RegisterOperand); ok {
		f.freeReg(regB.Register)
	}

	return resultReg
}

func compileArithmeticImm(op string, a Immediate, b Immediate) Immediate {
	immA, ok := a.(*ImmediateInt)
	if !ok {
		panic("compileArithmeticImm can only be called with *ImmediateInt!")
	}

	immB, ok := b.(*ImmediateInt)
	if !ok {
		panic("compileArithmeticImm can only be called with *ImmediateInt!")
	}

	switch op {
	case token.PLUS:
		return &ImmediateInt{Value: immA.Value + immB.Value}
	case token.MINUS:
		return &ImmediateInt{Value: immA.Value - immB.Value}
	case token.ASTERISK:
		return &ImmediateInt{Value: immA.Value * immB.Value}
	case token.SLASH:
		return &ImmediateInt{Value: immA.Value / immB.Value}
	default:
		panic(fmt.Sprintf(
			"Called compileArithmeticImm with unexpected operator!: %s",
			op))

	}

	return nil
}

func (f *Function) registerOperations() {
	f.operations = map[string]Operation{
		token.PLUS:     {compileArithmeticImm, f.compileArithmetic},
		token.MINUS:    {compileArithmeticImm, f.compileArithmetic},
		token.ASTERISK: {compileArithmeticImm, f.compileArithmetic},
		token.SLASH:    {compileArithmeticImm, f.compileArithmetic},
	}
}
