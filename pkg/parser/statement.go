package parser

import (
	"github.com/tjarjoura/cc/pkg/ast"
	"github.com/tjarjoura/cc/pkg/token"
)

func (p *Parser) parseDeclarationStatement(topLevel bool) *ast.DeclarationStatement {
	decls := p.parseDeclarations(topLevel)
	return &ast.DeclarationStatement{Declarations: decls}
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.currToken.Type {
	case token.RETURN:
		p.nextToken()
		returnStmt := &ast.ReturnStatement{ReturnValue: p.parseExpression(LOWEST)}
		if !p.expectPeek(token.SEMICOLON) {
			return nil
		}

		return returnStmt
	default:
		if p.currTokenIsStorageClass() || p.currTokenIsType() {
			return p.parseDeclarationStatement(false)
		}

		exprStmt := &ast.ExpressionStatement{Expression: p.parseExpression(LOWEST)}
		if !p.expectPeek(token.SEMICOLON) {
			return nil
		}

		return exprStmt
	}
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	blockStmt := &ast.BlockStatement{Statements: []ast.Statement{}}

	for !p.peekTokenIs(token.RBRACE, token.EOF) {
		p.nextToken()
		blockStmt.Statements = append(blockStmt.Statements, p.parseStatement())
	}

	if !p.expectPeek(token.RBRACE) {
		return nil
	}

	return blockStmt
}
