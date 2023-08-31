package derivator

import (
	"diff/lexer"
	"diff/parser"
	"log"
)

func (d *Derivator) derivNode(n parser.ASTNode) parser.ASTNode {

	log.Printf("started deriving %s", n)

	props := n.GetProps()

	if props.Computed != nil {
		log.Printf("ended deriving %s  got: %s (cached)",
			n, n.GetProps().Computed)
		return props.Computed
	}

	// TODO: possible to implement unused node deletition(but, I suppose, not in Go)

	n.Visit(d)

	d.lv.BeginEq()
	d.lv.GenTexForNode(&parser.DerivNode{Expr: n})
	d.lv.GenEqu()

	props.Computed = d.simplifyExpr(d.visitorResNode)

	d.lv.GenTexForNode(props.Computed)
	d.lv.EndEq()

	log.Printf("ended deriving %s  got: %s",
		n, n.GetProps().Computed)

	return props.Computed
}

func (d *Derivator) VisitBinOp(n *parser.BinOpNode) {

	switch n.Op {
	case lexer.Plus, lexer.Minus:

		derL := d.derivNode(n.Left)
		derR := d.derivNode(n.Right)

		d.visitorResNode = d.m.NewBinOpNode(
			n.Op,
			derL,
			derR,
		)
		return

	case lexer.Mult:
		derL := d.derivNode(n.Left)
		derR := d.derivNode(n.Right)

		multL := d.m.NewBinOpNode(
			lexer.Mult,
			derL,
			n.Right,
		)

		multR := d.m.NewBinOpNode(
			lexer.Mult,
			derR,
			n.Left,
		)

		d.visitorResNode = d.m.NewBinOpNode(
			lexer.Plus,
			multL,
			multR,
		)
		return

	case lexer.Div:
		derL := d.derivNode(n.Left)
		derR := d.derivNode(n.Right)

		multL := d.m.NewBinOpNode(
			lexer.Mult,
			derL,
			n.Right,
		)

		multR := d.m.NewBinOpNode(
			lexer.Mult,
			derR,
			n.Left,
		)

		minus := d.m.NewBinOpNode(
			lexer.Minus,
			multL,
			multR,
		)

		square := d.m.NewBinOpNode(
			lexer.Pow,
			n.Right,
			d.m.NewNumNode(2),
		)

		d.visitorResNode = d.m.NewBinOpNode(
			lexer.Div,
			minus,
			square,
		)
		return

	case lexer.Pow:
		isBaseHasVar := d.nodeHasVariable(n.Left)
		isExpHasVar := d.nodeHasVariable(n.Right)

		if isBaseHasVar && isExpHasVar {
			d.visitorResNode = d.buildMixedPowDeriv(n)
		} else if isBaseHasVar {
			d.visitorResNode = d.buildBasePowDeriv(n)
		} else if isExpHasVar {
			d.visitorResNode = d.buildExpPowDeriv(n)
		} else {
			d.visitorResNode = d.m.NewNumNode(0)
		}

	default:
		panic("unimplemented")

	}
}

func (d *Derivator) VisitUnOpNode(n *parser.UnOpNode) {

	switch n.Op {
	case lexer.Sin:
		f := d.m.NewUnOpNode(lexer.Cos, n.Expr)

		der := d.derivNode(n.Expr)

		d.visitorResNode = d.m.NewBinOpNode(
			lexer.Mult,
			f,
			der,
		)
		return

	case lexer.Cos:
		f := d.m.NewUnOpNode(lexer.Sin, n.Expr)
		m := d.m.NewUnOpNode(lexer.Minus, f)

		der := d.derivNode(n.Expr)

		d.visitorResNode = d.m.NewBinOpNode(
			lexer.Mult,
			m,
			der,
		)
		return

	case lexer.Tg:
		f := d.m.NewUnOpNode(lexer.Cos, n.Expr)
		pow := d.m.NewBinOpNode(
			lexer.Pow,
			f,
			d.m.NewNumNode(2),
		)
		div := d.m.NewBinOpNode(
			lexer.Div,
			d.m.NewNumNode(1),
			pow,
		)

		der := d.derivNode(n.Expr)

		d.visitorResNode = d.m.NewBinOpNode(
			lexer.Mult,
			div,
			der,
		)
		return

	case lexer.Ln:

		div := d.m.NewBinOpNode(
			lexer.Div,
			d.m.NewNumNode(1),
			n.Expr,
		)

		der := d.derivNode(n.Expr)

		d.visitorResNode = d.m.NewBinOpNode(
			lexer.Mult,
			div,
			der,
		)
		return

	case lexer.Minus:

		der := d.derivNode(n.Expr)

		d.visitorResNode = d.m.NewUnOpNode(
			lexer.Minus,
			der,
		)
		return

	default:
		panic("inimplemented")
	}

}

func (d *Derivator) VisitNumNode(n *parser.NumNode) {
	d.visitorResNode = d.m.NewNumNode(0)
}

func (d *Derivator) VisitConstNode(n *parser.ConstNode) {
	d.visitorResNode = d.m.NewNumNode(0)
}

func (d *Derivator) VisitVarNode(n *parser.VarNode) {

	if d.variable == n.Val {
		d.visitorResNode = d.m.NewNumNode(1)
	} else {
		d.visitorResNode = d.m.NewNumNode(0)
	}
}

func (d *Derivator) VisitDerivNode(*parser.DerivNode) {

}

func (d *Derivator) nodeHasVariable(n parser.ASTNode) bool {
	switch nt := n.(type) {
	case *parser.BinOpNode:
		return d.nodeHasVariable(nt.Left) || d.nodeHasVariable(nt.Right)

	case *parser.UnOpNode:
		return d.nodeHasVariable(nt.Expr)

	case *parser.VarNode:
		if nt.Val == d.variable {
			return true
		}
		return false

	default:
		return false
	}
}

func (d *Derivator) buildMixedPowDeriv(n *parser.BinOpNode) parser.ASTNode {
	ln := d.m.NewUnOpNode(lexer.Ln, n.Left)
	mult := d.m.NewBinOpNode(
		lexer.Mult,
		ln,
		n.Right,
	)

	der := d.derivNode(mult)

	return d.m.NewBinOpNode(
		lexer.Mult,
		n,
		der,
	)
}

func (d *Derivator) buildBasePowDeriv(n *parser.BinOpNode) parser.ASTNode {

	newExp := d.m.NewBinOpNode(
		lexer.Minus,
		n.Right,
		d.m.NewNumNode(1),
	)

	newPow := d.m.NewBinOpNode(
		lexer.Pow,
		n.Left,
		newExp,
	)

	derFunc := d.m.NewBinOpNode(
		lexer.Mult,
		n.Right,
		newPow,
	)

	der := d.derivNode(n.Left)

	return d.m.NewBinOpNode(
		lexer.Mult,
		derFunc,
		der,
	)
}

func (d *Derivator) buildExpPowDeriv(n *parser.BinOpNode) parser.ASTNode {

	ln := d.m.NewUnOpNode(
		lexer.Ln,
		n.Left,
	)

	derFunc := d.m.NewBinOpNode(
		lexer.Mult,
		n,
		ln,
	)

	der := d.derivNode(n.Right)

	return d.m.NewBinOpNode(
		lexer.Mult,
		derFunc,
		der,
	)
}
