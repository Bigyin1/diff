package derivator

import (
	"diff/lexer"
	"diff/parser"
	"diff/visualisers/graphviz"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TODO: tests )))

func TestDerivator(t *testing.T) {
	input := "x^(x^y * tg(x*cos(1/x^2))) / (5 *ln(x) * y) "
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

	d, err := NewDerivator(root, m, "x", "out.tex")
	assert.NoError(t, err)

	nr, err := d.Run()
	assert.NoError(t, err)

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
