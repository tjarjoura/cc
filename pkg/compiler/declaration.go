package compiler

import "github.com/tjarjoura/cc/pkg/ast"

func (f *Function) compileVariableDeclaration(varDecl *ast.VariableDeclaration) {
	t := varDecl.Type()
	f.frameSize += int64(SizeOf(t))
	address := &Address{
		Base:         REG_RBP,
		Displacement: -1 * f.frameSize,
		DataType:     t,
	}

	f.variables[varDecl.Name] = address

	if varDecl.Definition != nil {
		result := f.compileExpression(varDecl.Definition)
		f.Instructions = append(f.Instructions, Mov(address, result))
		if reg, ok := result.(*RegisterOperand); ok {
			f.freeReg(reg.Register)
		}
	}
}
