package lexer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLexer(t *testing.T) {

	testCases := []struct {
		name string

		input string
		toks  []Token
		err   error
	}{
		{
			name:  "correct input",
			input: "(23.09 * 2)^pi +sin(x)  ",
			toks: []Token{
				{
					TokenMeta: tokenNamesToTokenMeta[LParen],
					Value:     "(",
					Row:       1,
					Column:    1,
				},
				{
					TokenMeta: tokenNamesToTokenMeta[Number],
					Value:     "23.09",
					Row:       1,
					Column:    2,
				},
				{
					TokenMeta: tokenNamesToTokenMeta[Mult],
					Value:     "*",
					Row:       1,
					Column:    8,
				},
				{
					TokenMeta: tokenNamesToTokenMeta[Number],
					Value:     "2",
					Row:       1,
					Column:    10,
				},
				{
					TokenMeta: tokenNamesToTokenMeta[RParen],
					Value:     ")",
					Row:       1,
					Column:    11,
				},
				{
					TokenMeta: tokenNamesToTokenMeta[Pow],
					Value:     "^",
					Row:       1,
					Column:    12,
				},
				{
					TokenMeta: tokenNamesToTokenMeta[Pi],
					Value:     "pi",
					Row:       1,
					Column:    13,
				},
				{
					TokenMeta: tokenNamesToTokenMeta[Plus],
					Value:     "+",
					Row:       1,
					Column:    16,
				},
				{
					TokenMeta: tokenNamesToTokenMeta[Sin],
					Value:     "sin",
					Row:       1,
					Column:    17,
				},
				{
					TokenMeta: tokenNamesToTokenMeta[LParen],
					Value:     "(",
					Row:       1,
					Column:    20,
				},
				{
					TokenMeta: tokenNamesToTokenMeta[Variable],
					Value:     "x",
					Row:       1,
					Column:    21,
				},
				{
					TokenMeta: tokenNamesToTokenMeta[RParen],
					Value:     ")",
					Row:       1,
					Column:    22,
				},
				{
					TokenMeta: TokenMeta{ClassNone, EOF},
				},
			},
			err: nil,
		},
		{
			name:  "incorrect sinus",
			input: "tg / sinx",
			toks: []Token{
				{
					TokenMeta: tokenNamesToTokenMeta[Tg],
					Value:     "tg",
					Row:       1,
					Column:    1,
				},
				{
					TokenMeta: tokenNamesToTokenMeta[Div],
					Value:     "/",
					Row:       1,
					Column:    4,
				},
				{
					TokenMeta: TokenMeta{ClassNone, EOF},
				},
			},
			err: &LexError{
				Row:    1,
				Column: 6,
				Str:    "sinx",
			},
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			l := NewLexer(test.input)

			toks, err := l.Run()

			assert.Equal(t, test.err, err)
			assert.Equal(t, test.toks, toks)
		})
	}
}
