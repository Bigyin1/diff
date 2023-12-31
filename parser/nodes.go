package parser

import (
	"diff/lexer"
	"fmt"
	"log"
)

type NodeVisiter interface {
	VisitBinOp(n *BinOpNode)
	VisitUnOpNode(n *UnOpNode)
	VisitNumNode(n *NumNode)
	VisitConstNode(n *ConstNode)
	VisitVarNode(n *VarNode)
	VisitDerivNode(n *DerivNode)
}

type NodeType = lexer.TokenName

type ASTNode interface {
	String() string
	Signature() string
	Addr() string
	GetProps() *NodeProps
	Visit(NodeVisiter)
}

type NodeProps struct {
	Computed   ASTNode
	Simplified ASTNode
	// RefCount   uint TODO: unused node deletition
}

type BinOpNode struct {
	NodeProps

	Op NodeType

	Left  ASTNode
	Right ASTNode
}

func (n *BinOpNode) Visit(nv NodeVisiter) {
	nv.VisitBinOp(n)
}

func (n *BinOpNode) String() string {
	if n.Left == nil || n.Right == nil {
		return ""
	}

	return fmt.Sprintf("(%s %s %s)",
		n.Left.String(), n.Op.String(), n.Right.String())
}

func (n *BinOpNode) Signature() string {
	return fmt.Sprintf(
		"%s%s%s",
		n.Op.String(),
		n.Left.Addr(),
		n.Right.Addr(),
	)
}

func (n *BinOpNode) Addr() string {
	return fmt.Sprintf("%p", n)
}

func (n *BinOpNode) GetProps() *NodeProps {
	return &n.NodeProps
}

func (m NodesMap) NewBinOpNode(op NodeType, left, right ASTNode) ASTNode {
	node := BinOpNode{
		Op:    op,
		Left:  left,
		Right: right,
	}

	return m.GetOrCreateNode(&node)
}

type UnOpNode struct {
	NodeProps

	Op   NodeType
	Expr ASTNode
}

func (n *UnOpNode) Visit(nv NodeVisiter) {
	nv.VisitUnOpNode(n)
}

func (n *UnOpNode) String() string {
	if n.Expr == nil {
		return ""
	}

	return fmt.Sprintf("%s(%s)",
		n.Op.String(), n.Expr.String())

}

func (n *UnOpNode) Signature() string {
	return fmt.Sprintf("%s%s", n.Op.String(), n.Expr.Addr())
}

func (n *UnOpNode) Addr() string {
	return fmt.Sprintf("%p", n)
}

func (n *UnOpNode) GetProps() *NodeProps {
	return &n.NodeProps
}

func (m NodesMap) NewUnOpNode(op NodeType, expr ASTNode) ASTNode {
	node := UnOpNode{
		Op:   op,
		Expr: expr,
	}

	return m.GetOrCreateNode(&node)
}

type NumNode struct {
	NodeProps

	Val float64
}

func (n *NumNode) Visit(nv NodeVisiter) {
	nv.VisitNumNode(n)
}

func (n *NumNode) String() string {
	return fmt.Sprintf("%g", n.Val)
}

func (n *NumNode) Signature() string {
	return fmt.Sprintf("%g", n.Val)
}

func (n *NumNode) Addr() string {
	return fmt.Sprintf("%p", n)
}

func (n *NumNode) GetProps() *NodeProps {
	return &n.NodeProps
}

func (m NodesMap) NewNumNode(v float64) ASTNode {
	node := NumNode{
		Val: v,
	}

	return m.GetOrCreateNode(&node)
}

type ConstNode struct {
	NodeProps

	Val NodeType
}

func (n *ConstNode) Visit(nv NodeVisiter) {
	nv.VisitConstNode(n)
}

func (n *ConstNode) String() string {

	return n.Val.String()
}

func (n *ConstNode) Signature() string {
	return n.Val.String()
}

func (n *ConstNode) Addr() string {
	return fmt.Sprintf("%p", n)
}

func (n *ConstNode) GetProps() *NodeProps {
	return &n.NodeProps
}

func (m NodesMap) NewConstNode(val NodeType) ASTNode {
	node := ConstNode{
		Val: val,
	}

	return m.GetOrCreateNode(&node)
}

type VarNode struct {
	NodeProps

	Val string
}

func (n *VarNode) Visit(nv NodeVisiter) {
	nv.VisitVarNode(n)
}

func (n *VarNode) String() string {

	return n.Val
}

func (n *VarNode) Signature() string {
	return n.Val
}

func (n *VarNode) Addr() string {
	return fmt.Sprintf("%p", n)
}

func (n *VarNode) GetProps() *NodeProps {
	return &n.NodeProps
}

func NewVarNode(val string, m NodesMap) ASTNode {
	node := VarNode{
		Val: val,
	}

	return m.GetOrCreateNode(&node)
}

// DerivNode is an auxillary type for latex visualiser only
type DerivNode struct {
	Expr ASTNode
}

func (n *DerivNode) Visit(nv NodeVisiter) {
	nv.VisitDerivNode(n)
}

func (n *DerivNode) String() string {
	return ""
}

func (n *DerivNode) Signature() string {
	return ""
}

func (n *DerivNode) Addr() string {
	return ""
}

func (n *DerivNode) GetProps() *NodeProps {
	return nil
}

type NodesMap map[string]ASTNode

func NewNodesMap() NodesMap {
	return make(NodesMap)
}

func (m NodesMap) GetOrCreateNode(node ASTNode) ASTNode {
	log.Printf("try to get node: %s", node.String())

	n, ok := m[node.Signature()]
	if !ok {
		m[node.Signature()] = node

		log.Printf("new node: %s", node.String())

		return node
	}

	log.Printf("got cached node: %s", node.String())

	return n
}
