package parser

import (
	"diff/lexer"
	"fmt"
)

type Node interface {
	String() string
}

type BinOpNode struct {
	Op *lexer.Token

	Left  Node
	Right Node
}

func (n *BinOpNode) String() string {

	return fmt.Sprintf("(%s %s %s)",
		n.Left.String(), n.Op.String(), n.Right.String())
}

type UnOpNode struct {
	Op *lexer.Token

	Expr Node
}

func (n *UnOpNode) String() string {

	return fmt.Sprintf("%s(%s)",
		n.Op.String(), n.Expr.String())

}

type ValNode struct {
	Val *lexer.Token
}

func (n *ValNode) String() string {

	return n.Val.String()
}
