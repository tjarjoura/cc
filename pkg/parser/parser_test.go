package parser

import (
	"testing"

	"github.com/tjarjoura/cc/pkg/ast"
	"github.com/tjarjoura/cc/pkg/lexer"
)

func checkErrors(t *testing.T, p *Parser) {
	for _, err := range p.Errors() {
		t.Errorf("parser error: %s", err.String())
	}
}

func TestParseVariableDeclaration(t *testing.T) {
	tests := []struct {
		input              string
		expectedName       string
		expectedStorage    string
		expectedType       string
		expectedDefinition string
	}{
		{"int x;", "x", "", "int", ""},
		{"char ch;", "ch", "", "char", ""},
		{"long long int x;", "x", "", "long long int", ""},
		{"long const long int x;", "x", "", "const long long int", ""},
		{"volatile short x;", "x", "", "volatile short", ""},
		{"static volatile short x;", "x", "static", "volatile short", ""},
		{"extern int *x;", "x", "extern", "(int) *", ""},
		{"int **x;", "x", "", "((int) *) *", ""},
		{"const long int **x;", "x", "", "((const long int) *) *", ""},
		{"const int **const x;", "x", "", "((const int) *) * const", ""},
		{"const int *volatile const const const*const x;", "x", "", "((const int) * const volatile) * const", ""},
		{"int x[2];", "x", "", "(int)[]", ""},
		{"int x[2][3];", "x", "", "((int)[])[]", ""},
	}

	for _, tt := range tests {
		p := New(lexer.New(tt.input))

		tUnit := p.Parse()
		checkErrors(t, p)

		if len(tUnit.Declarations) != 1 {
			t.Fatalf("expected %d declarations, got=%d", 1,
				len(tUnit.Declarations))
		}

		decl := tUnit.Declarations[0]
		varDecl, ok := decl.(*ast.VariableDeclaration)
		if !ok {
			t.Fatalf("expected decl to be *ast.VariableDeclaration, got=%T",
				decl)
		}

		if varDecl.TypeSpec.String() != tt.expectedType {
			t.Fatalf("expected stmt.Type=%s, got=%s", tt.expectedType,
				varDecl.String())
		}

		if varDecl.Name != tt.expectedName {
			t.Fatalf("expected stmt.Name=%s, got=%s", tt.expectedName,
				varDecl.Name)
		}

		if varDecl.StorageClass != tt.expectedStorage {
			t.Fatalf("expected StorageClass=%s, got=%s",
				tt.expectedStorage, varDecl.StorageClass)
		}
	}
}
