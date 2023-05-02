package compiler

import (
	"fmt"
	"strings"
)

// TODO add binary encoding, right now we are relying on NASM to convert to raw machine code
type Instruction struct {
	neumonic string
	operandA Operand
	operandB Operand
}

type OperandType string

const (
	OP_REG_DIRECT   OperandType = "OP_REG_DIRECT"
	OP_REG_INDIRECT             = "OP_REG_INDIRECT"
	OP_IMM_DIRECT               = "OP_IMM_DIRECT"
	OP_IMM_INDIRECT             = "OP_IMM_INDIRECT"
)

// CPU Registers
type RegisterName string

const (
	REG_RAX RegisterName = "rax"
	REG_RBX              = "rbx"
	REG_RCX              = "rcx"
	REG_RDX              = "rdx"
	REG_EAX              = "eax"
	REG_EBX              = "ebx"
	REG_ECX              = "ecx"
	REG_EDX              = "edx"
)

type Operand interface {
	String() string
}

type Register struct {
	Register RegisterName
	Pointer  bool
}

func (r *Register) String() string {
	if r.Pointer {
		return fmt.Sprintf("[%s]", r.Register)
	}

	return string(r.Register)
}

type Immediate struct {
	Value   int64
	Pointer bool
}

func (i *Immediate) String() string {
	if i.Pointer {
		return fmt.Sprintf("[0x%x]", i.Value)
	}

	return fmt.Sprintf("0x%x", i.Value)
}

func (i *Instruction) Assembly() string {
	var out strings.Builder
	out.WriteString(fmt.Sprintf("\t%s", i.neumonic))
	if i.operandA != nil {
		out.WriteString(fmt.Sprintf("\t%s", i.operandA.String()))
		if i.operandB != nil {
			out.WriteString(fmt.Sprintf(", %s", i.operandB.String()))
		}
	}
	return out.String()
}

func Mov(opA Operand, opB Operand) *Instruction {
	return &Instruction{neumonic: "mov", operandA: opA, operandB: opB}
}

func Ret() *Instruction {
	return &Instruction{neumonic: "ret"}
}
