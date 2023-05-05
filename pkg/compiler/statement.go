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

		// TODO move to its own funtion
		if ast.ConvertError(f.fnType, returnType) {
			f.err(fmt.Sprintf(
				"incompatible types when converting from %s to %s",
				returnType.String(), f.fnType.String()))
			return
		} else if ast.ConvertWarn(f.fnType, returnType) {
			f.warn(fmt.Sprintf(
				"converting from %s to %s without a cast",
				returnType.String(), f.fnType.String()))
			return
		}

		var returnReg = REG_EAX // default 32 bit register for return val
		if SizeOf(f.fnType) > 4 {
			returnReg = REG_RAX // need 64 bit register to hold return val
		}

		//returnValue = scaleValue(f.fnType, returnValue)

		f.Instructions = append(f.Instructions,
			Mov(returnReg, returnValue),
			Leave(),
			Ret())
	}
}
