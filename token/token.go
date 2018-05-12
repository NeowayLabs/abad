// Package token exports the grammar's lexical tokens.
package token

import "fmt"

type (
	// Type of tokens
	Type int
)

const (
	Illegal Type = iota
	Decimal
	Hexadecimal
	Octal
	String

	Minus
	Plus

	Dot
	LParen
	RParen

	Ident

	EOF
)

var names = map[Type]string{
	Illegal:     "Illegal",
	Decimal:     "Decimal",
	Hexadecimal: "Hexadecimal",
	Octal:       "Octal",
	String:      "String",
	Minus:       "-",
	Plus:        "+",
	Dot:         ".",
	LParen:      "(",
	RParen:      ")",
	Ident:       "Ident",
	EOF:         "EOF",
}

func (t Type) String() string {
	str, ok := names[t]
	if !ok {
		panic(fmt.Sprintf("symbol %d not found", t))
	}
	return str
}

func IsNumber(t Type) bool {
	return t == Decimal ||
		t == Hexadecimal ||
		t == Octal
}

func IsUnaryOperator(t Type) bool {
	return t == Minus ||
		t == Plus
}