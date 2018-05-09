// Package token exports the grammar's lexical tokens.
package token

type (
	// Type of tokens
	Type int
)

const (
	Unknown Type = iota
)

var names = map[Type]string{
	Unknown: "Unknown",
}

func (t Type) String() string {
	return names[t]
}