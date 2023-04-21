package lexer

import (
	"testing"

	"github.com/tjarjoura/cc/pkg/token"
)

func TestLexer(t *testing.T) {
	input := `/* A nonsensical program that tries to use every language feature of C */
short d;
long e = 100L;
const long long f;
unsigned g;

typedef struct x {
	union {
		int g;
		short h;
	} y;
} xx;

static int fnc(int x, int y) {
	static long u=3;
	return e + u * x / y;
}

enum {
	A = 0,
	B
};

int main(void) { 
	auto double x = 40.23e+20;
	register float y = 40e+20;
	volatile int z = 0x3000;
	volatile int zz = 03000;
	static char f6 = 'c';
	z ^= 1;
	char nine[] = "nine";
	char *nine2 = nine;
	nine[2] = f6;
	(*nine2)++;
	*nine2 = 3;
	x += 2 >> 2&3^5<<2|9;
	do {
		z &= 4 * 99 / 6 + 7 - 0;
		z = f6 > 3 ? z : 0;
	} while (0);

	if (z < 1) {
		z /= 2;
	} else if (z <= 2) {
		z %= 2;
	} else if (z > 3) {
		z -= 2;
	} else if (z >= 3) {
		z *= 2;
	} else if (z && --d) {
		z <<= 1;
		z >>= 1;
		static struct x l;
		z = l.y.g++;
		struct x *m;
		z = --m->y.h;
	} else if (z || d) {}
	else if (z != d) {}
	else if (z == d) {}

	for (; ;) { if (!x) continue; else break; }

	switch (f6) {
		case '\n':
			goto exit;
		default:
			break;
	}

exit:
	// return /* */
	return fnc(z + 3, z - 2); 
}/* end comment */`
	expectedTokens := []token.Token{
		{Type: token.SHORT, Literal: "short", Line: 2, Column: 1},
		{Type: token.IDENTIFIER, Literal: "d", Line: 2, Column: 7},
		{Type: token.SEMICOLON, Literal: ";", Line: 2, Column: 8},
		{Type: token.LONG, Literal: "long", Line: 3, Column: 1},
		{Type: token.IDENTIFIER, Literal: "e", Line: 3, Column: 6},
		{Type: token.ASSIGN, Literal: "=", Line: 3, Column: 8},
		{Type: token.NUMERICL, Literal: "100L", Line: 3, Column: 10},
		{Type: token.SEMICOLON, Literal: ";", Line: 3, Column: 14},
		{Type: token.CONST, Literal: "const", Line: 4, Column: 1},
		{Type: token.LONG, Literal: "long", Line: 4, Column: 7},
		{Type: token.LONG, Literal: "long", Line: 4, Column: 12},
		{Type: token.IDENTIFIER, Literal: "f", Line: 4, Column: 17},
		{Type: token.SEMICOLON, Literal: ";", Line: 4, Column: 18},
		{Type: token.UNSIGNED, Literal: "unsigned", Line: 5, Column: 1},
		{Type: token.IDENTIFIER, Literal: "g", Line: 5, Column: 10},
		{Type: token.SEMICOLON, Literal: ";", Line: 5, Column: 11},
		{Type: token.TYPEDEF, Literal: "typedef", Line: 7, Column: 1},
		{Type: token.STRUCT, Literal: "struct", Line: 7, Column: 9},
		{Type: token.IDENTIFIER, Literal: "x", Line: 7, Column: 16},
		{Type: token.LBRACE, Literal: "{", Line: 7, Column: 18},
		{Type: token.UNION, Literal: "union", Line: 8, Column: 2},
		{Type: token.LBRACE, Literal: "{", Line: 8, Column: 8},
		{Type: token.INT, Literal: "int", Line: 9, Column: 3},
		{Type: token.IDENTIFIER, Literal: "g", Line: 9, Column: 7},
		{Type: token.SEMICOLON, Literal: ";", Line: 9, Column: 8},
		{Type: token.SHORT, Literal: "short", Line: 10, Column: 3},
		{Type: token.IDENTIFIER, Literal: "h", Line: 10, Column: 9},
		{Type: token.SEMICOLON, Literal: ";", Line: 10, Column: 10},
		{Type: token.RBRACE, Literal: "}", Line: 11, Column: 2},
		{Type: token.IDENTIFIER, Literal: "y", Line: 11, Column: 4},
		{Type: token.SEMICOLON, Literal: ";", Line: 11, Column: 5},
		{Type: token.RBRACE, Literal: "}", Line: 12, Column: 1},
		{Type: token.IDENTIFIER, Literal: "xx", Line: 12, Column: 3},
		{Type: token.SEMICOLON, Literal: ";", Line: 12, Column: 5},
		{Type: token.STATIC, Literal: "static", Line: 14, Column: 1},
		{Type: token.INT, Literal: "int", Line: 14, Column: 8},
		{Type: token.IDENTIFIER, Literal: "fnc", Line: 14, Column: 12},
		{Type: token.LPAREN, Literal: "(", Line: 14, Column: 15},
		{Type: token.INT, Literal: "int", Line: 14, Column: 16},
		{Type: token.IDENTIFIER, Literal: "x", Line: 14, Column: 20},
		{Type: token.COMMA, Literal: ",", Line: 14, Column: 21},
		{Type: token.INT, Literal: "int", Line: 14, Column: 23},
		{Type: token.IDENTIFIER, Literal: "y", Line: 14, Column: 27},
		{Type: token.RPAREN, Literal: ")", Line: 14, Column: 28},
		{Type: token.LBRACE, Literal: "{", Line: 14, Column: 30},
		{Type: token.STATIC, Literal: "static", Line: 15, Column: 2},
		{Type: token.LONG, Literal: "long", Line: 15, Column: 9},
		{Type: token.IDENTIFIER, Literal: "u", Line: 15, Column: 14},
		{Type: token.ASSIGN, Literal: "=", Line: 15, Column: 15},
		{Type: token.NUMERICL, Literal: "3", Line: 15, Column: 16},
		{Type: token.SEMICOLON, Literal: ";", Line: 15, Column: 17},
		{Type: token.RETURN, Literal: "return", Line: 16, Column: 2},
		{Type: token.IDENTIFIER, Literal: "e", Line: 16, Column: 9},
		{Type: token.PLUS, Literal: "+", Line: 16, Column: 11},
		{Type: token.IDENTIFIER, Literal: "u", Line: 16, Column: 13},
		{Type: token.ASTERISK, Literal: "*", Line: 16, Column: 15},
		{Type: token.IDENTIFIER, Literal: "x", Line: 16, Column: 17},
		{Type: token.SLASH, Literal: "/", Line: 16, Column: 19},
		{Type: token.IDENTIFIER, Literal: "y", Line: 16, Column: 21},
		{Type: token.SEMICOLON, Literal: ";", Line: 16, Column: 22},
		{Type: token.RBRACE, Literal: "}", Line: 17, Column: 1},
		{Type: token.ENUM, Literal: "enum", Line: 19, Column: 1},
		{Type: token.LBRACE, Literal: "{", Line: 19, Column: 6},
		{Type: token.IDENTIFIER, Literal: "A", Line: 20, Column: 2},
		{Type: token.ASSIGN, Literal: "=", Line: 20, Column: 4},
		{Type: token.NUMERICL, Literal: "0", Line: 20, Column: 6},
		{Type: token.COMMA, Literal: ",", Line: 20, Column: 7},
		{Type: token.IDENTIFIER, Literal: "B", Line: 21, Column: 2},
		{Type: token.RBRACE, Literal: "}", Line: 22, Column: 1},
		{Type: token.SEMICOLON, Literal: ";", Line: 22, Column: 2},
		{Type: token.INT, Literal: "int", Line: 24, Column: 1},
		{Type: token.IDENTIFIER, Literal: "main", Line: 24, Column: 5},
		{Type: token.LPAREN, Literal: "(", Line: 24, Column: 9},
		{Type: token.VOID, Literal: "void", Line: 24, Column: 10},
		{Type: token.RPAREN, Literal: ")", Line: 24, Column: 14},
		{Type: token.LBRACE, Literal: "{", Line: 24, Column: 16},
		{Type: token.AUTO, Literal: "auto", Line: 25, Column: 2},
		{Type: token.DOUBLE, Literal: "double", Line: 25, Column: 7},
		{Type: token.IDENTIFIER, Literal: "x", Line: 25, Column: 14},
		{Type: token.ASSIGN, Literal: "=", Line: 25, Column: 16},
		{Type: token.NUMERICL, Literal: "40.23e+20", Line: 25, Column: 18},
		{Type: token.SEMICOLON, Literal: ";", Line: 25, Column: 27},
		{Type: token.REGISTER, Literal: "register", Line: 26, Column: 2},
		{Type: token.FLOAT, Literal: "float", Line: 26, Column: 11},
		{Type: token.IDENTIFIER, Literal: "y", Line: 26, Column: 17},
		{Type: token.ASSIGN, Literal: "=", Line: 26, Column: 19},
		{Type: token.NUMERICL, Literal: "40e+20", Line: 26, Column: 21},
		{Type: token.SEMICOLON, Literal: ";", Line: 26, Column: 27},
		{Type: token.VOLATILE, Literal: "volatile", Line: 27, Column: 2},
		{Type: token.INT, Literal: "int", Line: 27, Column: 11},
		{Type: token.IDENTIFIER, Literal: "z", Line: 27, Column: 15},
		{Type: token.ASSIGN, Literal: "=", Line: 27, Column: 17},
		{Type: token.NUMERICL, Literal: "0x3000", Line: 27, Column: 19},
		{Type: token.SEMICOLON, Literal: ";", Line: 27, Column: 25},
		{Type: token.VOLATILE, Literal: "volatile", Line: 28, Column: 2},
		{Type: token.INT, Literal: "int", Line: 28, Column: 11},
		{Type: token.IDENTIFIER, Literal: "zz", Line: 28, Column: 15},
		{Type: token.ASSIGN, Literal: "=", Line: 28, Column: 18},
		{Type: token.NUMERICL, Literal: "03000", Line: 28, Column: 20},
		{Type: token.SEMICOLON, Literal: ";", Line: 28, Column: 25},
		{Type: token.STATIC, Literal: "static", Line: 29, Column: 2},
		{Type: token.CHAR, Literal: "char", Line: 29, Column: 9},
		{Type: token.IDENTIFIER, Literal: "f6", Line: 29, Column: 14},
		{Type: token.ASSIGN, Literal: "=", Line: 29, Column: 17},
		{Type: token.CHARL, Literal: "'c'", Line: 29, Column: 19},
		{Type: token.SEMICOLON, Literal: ";", Line: 29, Column: 22},
		{Type: token.IDENTIFIER, Literal: "z", Line: 30, Column: 2},
		{Type: token.BITXORA, Literal: "^=", Line: 30, Column: 4},
		{Type: token.NUMERICL, Literal: "1", Line: 30, Column: 7},
		{Type: token.SEMICOLON, Literal: ";", Line: 30, Column: 8},
		{Type: token.CHAR, Literal: "char", Line: 31, Column: 2},
		{Type: token.IDENTIFIER, Literal: "nine", Line: 31, Column: 7},
		{Type: token.LSQUARE, Literal: "[", Line: 31, Column: 11},
		{Type: token.RSQUARE, Literal: "]", Line: 31, Column: 12},
		{Type: token.ASSIGN, Literal: "=", Line: 31, Column: 14},
		{Type: token.STRINGL, Literal: "\"nine\"", Line: 31, Column: 16},
		{Type: token.SEMICOLON, Literal: ";", Line: 31, Column: 22},
		{Type: token.CHAR, Literal: "char", Line: 32, Column: 2},
		{Type: token.ASTERISK, Literal: "*", Line: 32, Column: 7},
		{Type: token.IDENTIFIER, Literal: "nine2", Line: 32, Column: 8},
		{Type: token.ASSIGN, Literal: "=", Line: 32, Column: 14},
		{Type: token.IDENTIFIER, Literal: "nine", Line: 32, Column: 16},
		{Type: token.SEMICOLON, Literal: ";", Line: 32, Column: 20},
		{Type: token.IDENTIFIER, Literal: "nine", Line: 33, Column: 2},
		{Type: token.LSQUARE, Literal: "[", Line: 33, Column: 6},
		{Type: token.NUMERICL, Literal: "2", Line: 33, Column: 7},
		{Type: token.RSQUARE, Literal: "]", Line: 33, Column: 8},
		{Type: token.ASSIGN, Literal: "=", Line: 33, Column: 10},
		{Type: token.IDENTIFIER, Literal: "f6", Line: 33, Column: 12},
		{Type: token.SEMICOLON, Literal: ";", Line: 33, Column: 14},
		{Type: token.LPAREN, Literal: "(", Line: 34, Column: 2},
		{Type: token.ASTERISK, Literal: "*", Line: 34, Column: 3},
		{Type: token.IDENTIFIER, Literal: "nine2", Line: 34, Column: 4},
		{Type: token.RPAREN, Literal: ")", Line: 34, Column: 9},
		{Type: token.INC, Literal: "++", Line: 34, Column: 10},
		{Type: token.SEMICOLON, Literal: ";", Line: 34, Column: 12},
		{Type: token.ASTERISK, Literal: "*", Line: 35, Column: 2},
		{Type: token.IDENTIFIER, Literal: "nine2", Line: 35, Column: 3},
		{Type: token.ASSIGN, Literal: "=", Line: 35, Column: 9},
		{Type: token.NUMERICL, Literal: "3", Line: 35, Column: 11},
		{Type: token.SEMICOLON, Literal: ";", Line: 35, Column: 12},
		{Type: token.IDENTIFIER, Literal: "x", Line: 36, Column: 2},
		{Type: token.PLUSA, Literal: "+=", Line: 36, Column: 4},
		{Type: token.NUMERICL, Literal: "2", Line: 36, Column: 7},
		{Type: token.RSHIFT, Literal: ">>", Line: 36, Column: 9},
		{Type: token.NUMERICL, Literal: "2", Line: 36, Column: 12},
		{Type: token.AMP, Literal: "&", Line: 36, Column: 13},
		{Type: token.NUMERICL, Literal: "3", Line: 36, Column: 14},
		{Type: token.BITXOR, Literal: "^", Line: 36, Column: 15},
		{Type: token.NUMERICL, Literal: "5", Line: 36, Column: 16},
		{Type: token.LSHIFT, Literal: "<<", Line: 36, Column: 17},
		{Type: token.NUMERICL, Literal: "2", Line: 36, Column: 19},
		{Type: token.BITOR, Literal: "|", Line: 36, Column: 20},
		{Type: token.NUMERICL, Literal: "9", Line: 36, Column: 21},
		{Type: token.SEMICOLON, Literal: ";", Line: 36, Column: 22},
		{Type: token.DO, Literal: "do", Line: 37, Column: 2},
		{Type: token.LBRACE, Literal: "{", Line: 37, Column: 5},
		{Type: token.IDENTIFIER, Literal: "z", Line: 38, Column: 3},
		{Type: token.BITANDA, Literal: "&=", Line: 38, Column: 5},
		{Type: token.NUMERICL, Literal: "4", Line: 38, Column: 8},
		{Type: token.ASTERISK, Literal: "*", Line: 38, Column: 10},
		{Type: token.NUMERICL, Literal: "99", Line: 38, Column: 12},
		{Type: token.SLASH, Literal: "/", Line: 38, Column: 15},
		{Type: token.NUMERICL, Literal: "6", Line: 38, Column: 17},
		{Type: token.PLUS, Literal: "+", Line: 38, Column: 19},
		{Type: token.NUMERICL, Literal: "7", Line: 38, Column: 21},
		{Type: token.MINUS, Literal: "-", Line: 38, Column: 23},
		{Type: token.NUMERICL, Literal: "0", Line: 38, Column: 25},
		{Type: token.SEMICOLON, Literal: ";", Line: 38, Column: 26},
		{Type: token.IDENTIFIER, Literal: "z", Line: 39, Column: 3},
		{Type: token.ASSIGN, Literal: "=", Line: 39, Column: 5},
		{Type: token.IDENTIFIER, Literal: "f6", Line: 39, Column: 7},
		{Type: token.GT, Literal: ">", Line: 39, Column: 10},
		{Type: token.NUMERICL, Literal: "3", Line: 39, Column: 12},
		{Type: token.QUESTION, Literal: "?", Line: 39, Column: 14},
		{Type: token.IDENTIFIER, Literal: "z", Line: 39, Column: 16},
		{Type: token.COLON, Literal: ":", Line: 39, Column: 18},
		{Type: token.NUMERICL, Literal: "0", Line: 39, Column: 20},
		{Type: token.SEMICOLON, Literal: ";", Line: 39, Column: 21},
		{Type: token.RBRACE, Literal: "}", Line: 40, Column: 2},
		{Type: token.WHILE, Literal: "while", Line: 40, Column: 4},
		{Type: token.LPAREN, Literal: "(", Line: 40, Column: 10},
		{Type: token.NUMERICL, Literal: "0", Line: 40, Column: 11},
		{Type: token.RPAREN, Literal: ")", Line: 40, Column: 12},
		{Type: token.SEMICOLON, Literal: ";", Line: 40, Column: 13},
		{Type: token.IF, Literal: "if", Line: 42, Column: 2},
		{Type: token.LPAREN, Literal: "(", Line: 42, Column: 5},
		{Type: token.IDENTIFIER, Literal: "z", Line: 42, Column: 6},
		{Type: token.LT, Literal: "<", Line: 42, Column: 8},
		{Type: token.NUMERICL, Literal: "1", Line: 42, Column: 10},
		{Type: token.RPAREN, Literal: ")", Line: 42, Column: 11},
		{Type: token.LBRACE, Literal: "{", Line: 42, Column: 13},
		{Type: token.IDENTIFIER, Literal: "z", Line: 43, Column: 3},
		{Type: token.SLASHA, Literal: "/=", Line: 43, Column: 5},
		{Type: token.NUMERICL, Literal: "2", Line: 43, Column: 8},
		{Type: token.SEMICOLON, Literal: ";", Line: 43, Column: 9},
		{Type: token.RBRACE, Literal: "}", Line: 44, Column: 2},
		{Type: token.ELSE, Literal: "else", Line: 44, Column: 4},
		{Type: token.IF, Literal: "if", Line: 44, Column: 9},
		{Type: token.LPAREN, Literal: "(", Line: 44, Column: 12},
		{Type: token.IDENTIFIER, Literal: "z", Line: 44, Column: 13},
		{Type: token.LTE, Literal: "<=", Line: 44, Column: 15},
		{Type: token.NUMERICL, Literal: "2", Line: 44, Column: 18},
		{Type: token.RPAREN, Literal: ")", Line: 44, Column: 19},
		{Type: token.LBRACE, Literal: "{", Line: 44, Column: 21},
		{Type: token.IDENTIFIER, Literal: "z", Line: 45, Column: 3},
		{Type: token.MODA, Literal: "%=", Line: 45, Column: 5},
		{Type: token.NUMERICL, Literal: "2", Line: 45, Column: 8},
		{Type: token.SEMICOLON, Literal: ";", Line: 45, Column: 9},
		{Type: token.RBRACE, Literal: "}", Line: 46, Column: 2},
		{Type: token.ELSE, Literal: "else", Line: 46, Column: 4},
		{Type: token.IF, Literal: "if", Line: 46, Column: 9},
		{Type: token.LPAREN, Literal: "(", Line: 46, Column: 12},
		{Type: token.IDENTIFIER, Literal: "z", Line: 46, Column: 13},
		{Type: token.GT, Literal: ">", Line: 46, Column: 15},
		{Type: token.NUMERICL, Literal: "3", Line: 46, Column: 17},
		{Type: token.RPAREN, Literal: ")", Line: 46, Column: 18},
		{Type: token.LBRACE, Literal: "{", Line: 46, Column: 20},
		{Type: token.IDENTIFIER, Literal: "z", Line: 47, Column: 3},
		{Type: token.MINUSA, Literal: "-=", Line: 47, Column: 5},
		{Type: token.NUMERICL, Literal: "2", Line: 47, Column: 8},
		{Type: token.SEMICOLON, Literal: ";", Line: 47, Column: 9},
		{Type: token.RBRACE, Literal: "}", Line: 48, Column: 2},
		{Type: token.ELSE, Literal: "else", Line: 48, Column: 4},
		{Type: token.IF, Literal: "if", Line: 48, Column: 9},
		{Type: token.LPAREN, Literal: "(", Line: 48, Column: 12},
		{Type: token.IDENTIFIER, Literal: "z", Line: 48, Column: 13},
		{Type: token.GTE, Literal: ">=", Line: 48, Column: 15},
		{Type: token.NUMERICL, Literal: "3", Line: 48, Column: 18},
		{Type: token.RPAREN, Literal: ")", Line: 48, Column: 19},
		{Type: token.LBRACE, Literal: "{", Line: 48, Column: 21},
		{Type: token.IDENTIFIER, Literal: "z", Line: 49, Column: 3},
		{Type: token.ASTERISKA, Literal: "*=", Line: 49, Column: 5},
		{Type: token.NUMERICL, Literal: "2", Line: 49, Column: 8},
		{Type: token.SEMICOLON, Literal: ";", Line: 49, Column: 9},
		{Type: token.RBRACE, Literal: "}", Line: 50, Column: 2},
		{Type: token.ELSE, Literal: "else", Line: 50, Column: 4},
		{Type: token.IF, Literal: "if", Line: 50, Column: 9},
		{Type: token.LPAREN, Literal: "(", Line: 50, Column: 12},
		{Type: token.IDENTIFIER, Literal: "z", Line: 50, Column: 13},
		{Type: token.AND, Literal: "&&", Line: 50, Column: 15},
		{Type: token.DEC, Literal: "--", Line: 50, Column: 18},
		{Type: token.IDENTIFIER, Literal: "d", Line: 50, Column: 20},
		{Type: token.RPAREN, Literal: ")", Line: 50, Column: 21},
		{Type: token.LBRACE, Literal: "{", Line: 50, Column: 23},
		{Type: token.IDENTIFIER, Literal: "z", Line: 51, Column: 3},
		{Type: token.LSHIFTA, Literal: "<<=", Line: 51, Column: 5},
		{Type: token.NUMERICL, Literal: "1", Line: 51, Column: 9},
		{Type: token.SEMICOLON, Literal: ";", Line: 51, Column: 10},
		{Type: token.IDENTIFIER, Literal: "z", Line: 52, Column: 3},
		{Type: token.RSHIFTA, Literal: ">>=", Line: 52, Column: 5},
		{Type: token.NUMERICL, Literal: "1", Line: 52, Column: 9},
		{Type: token.SEMICOLON, Literal: ";", Line: 52, Column: 10},
		{Type: token.STATIC, Literal: "static", Line: 53, Column: 3},
		{Type: token.STRUCT, Literal: "struct", Line: 53, Column: 10},
		{Type: token.IDENTIFIER, Literal: "x", Line: 53, Column: 17},
		{Type: token.IDENTIFIER, Literal: "l", Line: 53, Column: 19},
		{Type: token.SEMICOLON, Literal: ";", Line: 53, Column: 20},
		{Type: token.IDENTIFIER, Literal: "z", Line: 54, Column: 3},
		{Type: token.ASSIGN, Literal: "=", Line: 54, Column: 5},
		{Type: token.IDENTIFIER, Literal: "l", Line: 54, Column: 7},
		{Type: token.DOT, Literal: ".", Line: 54, Column: 8},
		{Type: token.IDENTIFIER, Literal: "y", Line: 54, Column: 9},
		{Type: token.DOT, Literal: ".", Line: 54, Column: 10},
		{Type: token.IDENTIFIER, Literal: "g", Line: 54, Column: 11},
		{Type: token.INC, Literal: "++", Line: 54, Column: 12},
		{Type: token.SEMICOLON, Literal: ";", Line: 54, Column: 14},
		{Type: token.STRUCT, Literal: "struct", Line: 55, Column: 3},
		{Type: token.IDENTIFIER, Literal: "x", Line: 55, Column: 10},
		{Type: token.ASTERISK, Literal: "*", Line: 55, Column: 12},
		{Type: token.IDENTIFIER, Literal: "m", Line: 55, Column: 13},
		{Type: token.SEMICOLON, Literal: ";", Line: 55, Column: 14},
		{Type: token.IDENTIFIER, Literal: "z", Line: 56, Column: 3},
		{Type: token.ASSIGN, Literal: "=", Line: 56, Column: 5},
		{Type: token.DEC, Literal: "--", Line: 56, Column: 7},
		{Type: token.IDENTIFIER, Literal: "m", Line: 56, Column: 9},
		{Type: token.ARROW, Literal: "->", Line: 56, Column: 10},
		{Type: token.IDENTIFIER, Literal: "y", Line: 56, Column: 12},
		{Type: token.DOT, Literal: ".", Line: 56, Column: 13},
		{Type: token.IDENTIFIER, Literal: "h", Line: 56, Column: 14},
		{Type: token.SEMICOLON, Literal: ";", Line: 56, Column: 15},
		{Type: token.RBRACE, Literal: "}", Line: 57, Column: 2},
		{Type: token.ELSE, Literal: "else", Line: 57, Column: 4},
		{Type: token.IF, Literal: "if", Line: 57, Column: 9},
		{Type: token.LPAREN, Literal: "(", Line: 57, Column: 12},
		{Type: token.IDENTIFIER, Literal: "z", Line: 57, Column: 13},
		{Type: token.OR, Literal: "||", Line: 57, Column: 15},
		{Type: token.IDENTIFIER, Literal: "d", Line: 57, Column: 18},
		{Type: token.RPAREN, Literal: ")", Line: 57, Column: 19},
		{Type: token.LBRACE, Literal: "{", Line: 57, Column: 21},
		{Type: token.RBRACE, Literal: "}", Line: 57, Column: 22},
		{Type: token.ELSE, Literal: "else", Line: 58, Column: 2},
		{Type: token.IF, Literal: "if", Line: 58, Column: 7},
		{Type: token.LPAREN, Literal: "(", Line: 58, Column: 10},
		{Type: token.IDENTIFIER, Literal: "z", Line: 58, Column: 11},
		{Type: token.NOTEQUALS, Literal: "!=", Line: 58, Column: 13},
		{Type: token.IDENTIFIER, Literal: "d", Line: 58, Column: 16},
		{Type: token.RPAREN, Literal: ")", Line: 58, Column: 17},
		{Type: token.LBRACE, Literal: "{", Line: 58, Column: 19},
		{Type: token.RBRACE, Literal: "}", Line: 58, Column: 20},
		{Type: token.ELSE, Literal: "else", Line: 59, Column: 2},
		{Type: token.IF, Literal: "if", Line: 59, Column: 7},
		{Type: token.LPAREN, Literal: "(", Line: 59, Column: 10},
		{Type: token.IDENTIFIER, Literal: "z", Line: 59, Column: 11},
		{Type: token.EQUALS, Literal: "==", Line: 59, Column: 13},
		{Type: token.IDENTIFIER, Literal: "d", Line: 59, Column: 16},
		{Type: token.RPAREN, Literal: ")", Line: 59, Column: 17},
		{Type: token.LBRACE, Literal: "{", Line: 59, Column: 19},
		{Type: token.RBRACE, Literal: "}", Line: 59, Column: 20},
		{Type: token.FOR, Literal: "for", Line: 61, Column: 2},
		{Type: token.LPAREN, Literal: "(", Line: 61, Column: 6},
		{Type: token.SEMICOLON, Literal: ";", Line: 61, Column: 7},
		{Type: token.SEMICOLON, Literal: ";", Line: 61, Column: 9},
		{Type: token.RPAREN, Literal: ")", Line: 61, Column: 10},
		{Type: token.LBRACE, Literal: "{", Line: 61, Column: 12},
		{Type: token.IF, Literal: "if", Line: 61, Column: 14},
		{Type: token.LPAREN, Literal: "(", Line: 61, Column: 17},
		{Type: token.NOT, Literal: "!", Line: 61, Column: 18},
		{Type: token.IDENTIFIER, Literal: "x", Line: 61, Column: 19},
		{Type: token.RPAREN, Literal: ")", Line: 61, Column: 20},
		{Type: token.CONTINUE, Literal: "continue", Line: 61, Column: 22},
		{Type: token.SEMICOLON, Literal: ";", Line: 61, Column: 30},
		{Type: token.ELSE, Literal: "else", Line: 61, Column: 32},
		{Type: token.BREAK, Literal: "break", Line: 61, Column: 37},
		{Type: token.SEMICOLON, Literal: ";", Line: 61, Column: 42},
		{Type: token.RBRACE, Literal: "}", Line: 61, Column: 44},
		{Type: token.SWITCH, Literal: "switch", Line: 63, Column: 2},
		{Type: token.LPAREN, Literal: "(", Line: 63, Column: 9},
		{Type: token.IDENTIFIER, Literal: "f6", Line: 63, Column: 10},
		{Type: token.RPAREN, Literal: ")", Line: 63, Column: 12},
		{Type: token.LBRACE, Literal: "{", Line: 63, Column: 14},
		{Type: token.CASE, Literal: "case", Line: 64, Column: 3},
		{Type: token.CHARL, Literal: "'\\n'", Line: 64, Column: 8},
		{Type: token.COLON, Literal: ":", Line: 64, Column: 12},
		{Type: token.GOTO, Literal: "goto", Line: 65, Column: 4},
		{Type: token.IDENTIFIER, Literal: "exit", Line: 65, Column: 9},
		{Type: token.SEMICOLON, Literal: ";", Line: 65, Column: 13},
		{Type: token.DEFAULT, Literal: "default", Line: 66, Column: 3},
		{Type: token.COLON, Literal: ":", Line: 66, Column: 10},
		{Type: token.BREAK, Literal: "break", Line: 67, Column: 4},
		{Type: token.SEMICOLON, Literal: ";", Line: 67, Column: 9},
		{Type: token.RBRACE, Literal: "}", Line: 68, Column: 2},
		{Type: token.IDENTIFIER, Literal: "exit", Line: 70, Column: 1},
		{Type: token.COLON, Literal: ":", Line: 70, Column: 5},
		{Type: token.RETURN, Literal: "return", Line: 72, Column: 2},
		{Type: token.IDENTIFIER, Literal: "fnc", Line: 72, Column: 9},
		{Type: token.LPAREN, Literal: "(", Line: 72, Column: 12},
		{Type: token.IDENTIFIER, Literal: "z", Line: 72, Column: 13},
		{Type: token.PLUS, Literal: "+", Line: 72, Column: 15},
		{Type: token.NUMERICL, Literal: "3", Line: 72, Column: 17},
		{Type: token.COMMA, Literal: ",", Line: 72, Column: 18},
		{Type: token.IDENTIFIER, Literal: "z", Line: 72, Column: 20},
		{Type: token.MINUS, Literal: "-", Line: 72, Column: 22},
		{Type: token.NUMERICL, Literal: "2", Line: 72, Column: 24},
		{Type: token.RPAREN, Literal: ")", Line: 72, Column: 25},
		{Type: token.SEMICOLON, Literal: ";", Line: 72, Column: 26},
		{Type: token.RBRACE, Literal: "}", Line: 73, Column: 1},
		{Type: token.EOF, Literal: "", Line: 73, Column: 19},
	}

	l := New(input)
	for i, expected := range expectedTokens {
		tok := l.NextToken()
		if tok.Type != expected.Type {
			t.Fatalf("[test %d] tok.Type != %s, got=%s", i, expected.Type,
				tok.Type)
		}

		if tok.Literal != expected.Literal {
			t.Fatalf("[%d] tok.Literal != %s, got=%s", i, expected.Literal,
				tok.Literal)
		}

		if tok.Line != expected.Line {
			t.Fatalf("[%d] tok.Line != %d, got=%d", i, expected.Line,
				tok.Line)
		}

		if tok.Column != expected.Column {
			t.Fatalf("[%d] tok.Column != %d, got=%d", i, expected.Column,
				tok.Column)
		}
	}
}
