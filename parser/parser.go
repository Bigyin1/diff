package parser

import (
	"diff/lexer"
	"strconv"
)

type Parser struct {
	tokens       []lexer.Token
	currTokenIdx uint

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

		node = &BinOpNode{Op: op.Name, Left: node, Right: p.term()}
	}

	return node
}

func (p *Parser) term() ASTNode {
	node := p.pow()

	for p.currTokenHasName(lexer.Mult, lexer.Div) {
		op := p.getCurrToken()

		p.eatToken(op.Name)

		node = &BinOpNode{Op: op.Name, Left: node, Right: p.pow()}
	}

	return node
}

func (p *Parser) pow() ASTNode {
	node := p.factor()

	for p.currTokenHasName(lexer.Pow) {
		op := p.getCurrToken()

		p.eatToken(op.Name)

		node = &BinOpNode{Op: op.Name, Left: node, Right: p.factor()}
	}

	return node
}

func (p *Parser) factor() ASTNode {

	if p.currTokenHasName(lexer.Plus, lexer.Minus) {

		op := p.getCurrToken()
		p.eatToken(op.Name)

		return &UnOpNode{Op: op.Name, Expr: p.factor()}
	}

	if p.currTokenHasClass(lexer.ClassConst,
		lexer.ClassNumber, lexer.ClassVariable) {

		v := p.getCurrToken()
		p.eatToken(v.Name)

		if v.Class == lexer.ClassConst {
			return &ConstNode{Val: v.Name}
		}
		if v.Class == lexer.ClassNumber {

			num, _ := strconv.ParseFloat(v.Value, 64)
			return &NumNode{Val: num}
		}

		return &VarNode{Val: v.Value}
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

		node := &UnOpNode{Op: fn.Name, Expr: p.expr()}

		p.eatToken(lexer.RParen)
		return node
	}

	p.err = &UnexpTokenError{tok: p.getCurrToken()}

	return nil
}

func NewParser(toks []lexer.Token) *Parser {
	return &Parser{
		tokens: toks,
	}
}

func (p *Parser) Run() (ASTNode, error) {

	return &DerivNode{Func: p.expr()}, p.err
}
