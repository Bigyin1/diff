package parser

import (
	"diff/lexer"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TODO: tests )))
func TestParser(t *testing.T) {
	input := "sin(ln(x)^ (x*cos(x^(3/pi) * e)))"
	l := lexer.NewLexer(input)

	toks, err := l.Run()
	if err != nil {
		t.Error(err)
		return
	}

	p := NewParser(toks)
	root, _, err := p.Run()
	if err != nil {
		t.Error(err)
		return
	}

	fmt.Print(root.String())
}

func TestParserErrors(t *testing.T) {
	input := "x^"
	l := lexer.NewLexer(input)

	toks, err := l.Run()
	if err != nil {
		t.Error(err)
		return
	}

	p := NewParser(toks)
	root, _, err := p.Run()

	assert.Error(t, err)
	fmt.Print(root.String())
}
