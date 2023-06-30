package latex

import (
	"diff/lexer"
	"diff/parser"
	_ "embed"
	"fmt"
	"os"
	"strings"
)

//go:embed header.tex
var header string

type LatexVisualiser struct {
	text *strings.Builder
	out  os.File
}

func NewLatexVisualiser(Out string) (*LatexVisualiser, error) {

	out, err := os.Create(Out)
	if err != nil {
		return nil, err
	}

	var file strings.Builder

	file.WriteString(header)

	return &LatexVisualiser{
		text: &file,
		out:  *out,
	}, nil
}

func (lv *LatexVisualiser) BeginDoc() {
	lv.text.WriteString("\n\\begin{document}\n")
}

func (lv *LatexVisualiser) EndDoc() error {
	lv.text.WriteString("\n\\end{document}\n")

	defer lv.out.Close()

	_, err := lv.out.WriteString(lv.text.String())

	return err
}

func (lv *LatexVisualiser) walkBinOpNode(n *parser.BinOpNode) string {

	res := strings.Builder{}

	dL := lv.walkNode(n.Left)
	dR := lv.walkNode(n.Right)

	switch n.Op {
	case lexer.Plus:
		res.WriteString(fmt.Sprintf("%s + %s", dL, dR))

	case lexer.Minus:
		res.WriteString(fmt.Sprintf("%s - %s", dL, dR))

	case lexer.Mult:
		switch nt := n.Right.(type) {
		case *parser.BinOpNode:
			if nt.Op == lexer.Plus || nt.Op == lexer.Minus {
				res.WriteString(fmt.Sprintf("%s \\cdot (%s)", dL, dR))
				return res.String()
			}
		}
		res.WriteString(fmt.Sprintf("%s \\cdot %s", dL, dR))

	case lexer.Div:
		res.WriteString(fmt.Sprintf("\\frac{%s}{%s}", dL, dR))

	case lexer.Pow:
		res.WriteString(fmt.Sprintf("{%s}^{%s}", dL, dR))
	}

	return res.String()

}

func (lv *LatexVisualiser) walkUnOpNode(n *parser.UnOpNode) string {
	res := strings.Builder{}

	dE := lv.walkNode(n.Expr)

	switch n.Op {
	case lexer.Plus:
		res.WriteString(fmt.Sprintf("+(%s)", dE))

	case lexer.Minus:
		res.WriteString(fmt.Sprintf("-(%s)", dE))

	case lexer.Sin:
		res.WriteString(fmt.Sprintf("\\sin(%s)", dE))

	case lexer.Cos:
		res.WriteString(fmt.Sprintf("\\cos(%s)", dE))

	case lexer.Ln:
		res.WriteString(fmt.Sprintf("\\ln(%s)", dE))
	}

	return res.String()
}

func (lv *LatexVisualiser) walkNumNode(n *parser.NumNode) string {
	res := strings.Builder{}

	res.WriteString(n.String())

	return res.String()
}

func (lv *LatexVisualiser) walkConstNode(n *parser.ConstNode) string {
	res := strings.Builder{}

	switch n.Val {
	case lexer.Pi:
		res.WriteString("\\pi")
	case lexer.Euler:
		res.WriteString("e")
	}

	return res.String()
}

func (lv *LatexVisualiser) walkVarNode(n *parser.VarNode) string {
	res := strings.Builder{}

	res.WriteString(n.String())

	return res.String()
}

func (lv *LatexVisualiser) walkDerivNode(n *parser.DerivNode) string {
	res := strings.Builder{}

	res.WriteString(fmt.Sprintf("(%s)'", lv.walkNode(n.Expr)))

	return res.String()
}

func (lv *LatexVisualiser) walkNode(n parser.ASTNode) string {
	switch nt := n.(type) {
	case *parser.BinOpNode:
		return lv.walkBinOpNode(nt)

	case *parser.UnOpNode:
		return lv.walkUnOpNode(nt)

	case *parser.NumNode:
		return lv.walkNumNode(nt)

	case *parser.ConstNode:
		return lv.walkConstNode(nt)

	case *parser.VarNode:
		return lv.walkVarNode(nt)

	case *parser.DerivNode:
		return lv.walkDerivNode(nt)

	default:
		panic("unknown node type")
	}
}

func (lv *LatexVisualiser) BeginEq() {
	lv.text.WriteString("\n\\begin{dmath}\n")
}

func (lv *LatexVisualiser) EndEq() {
	lv.text.WriteString("\n\\end{dmath}\n")
}

func (lv *LatexVisualiser) GenEqu() {
	lv.text.WriteString(" = ")
}

func (lv *LatexVisualiser) GenStr(str string) {
	lv.text.WriteString(str)
}

func (lv *LatexVisualiser) GenTexForNode(root parser.ASTNode) {

	lv.text.WriteString(lv.walkNode(root))
}
