// Code generated DO NOT EDIT

package lexer

import "fmt"

const (
	ClassNone TokenClass = iota

	ClassOperator
	ClassConst
	ClassParen
	ClassNumber
	ClassVariable
	ClassFunction
)

const (
	EOF TokenName = iota

	Sin
	Cos
	Tg
	Ctg
	Ln
	Plus
	Minus
	Mult
	Div
	Pow
	Euler
	Pi
	LParen
	RParen
	Number
	Variable
)

var reservedWordsToTokenMeta = map[string]TokenMeta{
	"sin": {ClassFunction, Sin},
	"cos": {ClassFunction, Cos},
	"tg":  {ClassFunction, Tg},
	"ctg": {ClassFunction, Ctg},
	"ln":  {ClassFunction, Ln},
	"+":   {ClassOperator, Plus},
	"-":   {ClassOperator, Minus},
	"*":   {ClassOperator, Mult},
	"/":   {ClassOperator, Div},
	"^":   {ClassOperator, Pow},
	"e":   {ClassConst, Euler},
	"pi":  {ClassConst, Pi},
	"(":   {ClassParen, LParen},
	")":   {ClassParen, RParen},
}

var tokenNamesToTokenMeta = map[TokenName]TokenMeta{
	Sin:    {ClassFunction, Sin},
	Cos:    {ClassFunction, Cos},
	Tg:     {ClassFunction, Tg},
	Ctg:    {ClassFunction, Ctg},
	Ln:     {ClassFunction, Ln},
	Plus:   {ClassOperator, Plus},
	Minus:  {ClassOperator, Minus},
	Mult:   {ClassOperator, Mult},
	Div:    {ClassOperator, Div},
	Pow:    {ClassOperator, Pow},
	Euler:  {ClassConst, Euler},
	Pi:     {ClassConst, Pi},
	LParen: {ClassParen, LParen},
	RParen: {ClassParen, RParen},

	Number:   {ClassNumber, Number},
	Variable: {ClassVariable, Variable},
}

func (t *Token) String() string {

	if Sin == t.Name {
		return "Sin"
	}

	if Cos == t.Name {
		return "Cos"
	}

	if Tg == t.Name {
		return "Tg"
	}

	if Ctg == t.Name {
		return "Ctg"
	}

	if Ln == t.Name {
		return "Ln"
	}

	if Plus == t.Name {
		return "Plus"
	}

	if Minus == t.Name {
		return "Minus"
	}

	if Mult == t.Name {
		return "Mult"
	}

	if Div == t.Name {
		return "Div"
	}

	if Pow == t.Name {
		return "Pow"
	}

	if Euler == t.Name {
		return "Euler"
	}

	if Pi == t.Name {
		return "Pi"
	}

	if LParen == t.Name {
		return "LParen"
	}

	if RParen == t.Name {
		return "RParen"
	}

	if Number == t.Name {
		return fmt.Sprintf("%s (%s)", "Number", t.Value)
	}

	if Variable == t.Name {
		return fmt.Sprintf("%s (%s)", "Variable", t.Value)
	}

	panic("bad")
}
