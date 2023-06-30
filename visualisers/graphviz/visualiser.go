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
	Nodes []gvNode
	Verts []gvVert

	deduplicateNum uint
}

func (c *gvCtx) addNode(ID, Val string) string {
	for _, n := range c.Nodes {
		if n.ID == ID {
			ID = n.ID + fmt.Sprintf("%d", c.deduplicateNum)
			c.deduplicateNum += 1
			break
		}
	}
	c.Nodes = append(c.Nodes, gvNode{ID: ID, Value: Val})
	return ID
}

func (c *gvCtx) addVert(From, To string) {
	for _, v := range c.Verts {
		if v.From == From && v.To == To {
			return
		}
	}
	c.Verts = append(c.Verts, gvVert{From: From, To: To})
}

func (c *gvCtx) walkBinOpNode(n *parser.BinOpNode) string {

	realID := c.addNode(n.Signature(), n.Op.String())

	c.addVert(realID, c.walkNode(n.Left))
	c.addVert(realID, c.walkNode(n.Right))

	return realID
}

func (c *gvCtx) walkUnOpNode(n *parser.UnOpNode) string {

	realID := c.addNode(n.Signature(), n.Op.String())

	c.addVert(realID, c.walkNode(n.Expr))

	return realID
}

func (c *gvCtx) walkNumNode(n *parser.NumNode) string {
	return c.addNode(n.Signature(), fmt.Sprintf("%.2f", n.Val))
}

func (c *gvCtx) walkConstNode(n *parser.ConstNode) string {
	return c.addNode(n.Signature(), n.Val.String())
}

func (c *gvCtx) walkVarNode(n *parser.VarNode) string {
	return c.addNode(n.Signature(), n.Val)
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
	default:
		panic("unknown node type")
	}
}

//go:embed *.dot.tmpl
var f embed.FS

func GenGraphViz(root parser.ASTNode) []byte {

	ctx := gvCtx{}
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
