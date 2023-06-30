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
	file *strings.Builder
}

func NewLatexVisualiser() *LatexVisualiser {

	var file strings.Builder

	file.WriteString(header)

	return &LatexVisualiser{
		file: &file,
	}
}

func (lv *LatexVisualiser) BeginDoc() {
	lv.file.WriteString("\n\\begin{document}\n")
}

func (lv *LatexVisualiser) EndDoc() {
	lv.file.WriteString("\n\\end{document}\n")

	f, _ := os.Create("out.tex")
	defer f.Close()

	_, _ = f.WriteString(lv.file.String())
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

	default:
		panic("unknown node type")
	}
}

func (lv *LatexVisualiser) GenTexForNode(root parser.ASTNode) {

	lv.file.WriteString("\n\\begin{math}\n")

	lv.file.WriteString(lv.walkNode(root))

	lv.file.WriteString("\n\\end{math}\n")
}
