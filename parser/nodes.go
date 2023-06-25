package parser

import (
	"diff/lexer"
	"fmt"
)

type NodeType = lexer.TokenName

type ASTNode interface {
	String() string
}

type BinOpNode struct {
	Op NodeType

	Left  ASTNode
	Right ASTNode
}

func (n *BinOpNode) String() string {
	if n.Left == nil || n.Right == nil {
		return ""
	}

	return fmt.Sprintf("(%s %s %s)",
		n.Left.String(), n.Op.String(), n.Right.String())
}

type UnOpNode struct {
	Op NodeType

	Expr ASTNode
}

func (n *UnOpNode) String() string {
	if n.Expr == nil {
		return ""
	}

	return fmt.Sprintf("%s(%s)",
		n.Op.String(), n.Expr.String())

}

type NumNode struct {
	Val float64
}

func (n *NumNode) String() string {

	return fmt.Sprintf("%.2f", n.Val)
}

type ConstNode struct {
	Val NodeType
}

func (n *ConstNode) String() string {

	return n.Val.String()
}

type VarNode struct {
	Val string
}

func (n *VarNode) String() string {

	return n.Val
}

type DerivNode struct {
	Func ASTNode
}

func (n *DerivNode) String() string {
	if n.Func == nil {
		return ""
	}

	return fmt.Sprintf("(%s)'", n.Func.String())
}
