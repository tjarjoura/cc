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
	COMMA        // ,
	BITANDA      // &= |= ^=
	BSHIFTA      // <<= >>=
	PRODA        // *= /= %=
	PLUSA        // += -=
	ASSIGN       // =
	TERNARY      // ?:
	OR           // ||
	AND          // &&
	BITOR        // |
	BITXOR       // ^
	BITAND       // &
	EQUALS       // == !=
	GT           // > >=
	LT           // < <=
	BSHIFT       // >> <<
	SUM          // + -
	PRODUCT      // * / %
	SIZEOF       // sizeof
	ADDRESSOF    // &
	DEREF        // *
	CAST         // (type) value
	NOT          // ! ~
	UNARYPLUS    // +(value) -(value)
	PREINC       // ++var
	STRUCTACCPTR // s->m
	STRUCTACC    // s.m
	SUBSCRIPT    // arr[idx]
	CALL         // fn(x)
	POSTINC      // var++
)

var infixPrecedenceMap = map[token.TokenType]int{
	token.COMMA:     COMMA,
	token.BITANDA:   BITANDA,
	token.BITORA:    BITANDA,
	token.BITXORA:   BITANDA,
	token.LSHIFTA:   BSHIFTA,
	token.RSHIFTA:   BSHIFTA,
	token.ASTERISKA: PRODA,
	token.SLASHA:    PRODA,
	token.MODA:      PRODA,
	token.PLUSA:     PLUSA,
	token.MINUSA:    PLUSA,
	token.ASSIGN:    ASSIGN,
	token.QUESTION:  TERNARY,
	token.COLON:     TERNARY,
	token.OR:        OR,
	token.AND:       AND,
	token.BITOR:     BITOR,
	token.BITXOR:    BITXOR,
	token.AMP:       BITAND,
	token.EQUALS:    EQUALS,
	token.NOTEQUALS: EQUALS,
	token.GT:        GT,
	token.LT:        LT,
	token.LSHIFT:    BSHIFT,
	token.RSHIFT:    BSHIFT,
	token.PLUS:      SUM,
	token.MINUS:     SUM,
	token.ASTERISK:  PRODUCT,
	token.SLASH:     PRODUCT,
	token.MOD:       PRODUCT,
	token.ARROW:     STRUCTACCPTR,
	token.DOT:       STRUCTACC,
	token.LSQUARE:   SUBSCRIPT,
	token.LPAREN:    CALL,
	token.INC:       POSTINC,
	token.DEC:       POSTINC,
}

var prefixPrecedenceMap = map[token.TokenType]int{
	token.SIZEOF:   SIZEOF,
	token.AMP:      ADDRESSOF,
	token.ASTERISK: DEREF,
	// TODO cast
	token.NOT:    NOT,
	token.BITNOT: NOT,
	token.PLUS:   UNARYPLUS,
	token.MINUS:  UNARYPLUS,
	token.INC:    PREINC,
}

func (p *Parser) registerParseFns() {
	p.prefixParseFns = map[token.TokenType]prefixParseFn{}
	p.infixParseFns = map[token.TokenType]infixParseFn{}

	p.prefixParseFns[token.IDENTIFIER] = p.parseIdentifier
	p.prefixParseFns[token.INTL] = p.parseIntegerLiteral
	p.prefixParseFns[token.FLOATL] = p.parseFloatLiteral

	p.prefixParseFns[token.LPAREN] = p.parseGroupedExpression // TODO handle type cast too

	p.prefixParseFns[token.BITNOT] = p.parsePrefixExpression
	p.prefixParseFns[token.NOT] = p.parsePrefixExpression
	p.prefixParseFns[token.ASTERISK] = p.parsePrefixExpression
	p.prefixParseFns[token.AMP] = p.parsePrefixExpression
	p.prefixParseFns[token.MINUS] = p.parsePrefixExpression
	p.prefixParseFns[token.PLUS] = p.parsePrefixExpression

	p.infixParseFns[token.ASSIGN] = p.parseInfixExpression
	p.infixParseFns[token.PLUS] = p.parseInfixExpression
	p.infixParseFns[token.PLUSA] = p.parseInfixExpression
	p.infixParseFns[token.MINUS] = p.parseInfixExpression
	p.infixParseFns[token.MINUSA] = p.parseInfixExpression
	p.infixParseFns[token.ASTERISK] = p.parseInfixExpression
	p.infixParseFns[token.ASTERISKA] = p.parseInfixExpression
	p.infixParseFns[token.SLASH] = p.parseInfixExpression
	p.infixParseFns[token.SLASHA] = p.parseInfixExpression
	p.infixParseFns[token.MOD] = p.parseInfixExpression
	p.infixParseFns[token.MODA] = p.parseInfixExpression
	p.infixParseFns[token.LSHIFT] = p.parseInfixExpression
	p.infixParseFns[token.LSHIFTA] = p.parseInfixExpression
	p.infixParseFns[token.RSHIFT] = p.parseInfixExpression
	p.infixParseFns[token.RSHIFTA] = p.parseInfixExpression
	p.infixParseFns[token.EQUALS] = p.parseInfixExpression
	p.infixParseFns[token.NOTEQUALS] = p.parseInfixExpression
	p.infixParseFns[token.GT] = p.parseInfixExpression
	p.infixParseFns[token.GTE] = p.parseInfixExpression
	p.infixParseFns[token.LT] = p.parseInfixExpression
	p.infixParseFns[token.LTE] = p.parseInfixExpression
	p.infixParseFns[token.AND] = p.parseInfixExpression
	p.infixParseFns[token.OR] = p.parseInfixExpression
	p.infixParseFns[token.AMP] = p.parseInfixExpression
	p.infixParseFns[token.BITANDA] = p.parseInfixExpression
	p.infixParseFns[token.BITOR] = p.parseInfixExpression
	p.infixParseFns[token.BITORA] = p.parseInfixExpression
	p.infixParseFns[token.BITXOR] = p.parseInfixExpression
	p.infixParseFns[token.BITXORA] = p.parseInfixExpression
}

func (p *Parser) noPrefixParseFnError(t token.TokenType) {
	msg := fmt.Sprintf("no prefix parse function for %s found", t)
	p.errors = append(p.errors, ParseError{msg: msg, token: p.currToken})
}

func (p *Parser) currPrecedence() int {
	if p, ok := infixPrecedenceMap[p.currToken.Type]; ok {
		return p
	}

	return LOWEST
}

func (p *Parser) peekPrecedence() int {
	if p, ok := infixPrecedenceMap[p.peekToken.Type]; ok {
		return p
	}

	return LOWEST
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()
	expr := p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return expr
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	prefixExpr := &ast.PrefixExpression{
		Token:    p.currToken,
		Operator: p.currToken.Literal,
	}

	precedence, ok := prefixPrecedenceMap[p.currToken.Type]
	if !ok {
		msg := fmt.Sprintf("couldn't find prefix precedence for %s operator",
			p.currToken.Literal)
		p.errors = append(p.errors, ParseError{token: p.currToken, msg: msg})
		return nil
	}

	p.nextToken()
	prefixExpr.Right = p.parseExpression(precedence)
	return prefixExpr
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	infixExpr := &ast.InfixExpression{
		Token:    p.currToken,
		Left:     left,
		Operator: p.currToken.Literal,
	}

	precedence := p.currPrecedence()
	p.nextToken()
	infixExpr.Right = p.parseExpression(precedence)

	return infixExpr
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.currToken, Value: p.currToken.Literal}
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
