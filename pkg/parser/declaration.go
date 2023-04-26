package parser

import (
	"bytes"
	"fmt"

	"github.com/tjarjoura/cc/pkg/ast"
	"github.com/tjarjoura/cc/pkg/token"
)

func (p *Parser) currTokenIsStorageClass() bool {
	return p.currTokenIs(token.STATIC) ||
		p.currTokenIs(token.EXTERN) ||
		p.currTokenIs(token.AUTO) ||
		p.currTokenIs(token.REGISTER)
}

func (p *Parser) peekTokenIsTypeQualifier() bool {
	return p.peekTokenIs(token.CONST) || p.peekTokenIs(token.VOLATILE)
}

func (p *Parser) currTokenIsType() bool {
	return p.currTokenIs(token.INT) ||
		p.currTokenIs(token.LONG) ||
		p.currTokenIs(token.CHAR) ||
		p.currTokenIs(token.SHORT) ||
		p.currTokenIs(token.SIGNED) ||
		p.currTokenIs(token.UNSIGNED) ||
		p.currTokenIs(token.FLOAT) ||
		p.currTokenIs(token.DOUBLE) ||
		p.currTokenIs(token.VOID)
}

func (p *Parser) peekTokenIsType() bool {
	return p.peekTokenIs(token.INT) ||
		p.peekTokenIs(token.LONG) ||
		p.peekTokenIs(token.CHAR) ||
		p.peekTokenIs(token.SHORT) ||
		p.peekTokenIs(token.SIGNED) ||
		p.peekTokenIs(token.UNSIGNED) ||
		p.peekTokenIs(token.FLOAT) ||
		p.peekTokenIs(token.DOUBLE) ||
		p.peekTokenIs(token.VOID)
}

func (p *Parser) parseDeclaratorLeft(t ast.Declaration) ast.Declaration {
	switch p.currToken.Type {
	case token.ASTERISK:
		pointer := &ast.Pointer{PointsTo: t}
		for {
			p.nextToken()
			if p.currTokenIs(token.CONST) {
				pointer.Const = true
			} else if p.currTokenIs(token.VOLATILE) {
				pointer.Volatile = true
			} else {
				break
			}
		}

		return p.parseDeclaratorLeft(pointer)
	case token.LPAREN:
		p.nextToken()
		interior := p.parseDeclaratorLeft(t)
		right := p.parseDeclaratorRight(t)

		// We need to swap out the inner most type with what we
		// parsed from the right
		t := interior
		for t.Type() != nil && t.Type().Type() != nil {
			t = t.Type()
		}

		t.SetType(right)
		return interior

	case token.IDENTIFIER:
		name := p.currToken.Literal
		right := p.parseDeclaratorRight(t)
		return &ast.VariableDeclaration{TypeSpec: right,
			Name: name}
	default:
		return nil
	}
}

func (p *Parser) parseDeclaratorRight(decl ast.Declaration) ast.Declaration {
	switch p.peekToken.Type {
	case token.LSQUARE:
		p.nextToken()
		// TODO parse expression for array size
		for !p.peekTokenIs(token.RSQUARE) && !p.peekTokenIs(token.EOF) {
			p.nextToken()
		}

		// TODO figure out how to handle syntax errors
		if !p.expectPeek(token.RSQUARE) {
			return nil
		}

		right := p.parseDeclaratorRight(decl)
		return &ast.Array{ArrayOf: right}
	case token.RPAREN:
		p.nextToken()
		return decl
	default:
		return decl
	}
}

func (p *Parser) parseTypeSpecificiation() *ast.TypeSpecification {
	typeSpec := &ast.TypeSpecification{}
	var typeName bytes.Buffer
	for {
		if p.currTokenIs(token.CONST) {
			typeSpec.Const = true
		} else if p.currTokenIs(token.VOLATILE) {
			typeSpec.Volatile = true
		} else if p.currTokenIsType() {
			typeName.WriteString(fmt.Sprintf("%s ", p.currToken.Literal))
		}

		if p.peekTokenIsType() || p.peekTokenIsTypeQualifier() {
			p.nextToken()
		} else {
			break
		}
	}

	tn := typeName.String()
	if len(tn) > 0 { // remove trailing space
		typeSpec.Name = tn[:len(tn)-1]
	}

	return typeSpec
}

func (p *Parser) parseDeclarations() []ast.Declaration {
	var decls = []ast.Declaration{}
	var storageClass string
	if p.currTokenIsStorageClass() {
		storageClass = p.currToken.Literal
		p.nextToken()
	}

	typeSpec := p.parseTypeSpecificiation()

	for !p.peekTokenIs(token.SEMICOLON) && !p.peekTokenIs(token.EOF) {
		if !p.expectPeek(token.IDENTIFIER, token.LPAREN, token.ASTERISK) {
			return nil
		}

		decl := p.parseDeclaratorLeft(typeSpec)
		if p.peekTokenIs(token.ASSIGN) { // also define the variable
			p.nextToken()
			p.nextToken()
			definition := p.parseExpression(LOWEST)
			varDecl, ok := decl.(*ast.VariableDeclaration)
			if ok {
				varDecl.Definition = definition
			}
		}

		decls = append(decls, decl)

		if p.peekTokenIs(token.COMMA) {
			p.nextToken()
		}
	}

	if !p.expectPeek(token.SEMICOLON) {
		return nil
	}

	for _, d := range decls {
		varDecl, ok := d.(*ast.VariableDeclaration)
		if ok {
			varDecl.StorageClass = storageClass
		}
	}

	return decls
}
