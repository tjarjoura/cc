package parser

import (
	"fmt"
	"strconv"

	"github.com/tjarjoura/cc/pkg/ast"
	"github.com/tjarjoura/cc/pkg/token"
)

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

// Operator precedence
const (
	_ int = iota
	LOWEST
	EQUALS      // ==
	LESSGREATER // > or <
	SUM         // +
	PRODUCT     // *
	PREFIX      // -X or !X
	CALL        // myFunction(x)
)

var precedenceMap = map[token.TokenType]int{}

func (p *Parser) registerParseFns() {
	p.prefixParseFns[token.INTL] = p.parseIntegerLiteral
}

func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

func (p *Parser) noPrefixParseFnError(t token.TokenType) {
	msg := fmt.Sprintf("no prefix parse function for %s found", t)
	p.errors = append(p.errors, ParseError{msg: msg, token: p.currToken})
}

func (p *Parser) peekPrecedence() int {
	if p, ok := precedenceMap[p.peekToken.Type]; ok {
		return p
	}

	return LOWEST
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	val, err := strconv.ParseInt(p.currToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.currToken.Literal)
		p.errors = append(p.errors, ParseError{token: p.currToken, msg: msg})
		return nil
	}

	return &ast.IntegerLiteral{Token: p.currToken, Value: val}
}

func (p *Parser) parseFloatLiteral() ast.Expression {
	val, err := strconv.ParseFloat(p.currToken.Literal, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as floating point", p.currToken.Literal)
		p.errors = append(p.errors, ParseError{token: p.currToken, msg: msg})
		return nil
	}

	return &ast.FloatLiteral{Token: p.currToken, Value: val}
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.currToken.Type]
	if prefix == nil {
		p.noPrefixParseFnError(p.currToken.Type)
		return nil
	}

	leftExp := prefix()
	for p.peekPrecedence() > precedence {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}

		p.nextToken()
		leftExp = infix(leftExp)
	}

	return leftExp
}
