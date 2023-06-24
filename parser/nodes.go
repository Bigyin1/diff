package parser

import (
	"diff/lexer"
	"fmt"
)

type Node interface {
	String() string
}

type BinOpNode struct {
	Op lexer.TokenName

	Left  Node
	Right Node
}

func (n *BinOpNode) String() string {
	if n.Left == nil || n.Right == nil {
		return ""
	}

	return fmt.Sprintf("(%s %s %s)",
		n.Left.String(), n.Op.String(), n.Right.String())
}

type UnOpNode struct {
	Op lexer.TokenName

	Expr Node
}

func (n *UnOpNode) String() string {
	if n.Expr == nil {
		return ""
	}

	return fmt.Sprintf("%s(%s)",
		n.Op.String(), n.Expr.String())

}

type ValNode struct {
	Type lexer.TokenName
	Val  string
}

func (n *ValNode) String() string {

	return fmt.Sprintf("%s(%s)", n.Type.String(), n.Val)
}

type DerivNode struct {
	Func Node
}

func (n *DerivNode) String() string {
	if n.Func == nil {
		return ""
	}

	return fmt.Sprintf("(%s)'", n.Func.String())
}
