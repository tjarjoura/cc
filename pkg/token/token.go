package token

type TokenType string
type Token struct {
	Type    TokenType
	Literal string
}

const (
	// keywords
	AUTO     = "AUTO"
	BREAK    = "BREAK"
	CASE     = "CASE"
	CHAR     = "CHAR"
	CONST    = "CONST"
	CONTINUE = "CONTINUE"
	DEFAULT  = "DEFAULT"
	DO       = "DO"
	DOUBLE   = "DOUBLE"
	ELSE     = "ELSE"
	ENUM     = "ENUM"
	EXTERN   = "EXTERN"
	FLOAT    = "FLOAT"
	FOR      = "FOR"
	GOTO     = "GOTO"
	IF       = "IF"
	INT      = "INT"
	LONG     = "LONG"
	REGISTER = "REGISTER"
	RETURN   = "RETURN"
	SHORT    = "SHORT"
	SIGNED   = "SIGNED"
	SIZEOF   = "SIZEOF"
	STATIC   = "STATIC"
	STRUCT   = "STRUCT"
	SWITCH   = "SWITCH"
	TYPEDEF  = "TYPEDEF"
	UNION    = "UNION"
	UNSIGNED = "UNSIGNED"
	VOID     = "VOID"
	VOLATILE = "VOLATILE"
	WHILE    = "WHILE"

	// identifiers and literals
	IDENT   = "IDENT"
	INTL    = "INTL"
	CHARL   = "CHARL"
	STRINGL = "STRINGL"
	REALL   = "REALL"

	// Operators and separators
	ASSIGN    = "="
	PLUS      = "+"
	PLUSA     = "+="
	MINUS     = "-"
	MINUSA    = "-="
	ASTERISK  = "*"
	ASTERISKA = "*="
	SLASH     = "/"
	SLASHA    = "/="
	MOD       = "%"
	MODA      = "%="
	LSHIFT    = "<<"
	LSHIFTA   = "<<="
	RSHIFT    = ">>"
	RSHIFTA   = ">>="

	EQUALS    = "=="
	NOTEQUALS = "!="
	GT        = ">"
	GTE       = ">="
	LT        = "<"
	LTE       = "<="
	AND       = "&&"
	OR        = "||"
	NOT       = "!"

	BITAND  = "&"
	BITANDA = "&="
	BITOR   = "|"
	BITORA  = "|="
	BITXOR  = "^"
	BITXORA = "^="

	INC = "++"
	DEC = "--"

	DOT   = "."
	ARROW = "->"

	EOF     = "EOF"
	ILLEGAL = "ILLEGAL"
)
