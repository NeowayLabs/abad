package lexer

import (
	"github.com/NeowayLabs/abad/internal/utf16"
	"github.com/NeowayLabs/abad/token"
)

type (
	Tokval struct {
		Type  token.Type
		Value utf16.Str
	}
)

func Lex(code utf16.Str) <-chan Tokval {
	tokens := make(chan Tokval)
	
	go func() {
		tokens <- Tokval{
			Type: token.Decimal,
			Value: code,
		}
		close(tokens)
	}()

	return tokens
}

func (t Tokval) Equal(other Tokval) bool {
	return t.Type == other.Type && t.Value.Equal(other.Value)
}