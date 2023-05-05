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

	if varDecl.VarType.String() != expectedType {
		t.Errorf("expected stmt.Type=%s, got=%s", expectedType,
			varDecl.VarType.String())
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
		{"int int int int int x;", "x", "", "int"},
		{"static int int int int int long x;", "x", "static", "long int"},
		{"long x;", "x", "", "long int"},
		{"long unsigned unsigned long x;", "x", "", "long long int"},
		{"double double x;", "x", "", "double"},
		{"double long double x;", "x", "", "long double"},
		{"void void x;", "x", "", "void"},
		{"short short x;", "x", "", "short int"},
		{"int longname;", "longname", "", "int"},
		{"char ch;", "ch", "", "char"},
		{"long long int x;", "x", "", "long long int"},
		{"long const long int x;", "x", "", "const long long int"},
		{"volatile short x;", "x", "", "volatile short int"},
		{"static volatile short x;", "x", "static", "volatile short int"},
		{"extern int *x;", "x", "extern", "(int) *"},
		{"register int *x;", "x", "register", "(int) *"},
		{"int **x;", "x", "", "((int) *) *"},
		{"const long int **x;", "x", "", "((const long int) *) *"},
		{"const int **const x;", "x", "", "((const int) *) * const"},
		{"const int *volatile const const const*const x;", "x", "", "((const int) * const volatile) * const"},
		{"int x[2];", "x", "", "(int)[2]"},
		{"int x[2+7-3];", "x", "", "(int)[((2 + 7) - 3)]"},
		{"int x[(ident + 2)];", "x", "", "(int)[(ident + 2)]"},
		{"int x[2][3];", "x", "", "((int)[3])[2]"},
		{"int *x[2][3];", "x", "", "(((int) *)[3])[2]"},
		{"int (*x)[2][3];", "x", "", "(((int)[3])[2]) *"},
		{"int (*x[2])[3];", "x", "", "(((int)[3]) *)[2]"},
		{";;;int (*fptr)();", "fptr", "", "(int ()) *"},
		{"int (*fptr)(int);", "fptr", "", "(int (int)) *"},
		{"char* (*(*foo[5])(char *))[];", "foo", "", "(((((char) *)[]) * ((char) *)) *)[5]"},
	}

	for _, tt := range tests {
		p := New(lexer.New(tt.input))

		tUnit := p.Parse()
		checkErrors(t, p)

		if len(tUnit.DeclarationStatements) != 1 {
			t.Fatalf("expected %d declaration statements, got=%d", 1,
				len(tUnit.DeclarationStatements))
		}

		if len(tUnit.DeclarationStatements[0].Declarations) != 1 {
			t.Fatalf("expected %d declarations, got=%d", 1,
				len(tUnit.DeclarationStatements))
		}

		decl := tUnit.DeclarationStatements[0].Declarations[0]
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
		{"z", "(long int)[3]"},
		{"a", "(long int) * const"},
	}

	p := New(lexer.New(input))
	tUnit := p.Parse()
	checkErrors(t, p)

	if len(tUnit.DeclarationStatements) != 1 {
		t.Fatalf("expected len(tUnit.DeclarationStatements)=1, got=%d\n",
			len(tUnit.DeclarationStatements))
	}

	decls := tUnit.DeclarationStatements[0].Declarations
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

func TestParseExpression(t *testing.T) {
	tests := []struct {
		input        string
		expectedExpr string
	}{
		{"int x = 3;", "3"},
		{"int x = 3 + 4;", "(3 + 4)"},
		{"int x = 3 + 4 + 5;", "((3 + 4) + 5)"},
		{"int x = 3 + 4 * 5;", "(3 + (4 * 5))"},
		{"int x = 3/ 9 + 4 * 5;", "((3 / 9) + (4 * 5))"},
		{"int x = 3|4&5-6>>4<<2*9+4/2%8^4<2>9==10!=11;",
			"(3 | ((4 & (((5 - 6) >> 4) << ((2 * 9) + ((4 / 2) % 8)))) ^ ((((4 < 2) > 9) == 10) != 11)))"},
		{"int x = 3.0;", "3.0"},
		{"int x = ~7 * 3;", "((~7) * 3)"},
		{"int x = ~(7 * 3);", "(~(7 * 3))"},
		{"int x = !~*&x*99;", "((!(~(*(&x)))) * 99)"},
	}

	for _, test := range tests {
		p := New(lexer.New(test.input))
		tUnit := p.Parse()
		checkErrors(t, p)

		declStmts := tUnit.DeclarationStatements
		if len(declStmts) != 1 {
			t.Fatalf("expected len(declStmts)=%d, got=%d\n", 1, len(declStmts))
		}

		if len(declStmts[0].Declarations) != 1 {
			t.Fatalf("expected len(declStmts[0].Declarations)=%d, got=%d\n", 1,
				len(declStmts[0].Declarations))
		}

		varDecl, ok := declStmts[0].Declarations[0].(*ast.VariableDeclaration)
		if !ok {
			t.Fatalf("expected decl to be *ast.VariableDeclaration, got=%T",
				declStmts[0].Declarations)
		}

		if varDecl.Definition == nil {
			t.Fatalf("variable definition was nil")
		}

		if varDecl.Definition.String() != test.expectedExpr {
			t.Fatalf("expected definition to be=%s, got=%s",
				test.expectedExpr, varDecl.Definition.String())
		}
	}
}

func TestParseFunctionDeclaration(t *testing.T) {
	tests := []struct {
		input      string
		expectedFn string
	}{
		{"int f();", "int f()"},
		{"int f(int x);", "int f(int x)"},
		{"int f(int *x);", "int f((int) * x)"},
		{"int f(int*);", "int f((int) *)"},
		{"int f(int **x);", "int f(((int) *) * x)"},
		{"int f(int, int );", "int f(int, int)"},
		{"int *f(int, int );", "(int) * f(int, int)"},
		{"int **f(int, int );", "((int) *) * f(int, int)"},
		{"int ***f();", "(((int) *) *) * f()"},
		{"char (*(*func())[5])();", "(((char ()) *)[5]) * func()"},
	}

	for _, tt := range tests {
		p := New(lexer.New(tt.input))
		tUnit := p.Parse()
		checkErrors(t, p)

		declStmts := tUnit.DeclarationStatements
		if len(declStmts) != 1 {
			t.Fatalf("expected len(declStmts)=%d, got=%d\n", 1, len(declStmts))
		}

		if len(declStmts[0].Declarations) != 1 {
			t.Fatalf("expected len(declStmts[0].Declarations)=%d, got=%d\n", 1,
				len(declStmts[0].Declarations))
		}

		fnDecl, ok := declStmts[0].Declarations[0].(*ast.FunctionDeclaration)
		if !ok {
			t.Fatalf("expected decl to be *ast.FunctionDeclaration, got=%T",
				declStmts[0].Declarations[0])
		}

		if fnDecl.String() != tt.expectedFn {
			t.Fatalf("expected fnDecl to=%s, got=%s", tt.expectedFn, fnDecl.String())
		}
	}
}

// make sure different syntax errors don't crash the program and are handled gracefully
func TestParseErrors(t *testing.T) {
	tests := []struct {
		input string
	}{
		{"int nosemicolon"},
		{"int fn(incomplete"},
		{"int (incomplete"},
		{"faketype g;"},
		{"int x = 3 +;"},
		{"int = 3;"},
		{"int x[3-] ;"},
		{"int x[3 ;"},
		{"int (*x)( ;"},
		{"int (((()))*x);"},
		{"int *const 3;"},
		{"int a, b(){};"},
	}

	for _, tt := range tests {
		p := New(lexer.New(tt.input))
		p.Parse()

		if len(p.Errors()) == 0 {
			t.Fatalf("Expected errors from %s but got 0.", tt.input)
		}
	}
}

func TestParseFunctionDefinition(t *testing.T) {
	input := `
int main(int argc, char **argv) {
	int x=2, y;
	char z=72-6*8, *h;
	return z*2;
}`

	testStmts := []string{
		"int x = 2, int y;",
		"char z = (72 - (6 * 8)), (char) * h;",
		"return (z * 2);",
	}

	p := New(lexer.New(input))
	tUnit := p.Parse()
	checkErrors(t, p)

	declStmts := tUnit.DeclarationStatements
	if len(declStmts) != 1 {
		t.Fatalf("expected len(declStmts)=%d, got=%d\n", 1, len(declStmts))
	}

	decls := declStmts[0].Declarations
	if len(decls) != 1 {
		t.Fatalf("expected len(decls)=%d, got=%d\n", 1, len(decls))
	}

	fnDecl, ok := declStmts[0].Declarations[0].(*ast.FunctionDeclaration)
	if !ok {
		t.Fatalf("expected decl to be *ast.FunctionDeclaration, got=%T",
			declStmts[0].Declarations[0])
	}

	if len(fnDecl.Parameters) != 2 {
		t.Fatalf("expected 2 parameters")
	}

	if fnDecl.Body == nil {
		t.Fatalf("expected fnDecl.Body != nil")
	}

	stmts := fnDecl.Body.Statements
	if len(stmts) != len(testStmts) {
		t.Fatalf("expected len(fnDecl.Body.Statements=1, got=%d",
			len(fnDecl.Body.Statements))
	}

	for i, stmt := range testStmts {
		if stmt != stmts[i].String() {
			t.Fatalf("expected %s=%s", testStmts[i], stmts[i])
		}
	}
}
