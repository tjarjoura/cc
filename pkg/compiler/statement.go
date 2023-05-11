package compiler

import (
	"fmt"

	"github.com/tjarjoura/cc/pkg/ast"
)

func (f *Function) compileStatement(stmt ast.Statement) {
	switch s := stmt.(type) {
	case *ast.ReturnStatement:
		returnValue, returnType := f.compileExpression(s.ReturnValue)
		if returnValue == nil {
			// TODO handle gracefully
			panic(fmt.Sprintf("Could not compile %s!", s.ReturnValue))
		}

		returnValue = f.compileTypeConversion(f.fnType, returnType, returnValue)
		if returnValue == nil { // error
			return
		}

		var returnReg = REG_EAX // default 32 bit register for return val
		if SizeOf(f.fnType) > 4 {
			returnReg = REG_RAX // need 64 bit register to hold return val
		}

		//returnValue = scaleValue(f.fnType, returnValue)

		f.Instructions = append(f.Instructions,
			Mov(returnReg, returnValue),
			//Leave(),
			Ret())
	}
}
