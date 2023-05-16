package compiler

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/tjarjoura/cc/pkg/ast"
)

/* AMD64 */
var (
	REG_RAX = &Register{map[uint64]string{8: "rax", 4: "eax", 2: "ax", 1: "al"}}
	REG_RBX = &Register{map[uint64]string{8: "rbx", 4: "ebx", 2: "bx", 1: "bl"}}
	REG_RCX = &Register{map[uint64]string{8: "rcx", 4: "ecx", 2: "cx", 1: "cl"}}
	REG_RDX = &Register{map[uint64]string{8: "rdx", 4: "edx", 2: "dx", 1: "dl"}}
	REG_RDI = &Register{map[uint64]string{8: "rdi", 4: "edi", 2: "di", 1: "dil"}}
	REG_RSI = &Register{map[uint64]string{8: "rsi", 4: "esi", 2: "si", 1: "sil"}}
	REG_RSP = &Register{map[uint64]string{8: "rsp", 4: "esp", 2: "sp", 1: "spl"}}
	REG_RBP = &Register{map[uint64]string{8: "rbp", 4: "ebp", 2: "bp", 1: "bpl"}}

	// order that registers will be used in for computing generic
	// expressions
	REG_ORDER = []*Register{REG_RAX, REG_RCX, REG_RDX}
)

type OperandType string

const (
	OP_TYPE_REGISTER  OperandType = "register"
	OP_TYPE_ADDRESS               = "address"
	OP_TYPE_IMMEDIATE             = "immediate"
)

type Register struct {
	NameMap map[uint64]string
}

func (r *Register) String() string { return r.NameMap[8] }

// TODO add binary encoding, right now we are relying on NASM to convert to raw machine code
type Instruction struct {
	neumonic string
	operandA Operand
	operandB Operand
}

type Operand interface {
	String() string
	Size() uint64
	Type() ast.Declaration
	OperandType() OperandType
}

func isImmediate(o Operand) bool {
	return o.OperandType() == OP_TYPE_IMMEDIATE
}

type RegisterOperand struct {
	Register *Register
	DataType ast.Declaration
}

func (r *RegisterOperand) String() string {
	name, ok := r.Register.NameMap[SizeOf(r.Type())]
	if !ok {
		return "???"
	}

	return name
}
func (r *RegisterOperand) Size() uint64             { return SizeOf(r.DataType) }
func (r *RegisterOperand) Type() ast.Declaration    { return r.DataType }
func (r *RegisterOperand) OperandType() OperandType { return OP_TYPE_REGISTER }

type Immediate interface {
	Operand
	immediateOperand()
}

type ImmediateInt struct {
	Value int64
}

func (i *ImmediateInt) immediateOperand()        {}
func (i *ImmediateInt) String() string           { return fmt.Sprintf("0x%x", i.Value) }
func (i *ImmediateInt) Size() uint64             { return IntSize(uint64(i.Value)) }
func (i *ImmediateInt) Type() ast.Declaration    { return IntType(i.Value) }
func (i *ImmediateInt) OperandType() OperandType { return OP_TYPE_IMMEDIATE }

type Address struct {
	Base         *Register
	Scale        *Register
	Index        *Register
	Displacement int64
	DataType     ast.Declaration
}

func (a *Address) Size() uint64             { return SizeOf(a.DataType) }
func (a *Address) Type() ast.Declaration    { return a.DataType }
func (a *Address) OperandType() OperandType { return OP_TYPE_ADDRESS }
func (a *Address) String() string {
	var result string
	if a.Base != nil {
		result = a.Base.String()
	}

	if a.Index != nil {
		var scaledIndex string
		if a.Scale != nil {
			scaledIndex = fmt.Sprintf("%s*%s", a.Scale.String(),
				a.Index.String())
		} else {
			scaledIndex = a.Index.String()
		}

		if result != "" {
			result = fmt.Sprintf("%s + %s", result, scaledIndex)
		} else {
			result = scaledIndex
		}
	}

	if a.Displacement != 0 {
		absDisplacement := a.Displacement
		if a.Displacement < 0 {
			absDisplacement *= -1
		}

		displacement := fmt.Sprintf("0x%s",
			strconv.FormatInt(absDisplacement, 16))

		if result != "" {
			if a.Displacement < 0 {
				result = fmt.Sprintf("%s - %s",
					result, displacement)
			} else {
				result = fmt.Sprintf("%s + %s",
					result, displacement)
			}
		} else {
			if a.Displacement < 0 {
				result = displacement
			} else {
				result = fmt.Sprintf("-%s", displacement)
			}
		}
	}

	sizeMap := map[uint64]string{8: "qword", 4: "dword", 2: "word", 1: "byte"}
	return fmt.Sprintf("%s [%s]", sizeMap[SizeOf(a.DataType)], result)
}

func (i *Instruction) Assembly() string {
	var out strings.Builder
	out.WriteString(fmt.Sprintf("%s", i.neumonic))
	if i.operandA != nil {
		out.WriteString(fmt.Sprintf("\t%s", i.operandA.String()))
		if i.operandB != nil {
			out.WriteString(fmt.Sprintf(", %s", i.operandB.String()))
		}
	}
	return out.String()
}

func Add(opA Operand, opB Operand) *Instruction {
	return &Instruction{neumonic: "add", operandA: opA, operandB: opB}
}

func Leave() *Instruction {
	return &Instruction{neumonic: "leave"}
}

func Mov(opA Operand, opB Operand) *Instruction {
	return &Instruction{neumonic: "mov", operandA: opA, operandB: opB}
}

func Neg(op Operand) *Instruction {
	return &Instruction{neumonic: "neg", operandA: op}
}

func Push(op Operand) *Instruction {
	return &Instruction{neumonic: "push", operandA: op}
}

func Ret() *Instruction {
	return &Instruction{neumonic: "ret"}
}

func Sub(opA Operand, opB Operand) *Instruction {
	return &Instruction{neumonic: "sub", operandA: opA, operandB: opB}
}
