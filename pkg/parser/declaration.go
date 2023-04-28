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

func (p *Parser) parseDeclaratorLeft(decl ast.Declaration, insideParen bool) ast.Declaration {
	switch p.currToken.Type {
	case token.ASTERISK:
		pointer := &ast.Pointer{PointsTo: decl}
		for {
			if p.peekTokenIs(token.CONST) {
				pointer.Const = true
			} else if p.peekTokenIs(token.VOLATILE) {
				pointer.Volatile = true
			} else {
				break
			}

			p.nextToken()
		}

		if !p.peekTokenIs(token.IDENTIFIER, token.LPAREN, token.ASTERISK) {
			return decl
		}

		p.nextToken()
		return p.parseDeclaratorLeft(pointer, insideParen)
	case token.LPAREN:
		if !p.expectPeek(token.IDENTIFIER, token.LPAREN, token.ASTERISK) {
			return nil
		}

		interior := p.parseDeclaratorLeft(decl, true)
		right := p.parseDeclaratorRight(decl, insideParen)

		// We need to insert what we parsed from the right into the
		// declaration tree
		t := interior
		for t.Type() != decl {
			t = t.Type()
		}

		t.SetType(right)

		return interior

	case token.IDENTIFIER:
		name := p.currToken.Literal
		right := p.parseDeclaratorRight(decl, insideParen)
		if fnDecl, ok := right.(*ast.FunctionDeclaration); ok {
			fnDecl.Name = name
			return fnDecl
		}
		return &ast.VariableDeclaration{VarType: right,
			Name: name}
	default:
		p.genericError(fmt.Sprintf("internal bug: called parseDeclaratorLeft with unexpected token: %s",
			p.currToken.Literal))
		return nil
	}
}

func (p *Parser) parseDeclaratorRight(decl ast.Declaration, insideParen bool) ast.Declaration {
	switch p.peekToken.Type {
	case token.LSQUARE:
		p.nextToken()
		// TODO parse expression for array size
		var expr ast.Expression
		if !p.peekTokenIs(token.RSQUARE) {
			expr = p.parseExpression(LOWEST)
		}

		if !p.expectPeek(token.RSQUARE) {
			return nil
		}

		right := p.parseDeclaratorRight(decl, insideParen)
		return &ast.Array{ArrayOf: right, ArraySize: expr}
	case token.LPAREN: // function declaration
		p.nextToken()

		params := []ast.Declaration{}
		for !p.peekTokenIs(token.RPAREN) && !p.peekTokenIs(token.EOF) {
			p.nextToken()
			param := p.parseFunctionParam()
			if param == nil {
				return nil
			}

			params = append(params, param)
			if p.peekTokenIs(token.COMMA) {
				p.nextToken()
			}
		}

		if !p.expectPeek(token.RPAREN) {
			return nil
		}

		fnDecl := &ast.FunctionDeclaration{ReturnType: decl, Parameters: params}
		return p.parseDeclaratorRight(fnDecl, insideParen)
	case token.RPAREN:
		if insideParen {
			p.nextToken()
		}
		return decl
	default:
		return decl
	}
}

func (p *Parser) parseFunctionParam() ast.Declaration {
	typeSpec := p.parseTypeSpecificiation()
	if len(typeSpec.Name) == 0 {
		p.genericError("expected parameter declarator")
		return nil
	}
	if !p.peekTokenIs(token.IDENTIFIER, token.ASTERISK, token.LPAREN) {
		return typeSpec
	}

	p.nextToken()
	return p.parseDeclaratorLeft(typeSpec, false)
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

		d := p.parseDeclaratorLeft(typeSpec, false)
		switch decl := d.(type) {
		case *ast.VariableDeclaration:
			if p.peekTokenIs(token.ASSIGN) { // also define the variable
				p.nextToken()
				p.nextToken()
				decl.Definition = p.parseExpression(LOWEST)
			}
		case *ast.FunctionDeclaration:
		}

		decls = append(decls, d)

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
