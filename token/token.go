// Package token exports the grammar's lexical tokens.
package token

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

	EOF
)

var names = map[Type]string{
	Illegal:     "Illegal",
	Decimal:     "Decimal",
	Hexadecimal: "Hexadecimal",
	Octal:       "Octal",
	String:      "String",
	EOF:         "EOF",
}

func (t Type) String() string {
	return names[t]
}