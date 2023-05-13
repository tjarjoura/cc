package compiler

import (
	"fmt"

	"github.com/tjarjoura/cc/pkg/ast"
)

func (f *Function) compileStatement(stmt ast.Statement) {
	switch s := stmt.(type) {
	case *ast.ReturnStatement:
		returnValue := f.compileExpression(s.ReturnValue)
		if returnValue == nil {
			// TODO handle gracefully
			panic(fmt.Sprintf("Could not compile %s!", s.ReturnValue))
		}

		returnValue = f.compileTypeConversion(f.Type, returnValue.Type(),
			returnValue)
		if returnValue == nil { // error
			return
		}

		// TODO if return value is already in RAX, no need to mov()
		returnReg := &RegisterOperand{Register: REG_RAX, _type: f.Type}
		f.Instructions = append(f.Instructions,
			Mov(returnReg, returnValue),
			//Leave(),
			Ret())
	}
}
