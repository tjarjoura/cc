package compiler

import (
	"fmt"
	"log"

	"github.com/tjarjoura/cc/pkg/token"
)

type (
	infixImmediateCompileFn  func(string, Immediate, Immediate) Immediate
	infixRuntimeCompileFn    func(string, Operand, Operand) Operand
	prefixImmediateCompileFn func(string, Immediate) Immediate
	prefixRuntimeCompileFn   func(string, Operand) Operand
)

type InfixOperation struct {
	CompileImmediate infixImmediateCompileFn
	CompileRuntime   infixRuntimeCompileFn
}

type PrefixOperation struct {
	CompileImmediate prefixImmediateCompileFn
	CompileRuntime   prefixRuntimeCompileFn
}

func (f *Function) compilePrefixArithmetic(operator string, operand Operand,
) Operand {
	var resultReg Operand = operand
	if operand.OperandType() != OP_TYPE_REGISTER {
		resultReg = &RegisterOperand{Register: f.allocNextReg(),
			DataType: operand.Type()}
		f.Instructions = append(f.Instructions, Mov(resultReg, operand))
	}

	switch operator {
	case token.MINUS:
		f.Instructions = append(f.Instructions, Neg(resultReg))
	default:
		f.err(fmt.Sprintf(
			"cannot handle prefix operator '%s' at runtime",
			operator))
		return nil
	}

	return resultReg
}

func (f *Function) compilePrefixArithmeticImm(operator string, operand Immediate,
) Immediate {
	log.Println("called compilePrefixArithmeticImm")
	if isFloat(operand.Type()) {
		f.err("floating point arithmetic is not supported")
		return nil
	}

	val := operand.(*ImmediateInt).Value

	switch operator {
	case token.MINUS:
		return &ImmediateInt{Value: -1 * val}
	}

	f.err(fmt.Sprintf(
		"cannot handle prefix operator '%s' at compile time",
		operator))
	return nil
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
	case token.PLUS:
		f.Instructions = append(f.Instructions, Add(resultReg, b))
	case token.MINUS:
		f.Instructions = append(f.Instructions, Sub(resultReg, b))
	default:
		f.err(fmt.Sprintf(
			"cannot handle infix operator %s at runtime", op))
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
	f.infixOperations = map[string]InfixOperation{
		token.PLUS:     {compileArithmeticImm, f.compileArithmetic},
		token.MINUS:    {compileArithmeticImm, f.compileArithmetic},
		token.ASTERISK: {compileArithmeticImm, f.compileArithmetic},
		token.SLASH:    {compileArithmeticImm, f.compileArithmetic},
	}

	f.prefixOperations = map[string]PrefixOperation{
		token.MINUS: {f.compilePrefixArithmeticImm, f.compilePrefixArithmetic},
	}
}
