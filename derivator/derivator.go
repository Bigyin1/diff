package derivator

import (
	"diff/parser"
	"diff/visualizers/latex"
)

type Derivator struct {
	root     parser.ASTNode
	m        parser.NodesMap
	variable string

	lv *latex.LatexVisualiser
}

func NewDerivator(root parser.ASTNode, m parser.NodesMap, Var string) *Derivator {
	return &Derivator{
		root:     root,
		variable: Var,
		m:        m,

		lv: latex.NewLatexVisualiser(),
	}
}

func (d *Derivator) Run() parser.ASTNode {

	d.lv.BeginDoc()
	d.lv.GenTexForNode(d.root)

	d.root = d.simplifyExpr(d.root)

	d.lv.GenTexForNode(d.root)

	res := d.derivNode(d.root)

	d.lv.GenTexForNode(res)
	d.lv.EndDoc()

	return res
}
