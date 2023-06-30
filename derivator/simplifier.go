package derivator

import (
	"diff/lexer"
	"diff/parser"
	"fmt"
	"log"
)

func arithmOp(l, r float64, op lexer.TokenName) float64 {
	switch op {
	case lexer.Plus:
		return l + r
	case lexer.Minus:
		return l - r
	case lexer.Mult:
		return l * r
	case lexer.Div:
		return l / r
	default:
		panic("unknown operator")
	}
}

func (d *Derivator) simlifyPureNumNode(n *parser.BinOpNode, l, r *parser.NumNode) parser.ASTNode {
	switch n.Op {
	case lexer.Plus, lexer.Minus, lexer.Mult:
		return d.m.NewNumNode(arithmOp(l.Val, r.Val, n.Op))

	case lexer.Div:
		if r.Val == 0 {
			fmt.Println("warning: division by zero")
			return n
		}
		return d.m.NewNumNode(arithmOp(l.Val, r.Val, n.Op))
	case lexer.Pow:
		if l.Val == 0 {
			if r.Val == 0 {
				return d.m.NewNumNode(1)
			}
			return d.m.NewNumNode(0)
		}
		if r.Val == 0 {
			return d.m.NewNumNode(1)
		}
		return n
	}

	return n
}

func (d *Derivator) simplifyBinOpNode(n *parser.BinOpNode) parser.ASTNode {

	delete(d.m, n.String())

	n.Left = d.simplifyExpr(n.Left)
	n.Right = d.simplifyExpr(n.Right)

	n = d.m.GetOrCreateNode(n).(*parser.BinOpNode)

	l, hasNumValL := n.Left.(*parser.NumNode)
	r, hasNumValR := n.Right.(*parser.NumNode)

	if hasNumValL && hasNumValR {

		return d.simlifyPureNumNode(n, l, r)
	}

	switch n.Op {
	case lexer.Plus:
		if hasNumValL && l.Val == 0 {
			return n.Right
		}

		if hasNumValR && r.Val == 0 {
			return n.Left
		}

	case lexer.Minus:
		if hasNumValL && l.Val == 0 {
			return d.m.NewUnOpNode(lexer.Minus, n.Right)
		}
		if hasNumValR && r.Val == 0 {
			return n.Left
		}

	case lexer.Mult:
		if hasNumValL && l.Val == 0 {
			return d.m.NewNumNode(0)
		}
		if hasNumValL && l.Val == 1 {
			return n.Right
		}
		if hasNumValR && r.Val == 0 {
			return d.m.NewNumNode(0)
		}
		if hasNumValR && r.Val == 1 {
			return n.Left
		}

	case lexer.Div:
		if hasNumValL && l.Val == 0 {
			return d.m.NewNumNode(0)
		}
		if hasNumValR && r.Val == 0 {
			fmt.Println("warning: division by zero")
		}
		if hasNumValR && r.Val == 1 {
			return n.Left
		}

	case lexer.Pow:
		if hasNumValL && l.Val == 0 {
			return d.m.NewNumNode(0)
		}
		if hasNumValL && l.Val == 1 {
			return d.m.NewNumNode(1)
		}
		if hasNumValR && r.Val == 0 {
			return d.m.NewNumNode(1)
		}
		if hasNumValR && r.Val == 1 {
			return n.Left
		}
	}

	return n

}

func (d *Derivator) simplifyUnOpNode(n *parser.UnOpNode) parser.ASTNode {

	delete(d.m, n.String())

	n.Expr = d.simplifyExpr(n.Expr)

	n = d.m.GetOrCreateNode(n).(*parser.UnOpNode)

	switch n.Op {
	case lexer.Plus:
		return n.Expr

	case lexer.Minus:
		switch ch := n.Expr.(type) {
		case *parser.UnOpNode:
			if ch.Op == lexer.Minus {
				return ch.Expr
			}
		case *parser.NumNode:
			return d.m.NewNumNode(-ch.Val)
		}

	case lexer.Ln:
		switch ch := n.Expr.(type) {
		case *parser.ConstNode:
			if ch.Val == lexer.Euler {
				return d.m.NewNumNode(1)
			}
		case *parser.NumNode:
			if ch.Val == 1 {
				return d.m.NewNumNode(0)
			}
		}
	}

	return n
}

func (d *Derivator) simplifyExpr(n parser.ASTNode) parser.ASTNode {

	log.Printf("started simplifying %s", n)

	if n.GetProps().Simplified != nil {
		log.Printf("ended simplifying %s  got: %s (cached)",
			n, n.GetProps().Simplified)
		return n.GetProps().Simplified
	}

	var smpl parser.ASTNode

	switch nt := n.(type) {
	case *parser.BinOpNode:
		smpl = d.simplifyBinOpNode(nt)

	case *parser.UnOpNode:
		smpl = d.simplifyUnOpNode(nt)

	case *parser.VarNode, *parser.NumNode, *parser.ConstNode:
		smpl = nt

	default:
		panic("unknown node type")
	}

	n.GetProps().Simplified = smpl
	smpl.GetProps().Simplified = smpl

	log.Printf("ended simplifying %s  got: %s",
		n, smpl)

	return smpl
}
