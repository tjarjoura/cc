package parser

import (
	"bytes"
	"fmt"

	"github.com/tjarjoura/cc/pkg/token"
)

type ParseError struct {
	msg   string
	token token.Token
}

func (p *ParseError) String() string {
	return fmt.Sprintf("[%d:%d] %s", p.token.Line, p.token.Column, p.msg)
}

func (p *Parser) genericError(msg string) {
	err := ParseError{msg: msg, token: p.currToken}
	p.errors = append(p.errors, err)
}

func (p *Parser) peekError(ts ...token.TokenType) {
	var expectedTokens bytes.Buffer
	for _, t := range ts[:len(ts)-1] {
		expectedTokens.WriteString(string(t))
		expectedTokens.WriteString(" or ")
	}
	expectedTokens.WriteString(string(ts[len(ts)-1]))

	//panic("ahhhhh!!!")
	msg := fmt.Sprintf("expected next token to be '%s', got '%s' instead",
		expectedTokens.String(), p.peekToken.Literal)
	p.errors = append(p.errors, ParseError{msg: msg, token: p.peekToken})
}
