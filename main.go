package main

import (
	"diff/derivator"
	"diff/lexer"
	"diff/parser"
	"diff/visualisers/graphviz"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
)

func main() {

	outFile := flag.String("f", "out.tex", "latex output file")
	derVar := flag.String("var", "x", "variable to derivate by")
	gv := flag.Bool("gv", false, "generate graphviz derivative graph")
	flag.Parse()

	log.SetOutput(io.Discard)
	input, err := io.ReadAll(os.Stdin)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	l := lexer.NewLexer(string(input))
	toks, err := l.Run()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	p := parser.NewParser(toks)
	root, nm, err := p.Run()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	d, err := derivator.NewDerivator(root, nm, *derVar, *outFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	gr, err := d.Run()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if *gv {
		gvData := graphviz.GenGraphViz(gr)
		file, err := os.Create("derivGraph.dot")
		defer file.Close()

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		_, err = file.Write(gvData)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
}
