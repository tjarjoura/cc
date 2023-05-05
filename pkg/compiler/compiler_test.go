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

	for _, co := range c.symbolMap {
		for _, err := range co.Errors() {
			t.Errorf("compiler error: %s", err.String())
		}
	}

	return ret
}

func checkStackFrameInstructions(t *testing.T, fn *Function, frameSize int) int {
	stackFrameInstructions := []string{}
	for i := 0; i < len(stackFrameInstructions); i++ {
		iStr := fn.Instructions[i].Assembly()
		if iStr != stackFrameInstructions[i] {
			t.Errorf("expected %s=%s\n", iStr,
				stackFrameInstructions[i])
			return -1
		}
	}

	return len(stackFrameInstructions)
}

func TestReturnStatement(t *testing.T) {
	tests := []struct {
		input                string
		expectedInstructions []string
	}{
		{"int f() { return 3; }", []string{"mov\teax, 0x3"}},
		{"int f() { return 4000; }", []string{"mov\teax, 0xfa0"}},
		{"long f() { return 4000; }", []string{"mov\trax, 0xfa0"}},
		//{"char f() { return 0xF40; }", []string{"mov\teax, 0x40"}},
	}

	expectedReturnInstructions := []string{"leave", "ret"}

	for _, tt := range tests {
		p := parser.New(lexer.New(tt.input))
		tUnit := p.Parse()
		if !checkParserErrors(t, p) {
			t.FailNow()
		}

		c := New(tUnit)
		c.Compile()
		if !checkCompilerErrors(t, c) {
			t.FailNow()
		}

		fn := c.symbolMap["f"].(*Function)

		offset := checkStackFrameInstructions(t, fn, 0)
		if offset < 0 {
			t.FailNow()
		}

		expectedLen := offset + len(expectedReturnInstructions) + len(tt.expectedInstructions)
		if len(fn.Instructions) != expectedLen {
			t.Fatalf("Expected there to be %d instructions, got=%d",
				expectedLen, len(fn.Instructions))
		}

		expectedInstructions := append(tt.expectedInstructions,
			expectedReturnInstructions...)
		for i, expected := range expectedInstructions {
			actual := fn.Instructions[i+offset].Assembly()
			if actual != expected {
				t.Fatalf("expected %s=%s", actual, expected)
			}
		}
	}
}

func TestAssemblyOut(t *testing.T) {
	tests := []struct {
		input       string
		expectedAsm string
	}{
		{
			input: `
int main() {
	return 0;
}
`,
			expectedAsm: `SECTION .text
GLOBAL main
main:
	mov	eax, 0x0
	leave
	ret
`},
		{
			input: `
int main() {
	return 7;
}
`,
			expectedAsm: `SECTION .text
GLOBAL main
main:
	mov	eax, 0x7
	leave
	ret
`},
		/*
		   		{
		   			input: `
		   int a = 3, b;
		   int main() {
		   	return 7;
		   }
		   `,
		   			expectedAsm: `SECTION .text
		   GLOBAL main
		   main:
		   	mov	eax, 0x7
		   	leave
		   	ret
		   SECTION .data
		   a: 	dw 0x7
		   SECTION .bss
		   b: 	resw
		   `},
		*/
	}

	for _, tt := range tests {
		p := parser.New(lexer.New(tt.input))
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

		if asmOut.String() != tt.expectedAsm {
			t.Fatalf("expected:\n%s\nactual:\n%s\n", tt.expectedAsm,
				asmOut.String())
		}
	}
}
