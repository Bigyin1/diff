// Code generated DO NOT EDIT

package lexer

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

	EOF: {ClassNone, EOF},
}

func (t TokenName) String() string {

	if Sin == t {
		return "Sin"
	}

	if Cos == t {
		return "Cos"
	}

	if Tg == t {
		return "Tg"
	}

	if Ctg == t {
		return "Ctg"
	}

	if Ln == t {
		return "Ln"
	}

	if Plus == t {
		return "Plus"
	}

	if Minus == t {
		return "Minus"
	}

	if Mult == t {
		return "Mult"
	}

	if Div == t {
		return "Div"
	}

	if Pow == t {
		return "Pow"
	}

	if Euler == t {
		return "Euler"
	}

	if Pi == t {
		return "Pi"
	}

	if LParen == t {
		return "LParen"
	}

	if RParen == t {
		return "RParen"
	}

	if Number == t {
		return "Number"
	}

	if Variable == t {
		return "Variable"
	}

	if EOF == t {
		return "EOF"
	}

	panic("bad")
}
