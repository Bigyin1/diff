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

	var der parser.ASTNode

	switch nt := n.(type) {
	case *parser.BinOpNode:
		der = d.binOpNodeWalker(nt)
	case *parser.UnOpNode:
		der = d.unOpNodeWalker(nt)
	case *parser.NumNode:
		der = d.numNodeWalker(nt)
	case *parser.ConstNode:
		der = d.constNodeWalker(nt)
	case *parser.VarNode:
		der = d.varNodeWalker(nt)
	default:
		panic("unknown node type")
	}

	d.lv.BeginEq()
	d.lv.GenTexForNode(&parser.DerivNode{Expr: n})
	d.lv.GenEqu()

	props.Computed = d.simplifyExpr(der)

	d.lv.GenTexForNode(props.Computed)
	d.lv.EndEq()

	log.Printf("ended deriving %s  got: %s",
		n, n.GetProps().Computed)

	return props.Computed
}

func (d *Derivator) binOpNodeWalker(n *parser.BinOpNode) parser.ASTNode {

	switch n.Op {
	case lexer.Plus, lexer.Minus:

		derL := d.derivNode(n.Left)
		derR := d.derivNode(n.Right)

		return d.m.NewBinOpNode(
			n.Op,
			derL,
			derR,
		)

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

		return d.m.NewBinOpNode(
			lexer.Plus,
			multL,
			multR,
		)

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

		return d.m.NewBinOpNode(
			lexer.Div,
			minus,
			square,
		)

	case lexer.Pow:
		isBaseHasVar := d.nodeHasVariable(n.Left)
		isExpHasVar := d.nodeHasVariable(n.Right)

		if isBaseHasVar && isExpHasVar {
			panic("unimplemented")
		}
		if isBaseHasVar {
			return d.buildBasePowDeriv(n)
		}
		if isExpHasVar {
			return d.buildExpPowDeriv(n)
		}
		return d.m.NewNumNode(0)

	default:
		panic("unimplemented")

	}
}

func (d *Derivator) unOpNodeWalker(n *parser.UnOpNode) parser.ASTNode {

	switch n.Op {
	case lexer.Sin:
		f := d.m.NewUnOpNode(lexer.Cos, n.Expr)

		der := d.derivNode(n.Expr)

		return d.m.NewBinOpNode(
			lexer.Mult,
			f,
			der,
		)

	case lexer.Cos:
		f := d.m.NewUnOpNode(lexer.Sin, n.Expr)
		m := d.m.NewUnOpNode(lexer.Minus, f)

		der := d.derivNode(n.Expr)

		return d.m.NewBinOpNode(
			lexer.Mult,
			m,
			der,
		)

	case lexer.Ln:

		div := d.m.NewBinOpNode(
			lexer.Div,
			d.m.NewNumNode(1),
			n.Expr,
		)

		der := d.derivNode(n.Expr)

		return d.m.NewBinOpNode(
			lexer.Mult,
			div,
			der,
		)

	case lexer.Minus:

		der := d.derivNode(n.Expr)

		return d.m.NewUnOpNode(
			lexer.Minus,
			der,
		)
	}

	return n

}

func (d *Derivator) numNodeWalker(n *parser.NumNode) parser.ASTNode {
	return d.m.NewNumNode(0)
}

func (d *Derivator) constNodeWalker(n *parser.ConstNode) parser.ASTNode {
	return d.m.NewNumNode(0)
}

func (d *Derivator) varNodeWalker(n *parser.VarNode) parser.ASTNode {

	if d.variable == n.Val {
		return d.m.NewNumNode(1)
	}
	return d.m.NewNumNode(0)
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
