package parser

import (
	"diff/lexer"
	"fmt"
)

type UnexpTokenError struct {
	tok *lexer.Token
	exp lexer.TokenName
}

func (e *UnexpTokenError) Error() string {

	expTok := lexer.Token{
		TokenMeta: lexer.TokenMeta{Class: lexer.ClassNone, Name: e.exp},
	}

	if e.exp != lexer.EOF {
		return fmt.Sprintf("unexpected token %s at %d, %d; wanted %s",
			e.tok.String(), e.tok.Row, e.tok.Column, expTok.String())
	}

	if e.tok.Name == lexer.EOF {
		return "unexpected end of expression"
	}

	return fmt.Sprintf("unexpected token %s at %d, %d",
		e.tok.String(), e.tok.Row, e.tok.Column)
}
