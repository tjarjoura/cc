package parser

import (
	"fmt"

	"github.com/tjarjoura/cc/pkg/ast"
	"github.com/tjarjoura/cc/pkg/lexer"
	"github.com/tjarjoura/cc/pkg/token"
)

type ParseError struct {
	msg   string
	token token.Token
}

func (p *ParseError) String() string {
	return fmt.Sprintf("[%d:%d] %s", p.token.Line, p.token.Column, p.msg)
}

type Parser struct {
	l      *lexer.Lexer
	errors []ParseError

	currToken token.Token
	peekToken token.Token

	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

func New(l *lexer.Lexer) *Parser {
	p := Parser{l: l, errors: []ParseError{}}

	p.nextToken()
	p.nextToken()
	return &p
}

func (p *Parser) Errors() []ParseError {
	return p.errors
}

func (p *Parser) peekError(ts ...token.TokenType) {
	msg := fmt.Sprintf("expected next token to be one of %v, got %s instead",
		ts, p.peekToken.Type)
	p.errors = append(p.errors, ParseError{msg: msg, token: p.peekToken})
}

func (p *Parser) nextToken() {
	p.currToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) currTokenIs(t token.TokenType) bool {
	return p.currToken.Type == t
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) expectPeek(ts ...token.TokenType) bool {
	for _, t := range ts {
		if p.peekTokenIs(t) {
			p.nextToken()
			return true
		}
	}

	p.peekError(ts...)
	return false
}

func (p *Parser) Parse() *ast.TranslationUnit {
	tUnit := &ast.TranslationUnit{Declarations: []ast.Declaration{}}
	for !p.currTokenIs(token.EOF) {
		decls := p.parseDeclarations()
		tUnit.Declarations = append(tUnit.Declarations, decls...)
		p.nextToken()
	}
	return tUnit
}
