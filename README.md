# Differentiator

Takes derivative of expression by chosen variable. Generates Latex file with step-by-step derivation process.

## Usage

Compile with: `go build -o diff main.go`

Execute `./diff` (add `--help` to check all options) and send expression to it's standard input.


## Internal Features

Differentiator uses Directed Acyclic Graph representation for expressions.

[Example graph](./graphs/iniExpr.dot.png) for initial expression `x^(x^y * tg(cos(ln(x^x)/x^2))) / (ln(x^x) * x^y)`
and [for it's derivative](./graphs/derivExpr.dot.png).
