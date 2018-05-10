package lexer

import (
	"fmt"

	"github.com/NeowayLabs/abad/internal/utf16"
	"github.com/NeowayLabs/abad/token"
)

type (
	Tokval struct {
		Type  token.Type
		Value utf16.Str
	}
)

// TODO: remove me
var mock = map[string][]Tokval{
	// NumericLiteral
	// https://es5.github.io/#x7.8.3
	// https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Lexical_grammar#Decimal
	"0":                  onetok(token.Decimal, "0"),
	"1":                  onetok(token.Decimal, "1"),
	"1234":               onetok(token.Decimal, "1234"),
	"1234567890":         onetok(token.Decimal, "1234567890"),
	"1a":                 onetok(token.Illegal, "1a"),
	"0x0":                onetok(token.Hexadecimal, "0x0"),
	"0x1234567890abcdef": onetok(token.Hexadecimal, "0x1234567890abcdef"),
	"0xff":               onetok(token.Hexadecimal, "0xff"),
	".0":                 onetok(token.Decimal, ".0"),
	".1":                 onetok(token.Decimal, ".1"),
	".0000":              onetok(token.Decimal, ".0000"),
	"1234.":              onetok(token.Decimal, "1234."),
	"0.12345":            onetok(token.Decimal, "0.12345"),
	"0.a":                onetok(token.Illegal, "0.a"),
	"12.13.":             onetok(token.Illegal, "12.13."),
	"1.0e10":             onetok(token.Decimal, "1.0e10"),
	".1e10":              onetok(token.Decimal, ".1e10"),
	"1e10":               onetok(token.Decimal, "1e10"),
	"1e-10":              onetok(token.Decimal, "1e-10"),

	"-1e-10": toks(
		tok(token.Minus, "-"), tok(token.Decimal, "1e-10"),
	),
	"+0":      toks(tok(token.Plus, "+"), tok(token.Decimal, "0")),
	"-1":      toks(tok(token.Minus, "-"), tok(token.Decimal, "1")),
	"-1234":   toks(tok(token.Minus, "-"), tok(token.Decimal, "1234")),
	"+1234":   toks(tok(token.Plus, "+"), tok(token.Decimal, "1234")),
	"-0x0":    toks(tok(token.Minus, "-"), tok(token.Decimal, "0x0")),
	"-0xff":   toks(tok(token.Minus, "-"), tok(token.Decimal, "0xff")),
	"+0xff":   toks(tok(token.Plus, "+"), tok(token.Decimal, "0xff")),
	"-.0":     toks(tok(token.Minus, "-"), tok(token.Decimal, ".0")),
	"+.0":     toks(tok(token.Plus, "+"), tok(token.Decimal, ".0")),
	"-.0e1":   toks(tok(token.Minus, "-"), tok(token.Decimal, ".0e1")),
	"-12.13.": toks(tok(token.Minus, "-"), tok(token.Illegal, "12.13.")),
}

// TODO: remove me
func onetok(t token.Type, v string) []Tokval {
	return []Tokval{tok(t, v)}
}

func tok(t token.Type, v string) Tokval {
	return Tokval{
		Type:  t,
		Value: utf16.Encode(v),
	}
}

func toks(toks ...Tokval) []Tokval {
	return toks
}

func Lex(code utf16.Str) <-chan Tokval {
	tokens := make(chan Tokval)
	tokvals := mock[code.String()]
	if len(tokvals) == 0 {
		panic(fmt.Errorf("mock not implemented for: %s", code))
	}

	go func() {
		for _, tok := range tokvals {
			tokens <- tok
		}
		close(tokens)
	}()

	return tokens
}