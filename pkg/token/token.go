package token

type TokenType string
type Token struct {
	Type      TokenType
	Literal   string
	HasAssign bool
	Line      int
	Column    int
}

// TokenTypes
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
	IDENTIFIER = "IDENTIFIER"
	INTL       = "INTL"
	FLOATL     = "FLOATL"
	CHARL      = "CHARL"
	STRINGL    = "STRINGL"

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
	QUESTION  = "?"

	AMP     = "&"
	BITANDA = "&="
	BITOR   = "|"
	BITORA  = "|="
	BITXOR  = "^"
	BITXORA = "^="
	BITNOT  = "~"

	INC = "++"
	DEC = "--"

	DOT   = "."
	ARROW = "->"

	LPAREN    = "("
	RPAREN    = ")"
	LSQUARE   = "["
	RSQUARE   = "]"
	LBRACE    = "{"
	RBRACE    = "}"
	COMMA     = ","
	SEMICOLON = ";"
	COLON     = ":"

	EOF     = "EOF"
	ILLEGAL = "ILLEGAL"
)

var keywords = map[string]TokenType{
	"auto":     AUTO,
	"break":    BREAK,
	"case":     CASE,
	"char":     CHAR,
	"const":    CONST,
	"continue": CONTINUE,
	"default":  DEFAULT,
	"do":       DO,
	"double":   DOUBLE,
	"else":     ELSE,
	"enum":     ENUM,
	"extern":   EXTERN,
	"float":    FLOAT,
	"for":      FOR,
	"goto":     GOTO,
	"if":       IF,
	"int":      INT,
	"long":     LONG,
	"register": REGISTER,
	"return":   RETURN,
	"short":    SHORT,
	"signed":   SIGNED,
	"sizeof":   SIZEOF,
	"static":   STATIC,
	"struct":   STRUCT,
	"switch":   SWITCH,
	"typedef":  TYPEDEF,
	"union":    UNION,
	"unsigned": UNSIGNED,
	"void":     VOID,
	"volatile": VOLATILE,
	"while":    WHILE,
}

func LookupIdent(ident string) TokenType {
	keyword, ok := keywords[ident]
	if !ok {
		return IDENTIFIER
	}

	return keyword
}
