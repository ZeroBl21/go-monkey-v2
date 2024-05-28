package parser

import (
	"fmt"
	"strconv"

	"github.com/ZeroBl21/go-monkey/src/ast"
	"github.com/ZeroBl21/go-monkey/src/token"
)

type BindingPower int

const (
	_ BindingPower = iota
	LOWEST
	EQUALS       // ==
	LESS_GREATER // > or <
	SUM          // +
	PRODUCT      // *
	PREFIX       // -X or !X
	CALL         // myFunction(X)
)

var precedences = map[token.TokenType]BindingPower{
	token.EQ:       EQUALS,
	token.NOT_EQ:   EQUALS,
	token.LT:       LESS_GREATER,
	token.RT:       LESS_GREATER,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.SLASH:    PRODUCT,
	token.ASTERISK: PRODUCT,
}

func (p *Parser) parseExpression(precedence BindingPower) ast.Expression {
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		p.noPrefixParseFn(p.curToken.Type)
		return nil
	}
	leftExp := prefix()

	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}

		p.nextToken()

		leftExp = infix(leftExp)
	}

	return leftExp
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{
		Token: p.curToken,
		Value: 0,
	}

	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}

	lit.Value = value

	return lit
}

func (p *Parser) parseBoolean() ast.Expression {
	return &ast.Boolean{
		Token: p.curToken,
		Value: p.curTokenIs(token.TRUE),
	}
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Right:    nil,
	}

	p.nextToken()

	expression.Right = p.parseExpression(PREFIX)

	return expression
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left,
	}

	precedence := p.curPrecedence()
	p.nextToken()
	expression.Right = p.parseExpression(precedence)

	return expression
}

func (p *Parser) noPrefixParseFn(t token.TokenType) {
	msg := fmt.Sprintf("No prefix parse function for token %s found", t)
	p.errors = append(p.errors, msg)
}

func (p *Parser) peekPrecedence() BindingPower {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}

	return LOWEST
}

func (p *Parser) curPrecedence() BindingPower {
	if p, ok := precedences[p.curToken.Type]; ok {
		return p
	}

	return LOWEST
}
