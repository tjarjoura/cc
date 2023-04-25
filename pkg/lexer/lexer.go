package lexer

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/tjarjoura/cc/pkg/token"
)

type Lexer struct {
	input  string
	pos    int
	peek   int
	peek2  int
	char   byte
	line   int
	column int
}

func (l *Lexer) peek2Char() byte {
	if l.peek2 > len(l.input) {
		return 0
	}
	return l.input[l.peek2]
}

func (l *Lexer) peekChar() byte {
	if l.peek > len(l.input) {
		return 0
	}
	return l.input[l.peek]
}

func (l *Lexer) readChar() {
	isNewline := (l.char == '\n')
	if l.peek >= len(l.input) {
		l.char = 0
	} else {
		l.char = l.input[l.peek]
	}

	if isNewline {
		l.line++
		l.column = 1
	} else {
		l.column++
	}

	l.pos = l.peek
	l.peek = l.peek2
	l.peek2 += 1
}

func isLetter(char byte) bool {
	return (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z')
}

func isDigit(char byte) bool {
	return char >= '0' && char <= '9'
}

func (l *Lexer) skipLineComment() {
	for l.char != '\n' && l.char != 0 {
		l.readChar()
	}
}

func (l *Lexer) skipBlockComment() {
	for !(l.char == '*' && l.peekChar() == '/') && l.char != 0 {
		l.readChar()
	}

	l.readChar()
	l.readChar()
}

func (l *Lexer) skipWhitespaceAndComments() {
	for {
		if l.char == '/' {
			if l.peekChar() == '/' {
				l.skipLineComment()
			} else if l.peekChar() == '*' {
				l.skipBlockComment()
			} else {
				break
			}
		} else if l.char == ' ' || l.char == '\t' || l.char == '\n' {
			l.readChar()
		} else {
			break
		}
	}
}

func (l *Lexer) checkMultiCharOp(next byte, next2 byte,
	tt token.TokenType) (token.Token, bool) {
	prev := l.char
	if next == '|' {
		fmt.Printf("checkMultiCharOp. l.char=%c l.peek=%c l.peek2=%c\n",
			l.char, l.peekChar(), l.peek2Char())
	}
	if next2 == 0 {
		if l.peekChar() == next {
			l.readChar()
			literal := fmt.Sprintf("%c%c", prev, l.char)
			return token.Token{Type: tt, Literal: literal}, true
		}
	} else {
		if l.peekChar() == next && l.peek2Char() == next2 {
			prev := l.char
			l.readChar()
			prev2 := l.char
			l.readChar()
			literal := fmt.Sprintf("%c%c%c", prev, prev2, l.char)
			return token.Token{Type: tt, Literal: literal}, true
		}
	}

	return token.Token{}, false
}

func (l *Lexer) readIdent() string {
	start := l.pos
	for isLetter(l.char) || isDigit(l.char) || l.char == '_' {
		l.readChar()
	}

	return l.input[start:l.pos]
}

func (l *Lexer) readNumber() string {
	re := regexp.MustCompile(`(^(0x)?(\d|[a-f]|[A-F])+(U|L|UL)?)|(^((\d*\.\d+)|(\d+\.*\d*))(e(\+|-)?\d+)?)`)
	re.Longest()
	indices := re.FindStringIndex(l.input[l.pos:])
	if indices == nil {
		return ""
	}

	start := l.pos
	for i := 0; i < indices[1]; i++ {
		l.readChar()
	}

	return l.input[start:l.pos]
}

func (l *Lexer) readCharLiteral() string {
	if l.char != '\'' {
		return ""
	}

	start := l.pos
	l.readChar()

	for l.char != '\'' && l.char != 0 {
		l.readChar()
	}

	return l.input[start:l.peek]
}

func (l *Lexer) readString() string {
	if l.char != '"' {
		return ""
	}

	start := l.pos
	l.readChar()

	for l.char != '"' && l.char != 0 {
		l.readChar()
	}

	return l.input[start:l.peek]
}

func New(input string) *Lexer {
	l := &Lexer{input: input, line: 1, column: 1}
	l.readChar()
	l.readChar()
	return l
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token
	var ok bool

	l.skipWhitespaceAndComments()

	line, column := l.line, l.column
	switch l.char {
	case '"':
		literal := l.readString()
		tok = token.Token{Type: token.STRINGL, Literal: literal}
	case '\'':
		literal := l.readCharLiteral()
		tok = token.Token{Type: token.CHARL, Literal: literal}
	case '=':
		tok, ok = l.checkMultiCharOp('=', 0, token.EQUALS)
		if !ok {
			tok = token.Token{Type: token.ASSIGN, Literal: string(l.char)}
		}
	case '+':
		tok, ok = l.checkMultiCharOp('=', 0, token.PLUSA)
		if !ok {
			tok, ok = l.checkMultiCharOp('+', 0, token.INC)
			if !ok {
				tok = token.Token{Type: token.PLUS,
					Literal: string(l.char)}
			}
		}
	case '-':
		tok, ok = l.checkMultiCharOp('=', 0, token.MINUSA)
		if !ok {
			tok, ok = l.checkMultiCharOp('-', 0, token.DEC)
			if !ok {
				tok, ok = l.checkMultiCharOp('>', 0, token.ARROW)
				if !ok {
					tok = token.Token{Type: token.MINUS,
						Literal: string(l.char)}
				}
			}
		}
	case '*':
		tok, ok = l.checkMultiCharOp('=', 0, token.ASTERISKA)
		if !ok {
			tok = token.Token{Type: token.ASTERISK, Literal: string(l.char)}
		}
	case '/':
		tok, ok = l.checkMultiCharOp('=', 0, token.SLASHA)
		if !ok {
			tok = token.Token{Type: token.SLASH, Literal: string(l.char)}
		}
	case '%':
		tok, ok = l.checkMultiCharOp('=', 0, token.MODA)
		if !ok {
			tok = token.Token{Type: token.MOD, Literal: string(l.char)}
		}
	case '!':
		tok, ok = l.checkMultiCharOp('=', 0, token.NOTEQUALS)
		if !ok {
			tok = token.Token{Type: token.NOT, Literal: string(l.char)}
		}
	case '<':
		tok, ok = l.checkMultiCharOp('<', '=', token.LSHIFTA)
		if !ok {
			tok, ok = l.checkMultiCharOp('<', 0, token.LSHIFT)
			if !ok {
				tok, ok = l.checkMultiCharOp('=', 0, token.LTE)
				if !ok {
					tok = token.Token{Type: token.LT,
						Literal: string(l.char)}
				}
			}
		}
	case '>':
		tok, ok = l.checkMultiCharOp('>', '=', token.RSHIFTA)
		if !ok {
			tok, ok = l.checkMultiCharOp('>', 0, token.RSHIFT)
			if !ok {
				tok, ok = l.checkMultiCharOp('=', 0, token.GTE)
				if !ok {
					tok = token.Token{Type: token.GT,
						Literal: string(l.char)}
				}
			}
		}
	case '&':
		tok, ok = l.checkMultiCharOp('&', 0, token.AND)
		if !ok {
			tok, ok = l.checkMultiCharOp('=', 0, token.BITANDA)
			if !ok {
				tok = token.Token{Type: token.AMP,
					Literal: string(l.char)}
			}
		}
	case '|':
		tok, ok = l.checkMultiCharOp('|', 0, token.OR)
		if !ok {
			tok, ok = l.checkMultiCharOp('=', 0, token.BITORA)
			if !ok {
				tok = token.Token{Type: token.BITOR,
					Literal: string(l.char)}
			}
		}
	case '^':
		tok, ok = l.checkMultiCharOp('=', 0, token.BITXORA)
		if !ok {
			tok = token.Token{Type: token.BITXOR, Literal: string(l.char)}
		}
	case '?':
		tok = token.Token{Type: token.QUESTION, Literal: string(l.char)}
	case '.':
		tok = token.Token{Type: token.DOT, Literal: string(l.char)}
	case '(':
		tok = token.Token{Type: token.LPAREN, Literal: string(l.char)}
	case ')':
		tok = token.Token{Type: token.RPAREN, Literal: string(l.char)}
	case '[':
		tok = token.Token{Type: token.LSQUARE, Literal: string(l.char)}
	case ']':
		tok = token.Token{Type: token.RSQUARE, Literal: string(l.char)}
	case '{':
		tok = token.Token{Type: token.LBRACE, Literal: string(l.char)}
	case '}':
		tok = token.Token{Type: token.RBRACE, Literal: string(l.char)}
	case ',':
		tok = token.Token{Type: token.COMMA, Literal: string(l.char)}
	case ';':
		tok = token.Token{Type: token.SEMICOLON, Literal: string(l.char)}
	case ':':
		tok = token.Token{Type: token.COLON, Literal: string(l.char)}
	case 0:
		tok = token.Token{Type: token.EOF, Literal: ""}
	default:
		if isLetter(l.char) || l.char == '_' {
			ident := l.readIdent()
			tokenType := token.LookupIdent(ident)
			return token.Token{Type: tokenType, Literal: ident,
				Line: line, Column: column}
		} else if isDigit(l.char) {
			number := l.readNumber()
			var tokenType token.TokenType = token.INTL
			if strings.IndexByte(number, '.') > -1 || strings.IndexByte(number, 'e') > -1 {
				tokenType = token.FLOATL
			}
			return token.Token{Type: tokenType, Literal: number,
				Line: line, Column: column}
		} else {
			tok = token.Token{Type: token.ILLEGAL,
				Literal: string(l.char)}
		}

	}

	l.readChar()
	tok.Line, tok.Column = line, column
	return tok
}
