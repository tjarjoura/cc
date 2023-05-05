package compiler

import "github.com/tjarjoura/cc/pkg/ast"

var (
	SizeToType = map[int]string{
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

func IntSize(val uint64) int {
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
