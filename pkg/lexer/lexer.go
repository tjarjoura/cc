package lexer

import "github.com/tjarjoura/cc/pkg/token"

type Lexer struct {
	input string
}

func New(input string) {
	return &Lexer{input: input}
}

func (l *Lexer) NextToken() token.Token {
}
