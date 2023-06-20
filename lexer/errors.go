package lexer

import "fmt"

type LexError struct {
	Row    int
	Column int
	Str    string
}

func (e *LexError) Error() string {
	return fmt.Sprintf("unknown lexeme at line %d column %d: %s", e.Row, e.Column, e.Str)
}

func (l *Lexer) fail() {
	l.err = &LexError{
		Row:    l.currRow,
		Column: l.currColumn,
		Str:    l.input[l.start:l.pos],
	}
}
