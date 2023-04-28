package parser

import (
	"github.com/tjarjoura/cc/pkg/ast"
	"github.com/tjarjoura/cc/pkg/lexer"
	"github.com/tjarjoura/cc/pkg/token"
)

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

	p.registerParseFns()

	p.nextToken()
	p.nextToken()
	return &p
}

func (p *Parser) Errors() []ParseError {
	return p.errors
}

func (p *Parser) nextToken() {
	p.currToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) currTokenIs(t token.TokenType) bool {
	return p.currToken.Type == t
}

func (p *Parser) peekTokenIs(ts ...token.TokenType) bool {
	for _, t := range ts {
		if p.peekToken.Type == t {
			return true
		}
	}
	return false
}

func (p *Parser) expectPeek(ts ...token.TokenType) bool {
	if p.peekTokenIs(ts...) {
		p.nextToken()
		return true
	}

	p.peekError(ts...)
	return false
}

func (p *Parser) Parse() *ast.TranslationUnit {
	tUnit := &ast.TranslationUnit{DeclarationStatements: []*ast.DeclarationStatement{}}
	for !p.currTokenIs(token.EOF) {
		for p.currTokenIs(token.SEMICOLON) { // skip blank statements
			p.nextToken()
		}

		decl := p.parseDeclarationStatement(true)
		tUnit.DeclarationStatements = append(tUnit.DeclarationStatements, decl)
		p.nextToken()
	}
	return tUnit
}
