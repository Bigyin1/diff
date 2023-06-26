package graphviz

import (
	"bytes"
	"diff/parser"
	"embed"
	"fmt"
	"text/template"
)

type gvNode struct {
	ID    string
	Value string
}

type gvVert struct {
	From string
	To   string
}

type gvCtx struct {
	Nodes    []gvNode
	Verts    []gvVert
	DerivVar string

	dedublicateNum uint
}

func (c *gvCtx) addNode(ID, Val string) string {
	for _, n := range c.Nodes {
		if n.ID == ID {
			ID = n.ID + fmt.Sprintf("%d", c.dedublicateNum)
			c.dedublicateNum += 1
			break
		}
	}
	c.Nodes = append(c.Nodes, gvNode{ID: ID, Value: Val})
	return ID
}

func (c *gvCtx) addVert(From, To string) {
	c.Verts = append(c.Verts, gvVert{From: From, To: To})
}

func (c *gvCtx) walkBinOpNode(n *parser.BinOpNode) string {

	realID := c.addNode(n.ID(), n.Op.String())

	c.addVert(realID, c.walkNode(n.Left))
	c.addVert(realID, c.walkNode(n.Right))

	return realID
}

func (c *gvCtx) walkUnOpNode(n *parser.UnOpNode) string {

	realID := c.addNode(n.ID(), n.Op.String())

	c.addVert(realID, c.walkNode(n.Expr))

	return realID
}

func (c *gvCtx) walkNumNode(n *parser.NumNode) string {
	return c.addNode(n.ID(), fmt.Sprintf("%.2f", n.Val))
}

func (c *gvCtx) walkConstNode(n *parser.ConstNode) string {
	return c.addNode(n.ID(), n.Val.String())
}

func (c *gvCtx) walkVarNode(n *parser.VarNode) string {
	return c.addNode(n.ID(), n.Val)
}

func (c *gvCtx) walkDerivNode(n *parser.DerivNode) string {
	realID := c.addNode(n.ID(), fmt.Sprintf("d()/d%s", c.DerivVar))

	c.addVert(realID, c.walkNode(n.Func))

	return realID
}

func (c *gvCtx) walkNode(n parser.ASTNode) string {
	switch nt := n.(type) {
	case *parser.BinOpNode:
		return c.walkBinOpNode(nt)
	case *parser.UnOpNode:
		return c.walkUnOpNode(nt)
	case *parser.NumNode:
		return c.walkNumNode(nt)
	case *parser.ConstNode:
		return c.walkConstNode(nt)
	case *parser.VarNode:
		return c.walkVarNode(nt)
	case *parser.DerivNode:
		return c.walkDerivNode(nt)
	default:
		panic("unknown node type")
	}
}

//go:embed *.dot.tmpl
var f embed.FS

func GenGraphViz(root parser.ASTNode, derVar string) []byte {

	ctx := gvCtx{DerivVar: derVar}
	ctx.walkNode(root)

	tmpl, err := template.ParseFS(f, "*.dot.tmpl")
	if err != nil {
		panic(err)
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, &ctx)
	if err != nil {
		panic(err)
	}

	return buf.Bytes()

}
