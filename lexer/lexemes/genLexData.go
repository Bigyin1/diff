package main

import (
	"bytes"
	"flag"
	"go/format"
	"os"
	"path"
	"text/template"
)

type lexGenCtx struct {
	Package      string
	ConstLexemes []LexemeMeta
	VarLexemes   []LexemeMeta

	TokClasses []string
	TokNames   []string
}

func collectTokClasses(ctx *lexGenCtx) {
	classSet := make(map[string]struct{})
	for _, v := range ctx.ConstLexemes {
		classSet[v.TokenClass] = struct{}{}
	}
	for _, v := range ctx.VarLexemes {
		classSet[v.TokenClass] = struct{}{}
	}

	for k := range classSet {
		ctx.TokClasses = append(ctx.TokClasses, k)
	}
}

func collectTokNames(ctx *lexGenCtx) {
	for _, v := range ctx.ConstLexemes {
		ctx.TokNames = append(ctx.TokNames, v.TokenName)
	}
	for _, v := range ctx.VarLexemes {
		ctx.TokNames = append(ctx.TokNames, v.TokenName)
	}
}

func main() {

	tmplFile := flag.String("tmpl", "", "")
	outFile := flag.String("o", "", "")
	flag.Parse()

	tmpl, err := template.New(path.Base(*tmplFile)).ParseFiles(*tmplFile)
	if err != nil {
		panic(err)
	}

	ctx := lexGenCtx{
		Package:      os.Getenv("GOPACKAGE"),
		ConstLexemes: ConstantLexemes,
		VarLexemes:   VariableLexemes,
	}

	collectTokClasses(&ctx)
	collectTokNames(&ctx)

	out, err := os.Create(*outFile)
	if err != nil {
		panic(err)
	}
	defer out.Close()

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, &ctx)
	if err != nil {
		panic(err)
	}

	p, err := format.Source(buf.Bytes())
	if err != nil {
		panic(err)
	}

	_, err = out.Write(p)
	if err != nil {
		panic(err)
	}
}
