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

func (c *Compiler) Errors() []CompileError { return c.errors }

type CompileError struct {
	token token.Token
	msg   string
	warn  bool
}

func (c *CompileError) String() string {
	var prefix = "[ERROR]"
	if c.warn {
		prefix = "[WARN]"
	}
	return fmt.Sprintf("%s %s", prefix, c.msg)
}

type CompilationObject interface {
	Assembly() string
	Errors() []CompileError
}

type Function struct {
	Instructions []*Instruction
	Type         ast.Declaration

	errors     []CompileError
	registers  map[*Register]bool
	operations map[string]Operation
}

func NewFunction(t ast.Declaration) *Function {
	fn := &Function{Type: t}
	fn.registerOperations()
	return fn
}

func (f *Function) Errors() []CompileError { return f.errors }

func (f *Function) err(msg string) {
	f.errors = append(f.errors, CompileError{msg: msg, warn: false})
}

func (f *Function) warn(msg string) {
	f.errors = append(f.errors, CompileError{msg: msg, warn: true})
}

func (f *Function) Assembly() string {
	var out strings.Builder

	for _, instr := range f.Instructions {
		out.WriteString("\t")
		out.WriteString(instr.Assembly())
		out.WriteString("\n")
	}
	return out.String()
}

func (f *Function) allocNextReg() *Register {
	for _, reg := range REG_ORDER {
		if !f.registers[reg] {
			f.registers[reg] = true
			return reg
		}
	}

	panic("ran out of registers to use! (will fix this in the future)")
}

func (f *Function) allocReg(r *Register) { f.registers[r] = true }

func (f *Function) freeReg(r *Register) { f.registers[r] = false }

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

func (c *Compiler) WriteAssembly(w io.StringWriter) error {
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
			sections[TEXT].WriteString(o.Assembly())
		}
	}

	for section, data := range sections {
		if _, err := w.WriteString(fmt.Sprintf("SECTION %s\n", section)); err != nil {
			return err
		}

		if _, err := w.WriteString(data.String()); err != nil {
			return err
		}
	}

	return nil
}

func (c *Compiler) compileFunction(fnDecl *ast.FunctionDeclaration) {
	f := NewFunction(fnDecl.Type())
	c.symbolMap[fnDecl.Name] = f

	if fnDecl.Body != nil {
		for _, stmt := range fnDecl.Body.Statements {
			f.compileStatement(stmt)
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
