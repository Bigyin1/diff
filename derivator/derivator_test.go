package derivator

import (
	"diff/lexer"
	"diff/parser"
	"diff/visualisers/graphviz"
	"fmt"
	"os"
	"testing"
)

// TODO: tests )))

func TestDerivator(t *testing.T) {
	input := "sin(1 / (x*cos(x^(3/pi) * ln(x^2)))) "
	l := lexer.NewLexer(input)

	toks, err := l.Run()
	if err != nil {
		t.Error(err)
		return
	}

	p := parser.NewParser(toks)
	root, m, err := p.Run()
	if err != nil {
		t.Error(err)
		return
	}

	gv := graphviz.GenGraphViz(root)
	file, _ := os.Create("iniExpr.dot")
	_, err = file.Write(gv)
	if err != nil {
		t.Error(err)
		return
	}
	file.Close()

	d, _ := NewDerivator(root, m, "x", "out.tex")
	nr, _ := d.Run()

	gv = graphviz.GenGraphViz(nr)
	file, _ = os.Create("derivExpr.dot")
	_, err = file.Write(gv)
	if err != nil {
		t.Error(err)
		return
	}
	file.Close()

	fmt.Println(nr.String())

}
