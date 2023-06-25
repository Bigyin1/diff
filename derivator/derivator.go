package derivator

import (
	"diff/parser"
	"fmt"
)

type Derivator struct {
	root     parser.ASTNode
	variable string
}

func NewDerivator(root parser.ASTNode, Var string) *Derivator {
	return &Derivator{
		root:     root,
		variable: Var,
	}
}

func (d *Derivator) Run() parser.ASTNode {

	d.root = simplifyExpr(d.root)
	fmt.Println(d.root)

	r, ok := d.root.(*parser.DerivNode)
	if !ok {
		return d.root
	}

	return d.derivNodeWalker(r)
}
