package compiler

import (
	"fmt"
	"strings"

	"github.com/tjarjoura/cc/pkg/ast"
)

var (
	SizeToType = map[uint64]string{
		1: "char",
		2: "short",
		4: "int",
		8: "long",
	}

	TypeToSize = map[string]uint64{
		"char":          1,
		"short int":     2,
		"int":           4,
		"long int":      8,
		"long long int": 8,
		"float":         4,
		"double":        8,
		"long double":   16,
	}

	PtrSize uint64 = 8
)

func IntSize(val uint64) uint64 {
	if val <= 0xFF {
		return 1
	} else if val <= 0xFFFF {
		return 2
	} else if val <= 0xFFFFFFFF {
		return 4
	} else {
		return 8
	}
}

// Return the smallest integer type that can fit the value
func IntType(val int64) *ast.BaseType {
	size := IntSize(uint64(val))
	return &ast.BaseType{Name: SizeToType[size]}
}

func SizeOf(decl ast.Declaration) uint64 {
	switch d := decl.(type) {
	case *ast.Pointer:
		return PtrSize
	case *ast.Array:
		panic("Sizeof(array) not implemented yet!")
	case *ast.BaseType:
		return TypeToSize[d.Name]
	case *ast.VariableDeclaration:
		return SizeOf(d.Type())
	case *ast.FunctionDeclaration:
		// This is a nonsense operation, but other compilers seem to return 1 here
		return 1
	default:
		// Should never get here
		return 0
	}
}

func isPointer(d ast.Declaration) bool {
	_, ok := d.(*ast.Pointer)
	return ok
}

func isFloat(d ast.Declaration) bool {
	b, ok := d.(*ast.BaseType)
	if !ok {
		return false
	}

	return strings.Contains(b.Name, "double") ||
		strings.Contains(b.Name, "float")

}

func biggestType(a ast.Declaration, b ast.Declaration) ast.Declaration {
	if SizeOf(b) > SizeOf(a) {
		return b
	}

	return a
}

func (f *Function) compileTypeConversion(toType ast.Declaration,
	fromType ast.Declaration, value Operand) Operand {
	if ast.ConvertError(toType, fromType) {
		f.err(fmt.Sprintf(
			"incompatible types when converting from %s to %s",
			fromType.String(), toType.String()))
		return nil
	} else if ast.ConvertWarn(toType, fromType) {
		f.warn(fmt.Sprintf(
			"converting from %s to %s without a cast",
			fromType.String(), toType.String()))
		return nil
	}

	return value

	sizeDifference := SizeOf(fromType) - SizeOf(toType)
	if sizeDifference == 0 {
		return value
	} else {
		panic(fmt.Sprintf("can't handle type conversions: %s to %s!",
			fromType.String(), toType.String()))
	}
	switch value.(type) {
	case *ImmediateInt: // TODO
		if sizeDifference > 0 { // getting smaller
		} else { // getting bigger
		}

	}

	/*
		sizeDifference := SizeOf(fromType) - SizeOf(toType)
		if !(fromType.Unsigned) {

		}
		if SizeOf(toType) == SizeOf(fromType) {
			// nothing to do
			return value
		} else if SizeOf(toType) < SizeOf(fromType) {
			switch v := value.(type) {
			case *ImmediateInt:
				return &ImmediateInt{Value: value & Bitmask(SizeOf(toType))}
			case *Register:
				f.Instructions = append(f.Instructions,
					Shl(v),
					Sar(v),
				)
			default:
				panic("Can't handle operands of this type!")
			}
		} else {
		}
	*/

	return value
}
