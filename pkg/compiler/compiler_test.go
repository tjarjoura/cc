package compiler

import (
	"strings"
	"testing"

	"github.com/tjarjoura/cc/pkg/lexer"
	"github.com/tjarjoura/cc/pkg/parser"
)

func checkParserErrors(t *testing.T, p *parser.Parser) bool {
	var ret = true
	for _, err := range p.Errors() {
		t.Errorf("parser error: %s", err.String())
		ret = false
	}

	return ret
}

func checkCompilerErrors(t *testing.T, c *Compiler) bool {
	var ret = true
	for _, err := range c.errors {
		t.Errorf("compiler error: %s", err.String())
	}

	return ret
}

func TestAssemblyOut(t *testing.T) {
	input := `
int main() {
	return 0;
}
`
	expectedAsm := `SECTION .text
GLOBAL main
main:
	mov	eax, 0x0
	ret
`

	p := parser.New(lexer.New(input))
	tUnit := p.Parse()
	if !checkParserErrors(t, p) {
		t.FailNow()
	}

	c := New(tUnit)
	c.Compile()
	if !checkCompilerErrors(t, c) {
		t.FailNow()
	}

	asmOut := &strings.Builder{}
	c.WriteAssembly(asmOut)

	if asmOut.String() != expectedAsm {
		t.Fatalf("expected:\n%s\nactual:\n%s\n", expectedAsm,
			asmOut.String())
	}
}
