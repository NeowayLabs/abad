// Package token exports the grammar's lexical tokens.
package token

import "fmt"

type (
	// Type of tokens
	Type int
)

const (
	Illegal Type = iota
	Bool
	Decimal
	Hexadecimal
	Octal
	String

	Minus
	Plus

	Dot
	LParen
	RParen
	Comma
	SemiColon

	Ident

	Null
	Undefined

	Newline

	EOF
)

var names = map[Type]string{
	Illegal:     "Illegal",
	Decimal:     "Decimal",
	Hexadecimal: "Hexadecimal",
	Octal:       "Octal",
	String:      "String",
	Bool:        "Bool",
	Minus:       "-",
	Plus:        "+",
	Dot:         ".",
	LParen:      "(",
	RParen:      ")",
	Comma:       ",",
	Ident:       "Ident",
	SemiColon:   "SemiColon",
	Null:        "Null",
	Undefined:   "Undefined",
	Newline:     "LineTerminator",
	EOF:         "EOF",
}

func (t Type) String() string {
	str, ok := names[t]
	if !ok {
		panic(fmt.Sprintf("unknown token type[%d]", t))
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
