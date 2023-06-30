package derivator

import (
	"diff/parser"
	"diff/visualisers/latex"
)

type Derivator struct {
	root     parser.ASTNode
	m        parser.NodesMap
	variable string

	lv *latex.LatexVisualiser
}

func NewDerivator(root parser.ASTNode, m parser.NodesMap, Var, Out string) (*Derivator, error) {

	lv, err := latex.NewLatexVisualiser(Out)
	if err != nil {
		return nil, err
	}

	return &Derivator{
		root:     root,
		variable: Var,
		m:        m,

		lv: lv,
	}, nil
}

func (d *Derivator) Run() (parser.ASTNode, error) {

	d.lv.BeginDoc()

	d.root = d.simplifyExpr(d.root)

	d.lv.GenStr("\nУпрощенное исходное выражение:\n")
	d.lv.BeginEq()
	d.lv.GenTexForNode(d.root)
	d.lv.EndEq()

	d.lv.GenStr("\nПромежуточные вычисления:\n")

	res := d.derivNode(d.root)

	d.lv.GenStr("\nПолучим:\n")

	d.lv.BeginEq()
	d.lv.GenTexForNode(res)
	d.lv.EndEq()

	err := d.lv.EndDoc()

	return res, err
}
