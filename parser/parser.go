package parser

import (
	"diff/lexer"
	"strconv"
)

type Parser struct {
	tokens       []lexer.Token
	currTokenIdx uint

	m NodesMap

	err error
}

func (p *Parser) getCurrToken() *lexer.Token {
	return &p.tokens[p.currTokenIdx]
}

func (p *Parser) currTokenHasName(nms ...lexer.TokenName) bool {
	if p.err != nil {
		return false
	}

	currName := p.getCurrToken().Name

	for _, t := range nms {
		if t == currName {
			return true
		}
	}

	return false
}

func (p *Parser) currTokenHasClass(cls ...lexer.TokenClass) bool {
	if p.err != nil {
		return false
	}

	currClass := p.getCurrToken().Class

	for _, t := range cls {
		if t == currClass {
			return true
		}
	}

	return false
}

func (p *Parser) eatToken(exp lexer.TokenName) {
	if p.err != nil {
		return
	}

	if p.currTokenHasName(exp) {
		p.currTokenIdx++
		return
	}

	p.err = &UnexpTokenError{tok: p.getCurrToken(), exp: exp}

}

func (p *Parser) expr() ASTNode {

	node := p.term()

	for p.currTokenHasName(lexer.Plus, lexer.Minus) {
		op := p.getCurrToken()

		p.eatToken(op.Name)

		node = p.m.NewBinOpNode(
			op.Name,
			node,
			p.term(),
		)
	}

	return node
}

func (p *Parser) term() ASTNode {

	node := p.pow()

	for p.currTokenHasName(lexer.Mult, lexer.Div) {
		op := p.getCurrToken()

		p.eatToken(op.Name)

		node = p.m.NewBinOpNode(
			op.Name,
			node,
			p.pow(),
		)
	}

	return node
}

func (p *Parser) pow() ASTNode {
	node := p.factor()

	for p.currTokenHasName(lexer.Pow) {
		op := p.getCurrToken()

		p.eatToken(op.Name)

		node = p.m.NewBinOpNode(
			op.Name,
			node,
			p.factor(),
		)
	}

	return node
}

func (p *Parser) factor() ASTNode {

	if p.currTokenHasName(lexer.Plus, lexer.Minus) {
		op := p.getCurrToken()

		p.eatToken(op.Name)

		node := p.m.NewUnOpNode(
			op.Name,
			p.factor(),
		)
		return node
	}

	if p.currTokenHasClass(lexer.ClassConst,
		lexer.ClassNumber, lexer.ClassVariable) {

		v := p.getCurrToken()
		p.eatToken(v.Name)

		if v.Class == lexer.ClassConst {
			return p.m.NewConstNode(v.Name)
		}
		if v.Class == lexer.ClassNumber {
			num, _ := strconv.ParseFloat(v.Value, 64)
			return p.m.NewNumNode(num)
		}

		return NewVarNode(v.Value, p.m)
	}

	if p.currTokenHasName(lexer.LParen) {
		p.eatToken(lexer.LParen)

		node := p.expr()

		p.eatToken(lexer.RParen)
		return node
	}

	if p.currTokenHasClass(lexer.ClassFunction) {
		fn := p.getCurrToken()

		p.eatToken(fn.Name)
		p.eatToken(lexer.LParen)

		node := p.m.NewUnOpNode(
			fn.Name,
			p.expr(),
		)

		p.eatToken(lexer.RParen)
		return node
	}

	p.err = &UnexpTokenError{tok: p.getCurrToken()}

	return nil
}

func NewParser(toks []lexer.Token) *Parser {

	return &Parser{
		tokens: toks,
		m:      NewNodesMap(),
	}
}

func (p *Parser) Run() (ASTNode, NodesMap, error) {

	res := p.expr()

	p.eatToken(lexer.EOF)

	return res, p.m, p.err
}
