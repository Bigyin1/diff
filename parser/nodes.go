package parser

import (
	"diff/lexer"
	"fmt"
)

type NodeType = lexer.TokenName

type ASTNode interface {
	String() string
	ID() string
	Copy() ASTNode
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

func (n *BinOpNode) ID() string {
	return fmt.Sprintf("%p", n)[1:]
}

func (n *BinOpNode) Copy() ASTNode {
	newNode := *n

	return &newNode
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

func (n *UnOpNode) Copy() ASTNode {
	newNode := *n

	return &newNode
}

func (n *UnOpNode) ID() string {
	return fmt.Sprintf("%p", n)[1:]
}

type NumNode struct {
	Val float64
}

func (n *NumNode) String() string {

	return fmt.Sprintf("%.2f", n.Val)
}

func (n *NumNode) ID() string {
	return fmt.Sprintf("%p", n)[1:]
}

func (n *NumNode) Copy() ASTNode {
	newNode := *n

	return &newNode
}

type ConstNode struct {
	Val NodeType
}

func (n *ConstNode) String() string {

	return n.Val.String()
}

func (n *ConstNode) ID() string {
	return fmt.Sprintf("%p", n)[1:]
}

func (n *ConstNode) Copy() ASTNode {
	newNode := *n

	return &newNode
}

type VarNode struct {
	Val string
}

func (n *VarNode) String() string {

	return n.Val
}

func (n *VarNode) ID() string {
	return fmt.Sprintf("%p", n)[1:]
}

func (n *VarNode) Copy() ASTNode {
	newNode := *n

	return &newNode
}

type DerivNode struct {
	Func ASTNode
}

func (n *DerivNode) ID() string {
	return fmt.Sprintf("%p", n)[1:]
}

func (n *DerivNode) String() string {
	if n.Func == nil {
		return ""
	}

	return fmt.Sprintf("(%s)'", n.Func.String())
}

func (n *DerivNode) Copy() ASTNode {
	newNode := *n

	return &newNode
}
