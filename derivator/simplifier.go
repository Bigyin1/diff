package derivator

import (
	"diff/lexer"
	"diff/parser"
	"fmt"
	"strconv"
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

func simplifyBinOpNode(n *parser.BinOpNode) parser.Node {

	n.Left = simplifyExpr(n.Left)
	n.Right = simplifyExpr(n.Right)

	l, hasNumValL := n.Left.(*parser.ValNode)
	r, hasNumValR := n.Right.(*parser.ValNode)
	if !hasNumValL && !hasNumValR {
		return n
	}
	var leftVal float64
	var rightVal float64

	if l != nil && l.Type != lexer.Number {
		hasNumValL = false
	} else if l != nil {
		leftVal, _ = strconv.ParseFloat(l.Val, 64)
	}

	if r != nil && r.Type != lexer.Number {
		hasNumValR = false
	} else if r != nil {
		rightVal, _ = strconv.ParseFloat(r.Val, 64)
	}

	if hasNumValL && hasNumValR {

		switch n.Op {
		case lexer.Plus, lexer.Minus, lexer.Mult:
			root := parser.ValNode{
				Type: lexer.Number,
				Val:  fmt.Sprintf("%f", arithmOp(leftVal, rightVal, n.Op)),
			}
			return &root
		case lexer.Div:
			if rightVal == 0 {
				fmt.Println("error: division by zero")
				return n
			}
			root := parser.ValNode{
				Type: lexer.Number,
				Val:  fmt.Sprintf("%f", arithmOp(leftVal, rightVal, n.Op)),
			}
			return &root
		case lexer.Pow:
			if leftVal == 0 {
				if rightVal == 0 {
					return &parser.ValNode{Type: lexer.Number, Val: "1"}
				}
				return &parser.ValNode{Type: lexer.Number, Val: "0"}
			}
			if rightVal == 0 {
				return &parser.ValNode{Type: lexer.Number, Val: "1"}
			}
			return n
		}
	}
	if hasNumValL {

		switch n.Op {
		case lexer.Plus:
			if leftVal == 0 {
				return n.Right
			}
			return n
		case lexer.Minus:
			if leftVal == 0 {
				root := parser.UnOpNode{
					Op:   lexer.Minus,
					Expr: n.Right,
				}
				return &root
			}
			return n
		case lexer.Mult:
			if leftVal == 0 {
				return &parser.ValNode{Type: lexer.Number, Val: "0"}
			}
			if leftVal == 1 {
				return n.Right
			}
			return n
		case lexer.Div, lexer.Pow:
			if leftVal == 0 {
				return &parser.ValNode{Type: lexer.Number, Val: "0"}
			}
			return n
		}
	}
	if hasNumValR {

		switch n.Op {
		case lexer.Plus:
			if rightVal == 0 {
				return n.Left
			}
			return n
		case lexer.Minus:
			if rightVal == 0 {
				return n.Left
			}
			return n
		case lexer.Mult:
			if rightVal == 0 {
				return &parser.ValNode{Type: lexer.Number, Val: "0"}
			}
			if rightVal == 1 {
				return n.Left
			}
			return n
		case lexer.Div:
			if rightVal == 0 {
				fmt.Println("error: division by zero")
			}
			return n
		case lexer.Pow:
			if rightVal == 0 {
				return &parser.ValNode{Type: lexer.Number, Val: "1"}
			}
			return n
		}
	}

	return n

}

func simplifyExpr(n parser.Node) parser.Node {

	switch nt := n.(type) {
	case *parser.BinOpNode:
		return simplifyBinOpNode(nt)
	case *parser.UnOpNode:
		nt.Expr = simplifyExpr(nt.Expr)
		return nt
	case *parser.DerivNode:
		nt.Func = simplifyExpr(nt.Func)
		return nt
	case *parser.ValNode:
		return nt
	default:
		panic("unknown node type")
	}
}
