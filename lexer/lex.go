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
	
		currentState := lexerInitialState(code)
		
		for currentState != nil {
			token, newState := currentState()
			if token != nil {
				tokens <- *token
			}
			currentState = newState
		}
		
		close(tokens)
	}()

	return tokens
}

type lexerState func() (*Tokval, lexerState)

func lexerInitialState(code utf16.Str) lexerState {
	return func() (*Tokval, lexerState) {
		if isNumber(code[0]) {
			return nil, decimalOrHexadecimalState(code, 1)
		}
		// TODO: Almost everything =)
		return nil, nil
	}
}

func decimalOrHexadecimalState(code utf16.Str, position uint) lexerState {
	return func() (*Tokval, lexerState) {
		if isEOF(code, position) {
			return &Tokval{
				Type: token.Decimal,
				Value: code,
			}, nil
		}
		
		if isNumber(code[position]) {
			return nil, decimalState(code, position + 1)
		}
		// TODO: need more tests to validate input 0x and 0X
		return nil, hexadecimalState(code, position) 
	}
}

func hexadecimalState(code utf16.Str, position uint) lexerState {
	// TODO: need more tests to validate x/X before continuing
	// TODO: tests validating invalid hexadecimals
	return func() (*Tokval, lexerState) {
		if isEOF(code, position) {
			return &Tokval{
				Type: token.Hexadecimal,
				Value: code,
			}, nil
		}
		
		return nil, hexadecimalState(code, position + 1)
	}
}

func decimalState(code utf16.Str, position uint) lexerState {
	// TODO: tests validating invalid decimals
	return func() (*Tokval, lexerState) {
		if isEOF(code, position) {
			return &Tokval{
				Type: token.Decimal,
				Value: code,
			}, nil
		}
		
		return nil, decimalState(code, position + 1)
	}
}

func isNumber(utf16char uint16) bool {
	// TODO
	return true
}

func isEOF(code utf16.Str, position uint) bool {
	return position >= uint(len(code))
}

func (t Tokval) Equal(other Tokval) bool {
	return t.Type == other.Type && t.Value.Equal(other.Value)
}