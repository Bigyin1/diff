package derivator

import (
	"diff/parser"
	"fmt"
)

type Derivator struct {
	root     parser.Node
	variable string
}

func NewDerivator(root parser.Node, Var string) *Derivator {
	return &Derivator{
		root:     root,
		variable: Var,
	}
}

func (d *Derivator) Run() parser.Node {

	d.root = simplifyExpr(d.root)
	fmt.Println(d.root)

	r, ok := d.root.(*parser.DerivNode)
	if !ok {
		return d.root
	}

	return d.derivNodeWalker(r)
}
