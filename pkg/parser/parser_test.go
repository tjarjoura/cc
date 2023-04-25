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

func testVariableDeclaration(t *testing.T, decl ast.Declaration,
	expectedType string, expectedName string, expectedStorage string) bool {

	varDecl, ok := decl.(*ast.VariableDeclaration)
	if !ok {
		t.Errorf("expected decl to be *ast.VariableDeclaration, got=%T",
			decl)
		return false
	}

	if varDecl.TypeSpec.String() != expectedType {
		t.Errorf("expected stmt.Type=%s, got=%s", expectedType,
			varDecl.TypeSpec.String())
		return false
	}

	if varDecl.Name != expectedName {
		t.Errorf("expected stmt.Name=%s, got=%s", expectedName,
			varDecl.Name)
		return false
	}

	if varDecl.StorageClass != expectedStorage {
		t.Errorf("expected StorageClass=%s, got=%s",
			expectedStorage, varDecl.StorageClass)
		return false
	}

	return true
}

func TestParseVariableDeclaration(t *testing.T) {
	tests := []struct {
		input           string
		expectedName    string
		expectedStorage string
		expectedType    string
	}{
		{"int x;", "x", "", "int"},
		{"char ch;", "ch", "", "char"},
		{"long long int x;", "x", "", "long long int"},
		{"long const long int x;", "x", "", "const long long int"},
		{"volatile short x;", "x", "", "volatile short"},
		{"static volatile short x;", "x", "static", "volatile short"},
		{"extern int *x;", "x", "extern", "(int) *"},
		{"int **x;", "x", "", "((int) *) *"},
		{"const long int **x;", "x", "", "((const long int) *) *"},
		{"const int **const x;", "x", "", "((const int) *) * const"},
		{"const int *volatile const const const*const x;", "x", "", "((const int) * const volatile) * const"},
		{"int x[2];", "x", "", "(int)[]"},
		{"int x[2][3];", "x", "", "((int)[])[]"},
		{"int *x[2][3];", "x", "", "(((int) *)[])[]"},
		{"int (*x)[2][3];", "x", "", "(((int)[])[]) *"},
		{"int (*x[2])[3];", "x", "", "(((int)[]) *)[]"},
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
		testVariableDeclaration(t, decl, tt.expectedType,
			tt.expectedName, tt.expectedStorage)
	}
}

func TestParseMultipleDeclarations(t *testing.T) {
	input := "static long int x=5, *y, z[3], *const a;"

	expectedDecls := []struct {
		name     string
		typeSpec string
	}{
		{"x", "long int"},
		{"y", "(long int) *"},
		{"z", "(long int)[]"},
		{"a", "(long int) * const"},
	}

	p := New(lexer.New(input))
	tUnit := p.Parse()
	decls := tUnit.Declarations
	checkErrors(t, p)

	if len(decls) != len(expectedDecls) {
		t.Fatalf("expected len(decls)=%d, got=%d\n",
			len(expectedDecls), len(decls))
	}

	for i, decl := range decls {
		if !testVariableDeclaration(t, decl, expectedDecls[i].typeSpec,
			expectedDecls[i].name, "static") {
			t.FailNow()
		}
	}
}
