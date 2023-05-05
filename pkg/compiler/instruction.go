package compiler

import (
	"fmt"
	"strings"
)

/* AMD64 */

// TODO add binary encoding, right now we are relying on NASM to convert to raw machine code
type Instruction struct {
	neumonic string
	operandA Operand
	operandB Operand
}

type Operand interface {
	String() string
	Size() int
}

type Register struct {
	name string
	size int
}

func (r *Register) String() string { return r.name }
func (r *Register) Size() int      { return r.size }

var (
	// 64 bit registers
	REG_RAX = &Register{"rax", 8}
	REG_RBX = &Register{"rbx", 8}
	REG_RCX = &Register{"rcx", 8}
	REG_RDX = &Register{"rdx", 8}
	REG_RDI = &Register{"rdi", 8}
	REG_RSI = &Register{"rsi", 8}
	REG_RSP = &Register{"rsp", 8}
	REG_RBP = &Register{"rbp", 8}
	// ...
	// 32 bit registers
	REG_EAX = &Register{"eax", 4}
	REG_EBX = &Register{"ebx", 4}
	REG_ECX = &Register{"ecx", 4}
	REG_EDX = &Register{"edx", 4}
	// ...
	// TODO add more as needed
)

type ImmediateInt struct {
	Value int64
}

func (i *ImmediateInt) String() string { return fmt.Sprintf("0x%x", i.Value) }
func (i *ImmediateInt) Size() int      { return IntSize(uint64(i.Value)) }

type Address struct {
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

func Leave() *Instruction {
	return &Instruction{neumonic: "leave"}
}

func Ret() *Instruction {
	return &Instruction{neumonic: "ret"}
}
