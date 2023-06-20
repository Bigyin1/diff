package main

type LexemeMeta struct {
	Lexeme     string
	TokenClass string
	TokenName  string
}

var ConstantLexemes = []LexemeMeta{

	{
		Lexeme:     "sin",
		TokenClass: "ClassFunction",
		TokenName:  "Sin",
	},
	{
		Lexeme:     "cos",
		TokenClass: "ClassFunction",
		TokenName:  "Cos",
	},
	{
		Lexeme:     "tg",
		TokenClass: "ClassFunction",
		TokenName:  "Tg",
	},
	{
		Lexeme:     "ctg",
		TokenClass: "ClassFunction",
		TokenName:  "Ctg",
	},
	{
		Lexeme:     "ln",
		TokenClass: "ClassFunction",
		TokenName:  "Ln",
	},
	{
		Lexeme:     "+",
		TokenClass: "ClassOperator",
		TokenName:  "Plus",
	},
	{
		Lexeme:     "-",
		TokenClass: "ClassOperator",
		TokenName:  "Minus",
	},
	{
		Lexeme:     "*",
		TokenClass: "ClassOperator",
		TokenName:  "Mult",
	},
	{
		Lexeme:     "/",
		TokenClass: "ClassOperator",
		TokenName:  "Div",
	},
	{
		Lexeme:     "^",
		TokenClass: "ClassOperator",
		TokenName:  "Pow",
	},
	{
		Lexeme:     "e",
		TokenClass: "ClassConst",
		TokenName:  "Euler",
	},
	{
		Lexeme:     "pi",
		TokenClass: "ClassConst",
		TokenName:  "Pi",
	},
	{
		Lexeme:     "(",
		TokenClass: "ClassParen",
		TokenName:  "LParen",
	},
	{
		Lexeme:     ")",
		TokenClass: "ClassParen",
		TokenName:  "RParen",
	},
}

var VariableLexemes = []LexemeMeta{
	{
		Lexeme:     "",
		TokenClass: "ClassNumber",
		TokenName:  "Number",
	},
	{
		Lexeme:     "",
		TokenClass: "ClassVariable",
		TokenName:  "Variable",
	},
}
