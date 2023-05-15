package compiler

import (
	"fmt"

	"github.com/tjarjoura/cc/pkg/ast"
)

func (f *Function) compileStatement(stmt ast.Statement) {
	switch s := stmt.(type) {
	case *ast.DeclarationStatement:
		for _, d := range s.Declarations {
			switch decl := d.(type) {
			case *ast.VariableDeclaration:
				f.compileVariableDeclaration(decl)
			case *ast.FunctionDeclaration:
				f.err("cannot declare a function inside another function")
			default:
				f.warn("declaration does not declare anything")
			}
		}
	case *ast.ReturnStatement:
		returnValue := f.compileExpression(s.ReturnValue)
		if returnValue == nil {
			f.err(fmt.Sprintf("Could not compile '%s'", s.ReturnValue))
			return
		}

		returnValue = f.compileTypeConversion(f.Type, returnValue.Type(),
			returnValue)
		if returnValue == nil { // error
			return
		}

		// TODO if return value is already in RAX, no need to mov()
		returnReg := &RegisterOperand{Register: REG_RAX, DataType: f.Type}
		f.Instructions = append(f.Instructions,
			Mov(returnReg, returnValue),
			Leave(),
			Ret())
	}
}
