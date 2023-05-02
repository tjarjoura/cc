package compiler

import "github.com/tjarjoura/cc/pkg/ast"

func (f *Function) compileStatement(stmt ast.Statement) {
	switch s := stmt.(type) {
	case *ast.ReturnStatement:
		// TODO always returning zero for now...
		f.Instructions = append(f.Instructions,
			Mov(&Register{Register: REG_EAX, Pointer: false},
				&Immediate{Value: 0, Pointer: false}))
		f.Instructions = append(f.Instructions, Ret())
	}
}
