package derivator

import (
	"diff/lexer"
	"diff/parser"
	"fmt"
)

func (d *Derivator) derivNodeWalker(n *parser.DerivNode) parser.ASTNode {

	switch nt := n.Func.(type) {
	case *parser.BinOpNode:
		return d.binOpNodeWalker(nt)
	case *parser.UnOpNode:
		return d.unOpNodeWalker(nt)
	case *parser.NumNode:
		return d.numNodeWalker(nt)
	case *parser.ConstNode:
		return d.constNodeWalker(nt)
	case *parser.VarNode:
		return d.varNodeWalker(nt)
	default:
		panic("unknown node type")
	}
}

func (d *Derivator) binOpNodeWalker(n *parser.BinOpNode) parser.ASTNode {
	switch n.Op {
	case lexer.Plus, lexer.Minus:

		derL := parser.DerivNode{Func: n.Left}
		derR := parser.DerivNode{Func: n.Right}

		root := parser.BinOpNode{
			Op:    n.Op,
			Left:  &derL,
			Right: &derR,
		}

		fmt.Println(root.String())

		root.Left = d.derivNodeWalker(&derL)
		root.Right = d.derivNodeWalker(&derR)

		return simplifyExpr(&root)

	case lexer.Mult:
		derL := parser.DerivNode{Func: n.Left}
		derR := parser.DerivNode{Func: n.Right}

		root := parser.BinOpNode{
			Op: lexer.Plus,
		}

		multL := parser.BinOpNode{
			Op:    lexer.Mult,
			Left:  &derL,
			Right: n.Right,
		}
		multR := parser.BinOpNode{
			Op:    lexer.Mult,
			Left:  &derR,
			Right: n.Left,
		}

		root.Left = &multL
		root.Right = &multR

		fmt.Println(root.String())

		multL.Left = d.derivNodeWalker(&derL)
		multR.Left = d.derivNodeWalker(&derR)

		return simplifyExpr(&root)
	case lexer.Div:
		derL := parser.DerivNode{Func: n.Left}
		derR := parser.DerivNode{Func: n.Right}

		minus := parser.BinOpNode{
			Op: lexer.Minus,
		}

		multL := parser.BinOpNode{
			Op:    lexer.Mult,
			Left:  &derL,
			Right: n.Right,
		}
		multR := parser.BinOpNode{
			Op:    lexer.Mult,
			Left:  &derR,
			Right: n.Left,
		}

		minus.Left = &multL
		minus.Right = &multR

		square := parser.BinOpNode{
			Op:    lexer.Pow,
			Left:  n.Right,
			Right: &parser.NumNode{Val: 2},
		}

		root := parser.BinOpNode{
			Op:    lexer.Div,
			Left:  &minus,
			Right: &square,
		}

		fmt.Println(root.String())

		multL.Left = d.derivNodeWalker(&derL)
		multR.Left = d.derivNodeWalker(&derR)

		return simplifyExpr(&root)

	case lexer.Pow:
		isBaseHasVar := d.nodeHasVariable(n.Left)
		isExpHasVar := d.nodeHasVariable(n.Right)
		if isBaseHasVar && isExpHasVar {
			panic("unimplemented")
		}
		if isBaseHasVar {
			return simplifyExpr(d.buildBasePowDeriv(n))
		}
		if isExpHasVar {
			return simplifyExpr(d.buildExpPowDeriv(n))
		}
		return &parser.NumNode{Val: 0}

	default:
		panic("unimplemented")

	}
}

func (d *Derivator) unOpNodeWalker(n *parser.UnOpNode) parser.ASTNode {

	switch n.Op {
	case lexer.Sin:
		f := parser.UnOpNode{Op: lexer.Cos, Expr: n.Expr}

		der := &parser.DerivNode{Func: n.Expr}
		root := parser.BinOpNode{
			Op:   lexer.Mult,
			Left: &f, Right: der}

		fmt.Println(root.String())

		root.Right = d.derivNodeWalker(der)
		return simplifyExpr(&root)

	case lexer.Cos:
		f := parser.UnOpNode{Op: lexer.Sin, Expr: n.Expr}
		m := parser.UnOpNode{Op: lexer.Minus, Expr: &f}

		der := &parser.DerivNode{Func: n.Expr}
		root := parser.BinOpNode{
			Op:   lexer.Mult,
			Left: &m, Right: der,
		}

		fmt.Println(root.String())

		root.Right = d.derivNodeWalker(der)

		return simplifyExpr(&root)

	case lexer.Ln:
		div := parser.BinOpNode{
			Op:    lexer.Div,
			Left:  &parser.NumNode{Val: 1},
			Right: n.Expr,
		}

		der := &parser.DerivNode{Func: n.Expr}
		root := parser.BinOpNode{
			Op:   lexer.Mult,
			Left: &div, Right: der}

		fmt.Println(root.String())

		root.Right = d.derivNodeWalker(der)
		return simplifyExpr(&root)

	case lexer.Minus:
		der := &parser.DerivNode{Func: n.Expr}
		n.Expr = d.derivNodeWalker(der)
		return simplifyExpr(n)
	}

	return n

}

func (d *Derivator) numNodeWalker(n *parser.NumNode) parser.ASTNode {
	return &parser.NumNode{Val: 0}
}

func (d *Derivator) constNodeWalker(n *parser.ConstNode) parser.ASTNode {
	return &parser.NumNode{Val: 0}
}

func (d *Derivator) varNodeWalker(n *parser.VarNode) parser.ASTNode {

	if d.variable == n.Val {
		return &parser.NumNode{Val: 1}
	}

	return &parser.NumNode{Val: 0}
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

	newExp := parser.BinOpNode{
		Op:    lexer.Minus,
		Left:  n.Right,
		Right: &parser.NumNode{Val: 1},
	}

	newPow := parser.BinOpNode{
		Op:    lexer.Pow,
		Left:  n.Left,
		Right: &newExp,
	}

	derFunc := parser.BinOpNode{
		Op:    lexer.Mult,
		Left:  n.Right,
		Right: &newPow,
	}

	der := parser.DerivNode{Func: n.Left}
	root := parser.BinOpNode{
		Op:    lexer.Mult,
		Left:  &derFunc,
		Right: &der,
	}

	fmt.Println(root.String())

	root.Right = d.derivNodeWalker(&der)
	return &root
}

func (d *Derivator) buildExpPowDeriv(n *parser.BinOpNode) parser.ASTNode {

	ln := parser.UnOpNode{
		Op:   lexer.Ln,
		Expr: n.Left,
	}

	derFunc := parser.BinOpNode{
		Op:    lexer.Mult,
		Left:  n,
		Right: &ln,
	}

	der := parser.DerivNode{Func: n.Right}
	root := parser.BinOpNode{
		Op:    lexer.Mult,
		Left:  &derFunc,
		Right: &der,
	}

	fmt.Println(root.String())

	root.Right = d.derivNodeWalker(&der)
	return &root
}
