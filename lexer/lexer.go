package lexer

import (
	"strings"
	"unicode/utf8"
)

type TokenClass uint16

type TokenName uint16

type TokenMeta struct {
	Class TokenClass
	Name  TokenName
}

type Token struct {
	TokenMeta

	Value string

	Row    int
	Column int
}

//go:generate go run lexemes/genLex.go lexemes/lexemes.go -tmpl lexemes/lex.tmpl -o generatedLexemes.go

type Lexer struct {
	input string

	start int // start position for the current lexeme
	pos   int // current position (>= start)
	width int // width in bytes of last rune

	tokens []Token
	err    *LexError

	currRow    int
	currColumn int
}

func (l *Lexer) next() (r rune) {

	if l.pos >= len(l.input) {
		l.width = 0
		return 0
	}

	r, l.width = utf8.DecodeRuneInString(l.input[l.pos:])
	l.pos += l.width
	return r
}

func (l *Lexer) currWord() string {
	return l.input[l.start:l.pos]
}

func (l *Lexer) ignore() {
	l.currColumn += utf8.RuneCountInString(l.input[l.start:l.pos])
	l.start = l.pos
}

func (l *Lexer) backup() {
	l.pos -= l.width
}

func (l *Lexer) peek() rune {
	r := l.next()
	l.backup()

	return r
}

type acceptFn func(r rune) bool

func (l *Lexer) accept(fn acceptFn) bool {
	if fn(l.next()) {
		return true
	}
	l.backup()
	return false
}

func (l *Lexer) acceptRun(fn acceptFn) {
	for fn(l.next()) {
	}
	l.backup()
}

func (l *Lexer) acceptSet(valid string) bool {
	if strings.ContainsRune(valid, l.next()) {
		return true
	}
	l.backup()
	return false
}

func (l *Lexer) emit(tMeta TokenMeta) {

	l.tokens = append(l.tokens,
		Token{
			TokenMeta: tMeta,
			Value:     l.input[l.start:l.pos],
			Column:    l.currColumn,
			Row:       l.currRow,
		})

	l.ignore()
}

func NewLexer(input string) *Lexer {

	return &Lexer{
		input:      input,
		tokens:     []Token{},
		currRow:    1,
		currColumn: 1,
	}
}

func (l *Lexer) Run() ([]Token, error) {

	for state := lexFunction; state != nil; {
		state = state(l)
	}

	return l.tokens, l.err
}
