package derivator

import (
	"diff/lexer"
	"diff/parser"
	"fmt"
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

func simplifyBinOpNode(n *parser.BinOpNode) parser.ASTNode {

	n.Left = simplifyExpr(n.Left)
	n.Right = simplifyExpr(n.Right)

	l, hasNumValL := n.Left.(*parser.NumNode)
	r, hasNumValR := n.Right.(*parser.NumNode)

	if !hasNumValL && !hasNumValR {
		return n
	}

	if hasNumValL && hasNumValR {

		switch n.Op {
		case lexer.Plus, lexer.Minus, lexer.Mult:
			root := parser.NumNode{
				Val: arithmOp(l.Val, r.Val, n.Op),
			}
			return &root
		case lexer.Div:
			if r.Val == 0 {
				fmt.Println("error: division by zero")
				return n
			}
			root := parser.NumNode{
				Val: arithmOp(l.Val, r.Val, n.Op),
			}
			return &root
		case lexer.Pow:
			if l.Val == 0 {
				if r.Val == 0 {
					return &parser.NumNode{Val: 1}
				}
				return &parser.NumNode{Val: 0}
			}
			if r.Val == 0 {
				return &parser.NumNode{Val: 1}
			}
			return n
		}
	}
	if hasNumValL {

		switch n.Op {
		case lexer.Plus:
			if l.Val == 0 {
				return n.Right
			}
			return n
		case lexer.Minus:
			if l.Val == 0 {
				root := parser.UnOpNode{
					Op:   lexer.Minus,
					Expr: n.Right,
				}
				return &root
			}
			return n
		case lexer.Mult:
			if l.Val == 0 {
				return &parser.NumNode{Val: 0}
			}
			if l.Val == 1 {
				return n.Right
			}
			return n
		case lexer.Div:
			if l.Val == 0 {
				return &parser.NumNode{Val: 0}
			}
			return n
		case lexer.Pow:
			if l.Val == 0 {
				return &parser.NumNode{Val: 0}
			}
			if l.Val == 1 {
				return &parser.NumNode{Val: 1}
			}
		}
	}
	if hasNumValR {

		switch n.Op {
		case lexer.Plus:
			if r.Val == 0 {
				return n.Left
			}
			return n
		case lexer.Minus:
			if r.Val == 0 {
				return n.Left
			}
			return n
		case lexer.Mult:
			if r.Val == 0 {
				return &parser.NumNode{Val: 0}
			}
			if r.Val == 1 {
				return n.Left
			}
			return n
		case lexer.Div:
			if r.Val == 0 {
				fmt.Println("error: division by zero")
			}
			return n
		case lexer.Pow:
			if r.Val == 0 {
				return &parser.NumNode{Val: 1}
			}
			if r.Val == 1 {
				return n.Left
			}
			return n
		}
	}

	return n

}

func simplifyExpr(n parser.ASTNode) parser.ASTNode {

	switch nt := n.(type) {
	case *parser.BinOpNode:
		return simplifyBinOpNode(nt)
	case *parser.UnOpNode:
		nt.Expr = simplifyExpr(nt.Expr)
		return nt
	case *parser.DerivNode:
		nt.Func = simplifyExpr(nt.Func)
		return nt
	case *parser.VarNode, *parser.NumNode, *parser.ConstNode:
		return nt
	default:
		panic("unknown node type")
	}
}
