package parser

import (
	"fmt"
	"strconv"

	"github.com/ZeroBl21/go-monkey/src/ast"
)

// TODO: Handle expression precedence
func (p *Parser) parseExpression(_ BindingPower) ast.Expression {
	prefix := p.prefixParseFn[p.curToken.Type]
	if prefix == nil {
		p.handlerError(p.curToken.Type)
		return nil
	}

	leftExp := prefix()

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
