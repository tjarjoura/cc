package compiler

import (
	"fmt"
	"go/token"
	"io"
	"strings"

	"github.com/tjarjoura/cc/pkg/ast"
)

type Compiler struct {
	translationUnit *ast.TranslationUnit
	symbolMap       map[string]CompilationObject
	functions       []*Function
	registers       []bool

	errors []CompileError
}

type CompileError struct {
	token token.Token
	msg   string
}

func (c *CompileError) String() string {
	return c.msg
}

type CompilationObject interface {
	assembly() string
}

type Function struct {
	Instructions []*Instruction
}

func (f *Function) assembly() string {
	var out strings.Builder

	for _, instr := range f.Instructions {
		out.WriteString(instr.Assembly())
		out.WriteString("\n")
	}
	return out.String()
}

type Variable struct {
	size    int
	initial []byte
}

const (
	TEXT = ".text"
	DATA = ".data"
	BSS  = ".bss"
)

// Needs to generate the raw bytes itself and also the relocation entries
func (v *Variable) generateBinary() []byte {
	return v.initial
}

func New(tUnit *ast.TranslationUnit) *Compiler {
	compiler := &Compiler{translationUnit: tUnit, symbolMap: map[string]CompilationObject{}}
	return compiler
}

func (c *Compiler) WriteAssembly(w io.StringWriter) {
	sections := map[string]*strings.Builder{
		TEXT: &strings.Builder{},
		//DATA: &strings.Builder{},
	}
	for symbol, obj := range c.symbolMap {
		switch o := obj.(type) {
		case *Function:
			// TODO only if storage class != "static"
			sections[TEXT].WriteString(fmt.Sprintf("GLOBAL %s\n", symbol))
			sections[TEXT].WriteString(fmt.Sprintf("%s:\n", symbol))
			sections[TEXT].WriteString(o.assembly())
		}
	}

	for section, data := range sections {
		w.WriteString(fmt.Sprintf("SECTION %s\n", section))
		w.WriteString(data.String())
	}
}

func (c *Compiler) compileFunction(fnDecl *ast.FunctionDeclaration) {
	fn := &Function{}
	c.symbolMap[fnDecl.Name] = fn

	if fnDecl.Body != nil {
		for _, stmt := range fnDecl.Body.Statements {
			fn.compileStatement(stmt)
		}
	}
}

func (c *Compiler) Compile() {
	for _, declStmt := range c.translationUnit.DeclarationStatements {
		for _, decl := range declStmt.Declarations {
			switch d := decl.(type) {
			case *ast.VariableDeclaration:
				continue
			case *ast.FunctionDeclaration:
				c.compileFunction(d)
			}
		}
	}
}
