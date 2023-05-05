package parser

import (
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
			return pointer
		}

		p.nextToken()
		return p.parseDeclaratorLeft(pointer, insideParen)
	case token.LPAREN:
		if !p.expectPeek(token.IDENTIFIER, token.LPAREN, token.ASTERISK) {
			return nil
		}

		interior := p.parseDeclaratorLeft(decl, true)
		if interior == nil {
			return nil
		}

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
		var expr ast.Expression
		if !p.peekTokenIs(token.RSQUARE) {
			p.nextToken()
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
	typeSpec := p.parseBaseType()
	if typeSpec == nil {
		return nil
	}
	if !p.peekTokenIs(token.IDENTIFIER, token.ASTERISK, token.LPAREN) {
		return typeSpec
	}

	p.nextToken()
	return p.parseDeclaratorLeft(typeSpec, false)
}

// Kind of spaghetti code but the type naming rules in C are a bit all over the place
func (p *Parser) combineTypeSpecifier(typeSpec string) string {
	switch p.currToken.Type {
	case token.SHORT:
		if typeSpec == "int" {
			return "short int"
		}
	case token.INT:
		if typeSpec == "short" || typeSpec == "short int" {
			return "short int"
		} else if typeSpec == "long" || typeSpec == "long int" {
			return "long int"
		} else if typeSpec == "long long" || typeSpec == "long long int" {
			return "long long int"
		}
	case token.LONG:
		if typeSpec == "int" {
			return "long int"
		} else if typeSpec == "long int" || typeSpec == "long" {
			return "long long int"
		} else if typeSpec == "double" {
			return "long double"
		}
	case token.DOUBLE:
		if typeSpec == "long" || typeSpec == "long double" {
			return "long double"
		}
	}

	if typeSpec != p.currToken.Literal {
		p.genericError(fmt.Sprintf(
			"conflicting type specifier %s",
			p.currToken.Literal))
		return ""

	}
	return typeSpec

}

func (p *Parser) parseBaseType() *ast.BaseType {
	typeSpec := &ast.BaseType{}
	alreadySigned, alreadyTyped := false, false

	for {
		if p.currTokenIs(token.CONST) {
			typeSpec.Const = true
		} else if p.currTokenIs(token.VOLATILE) {
			typeSpec.Volatile = true
		} else if p.currTokenIs(token.UNSIGNED) || p.currTokenIs(token.SIGNED) {
			if alreadySigned {
				if typeSpec.Signed != p.currTokenIs(token.SIGNED) {
					p.genericError(fmt.Sprintf(
						"conflicting type specifier %s",
						p.currToken.Literal))
					return nil
				}
			} else {
				alreadySigned = true
				typeSpec.Signed = p.currTokenIs(token.SIGNED)
			}
		} else if p.currTokenIsType() {
			if !alreadyTyped {
				typeSpec.Name = p.currToken.Literal
				alreadyTyped = true
			} else {
				typeSpec.Name = p.combineTypeSpecifier(typeSpec.Name)
				if typeSpec.Name == "" {
					return nil
				}
			}
		}

		if p.peekTokenIsType() || p.peekTokenIsTypeQualifier() {
			p.nextToken()
		} else {
			break
		}
	}

	if typeSpec.Name == "" {
		p.genericError("type specifier missing. implicit int is not supported by this compiler")
		return nil
	} else if typeSpec.Name == "long" || typeSpec.Name == "short" { // append "int" for consistency
		typeSpec.Name = fmt.Sprintf("%s int", typeSpec.Name)
	}

	return typeSpec
}

func (p *Parser) parseDeclarations(topLevel bool) []ast.Declaration {
	var decls = []ast.Declaration{}
	var storageClass string
	if p.currTokenIsStorageClass() {
		storageClass = p.currToken.Literal
		p.nextToken()
	}

	typeSpec := p.parseBaseType()
	if typeSpec == nil {
		return decls
	}

	for !p.peekTokenIs(token.SEMICOLON) && !p.peekTokenIs(token.EOF) {
		if !p.expectPeek(token.IDENTIFIER, token.LPAREN, token.ASTERISK) {
			return nil
		}

		d := p.parseDeclaratorLeft(typeSpec, false)
		decls = append(decls, d)

		switch decl := d.(type) {
		case *ast.VariableDeclaration:
			decl.StorageClass = storageClass
			if p.peekTokenIs(token.ASSIGN) { // also define the variable
				p.nextToken()
				p.nextToken()
				decl.Definition = p.parseExpression(LOWEST)
			}
		case *ast.FunctionDeclaration:
			decl.StorageClass = storageClass

			// if this is the first declaration, there can also be a function definition
			if len(decls) == 1 && p.peekTokenIs(token.LBRACE) {
				if !topLevel {
					p.genericError("function definition not allowed here")
					return decls
				}
				p.nextToken()
				decl.Body = p.parseBlockStatement()
				goto end // we don't allow more than one declaration if we defined a function
			}
		}

		if p.peekTokenIs(token.COMMA) {
			p.nextToken()
		}
	}

	if !p.expectPeek(token.SEMICOLON) {
		return nil
	}

end:
	return decls
}
