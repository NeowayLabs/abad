package lexer

import (
	"github.com/NeowayLabs/abad/internal/utf16"
	"github.com/NeowayLabs/abad/token"
)

type Tokval struct {
	Type  token.Type
	Value utf16.Str
}

func (t Tokval) Equal(other Tokval) bool {
	return t.Type == other.Type && t.Value.Equal(other.Value)
}

// Lex will lex the given crappy JS code (utf16 yay) and provide a
// stream of tokens as a result (the returned channel).
//
// The caller should iterate on the given channel until it is
// closed indicating a EOF (or an error). Errors should be
// handled by checking the type of the token.
//
// A goroutine will be started to lex the given code, if you
// do not iterate the returned channel the goroutine will leak,
// you MUST drain the provided channel.
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
		// TODO: handle empty input
		
		if len(code) == 0 {
			return nil, nil
		}
		
		if isNumber(code[0]) {
			return nil, decimalOrHexadecimalState(code, 1)
		}
		
		if isDot(code[0]) {
			return nil, decimalState(code, 1)
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
		
		if isNumber(code[position]) || isDot(code[position]) {
			return nil, decimalState(code, position + 1)
		}
		
		if isHexStart(code[position]) {
			if isEOF(code, position + 1) {
				return illegalToken(code)
			}
			return nil, hexadecimalState(code, position)
		}
		
		return illegalToken(code)
	}
}

func illegalToken(code utf16.Str) (*Tokval, lexerState) {
	return &Tokval{
		Type: token.Illegal,
		Value: code,
	}, nil
}

func hexadecimalState(code utf16.Str, position uint) lexerState {
	// TODO: need more tests to validate x/X before continuing
	// TODO: tests validating invalid hexadecimals
	
	return func() (*Tokval, lexerState) {
	
		for !isEOF(code, position) {
			position += 1
		}
		
		return &Tokval{
			Type: token.Hexadecimal,
			Value: code,
		}, nil
	}
}

func decimalState(code utf16.Str, position uint) lexerState {
	// TODO: tests validating invalid decimals
	return func() (*Tokval, lexerState) {
	
		for !isEOF(code, position) {
			position += 1
		}
	
		return &Tokval{
			Type: token.Decimal,
			Value: code,
		}, nil
	}
}

func isNumber(utf16char uint16) bool {
	str := strFromChar(utf16char)
	return numbers.Contains(str)
}

func isEOF(code utf16.Str, position uint) bool {
	return position >= uint(len(code))
}

func isDot(utf16char uint16) bool {
	str := strFromChar(utf16char)
	return dot.Equal(str)
}

func isHexStart(utf16char uint16) bool {
	str := strFromChar(utf16char)
	return hexStart.Contains(str)
}

func strFromChar(utf16char uint16) utf16.Str {
	return utf16.Str([]uint16{utf16char})
}


var numbers utf16.Str
var dot utf16.Str
var exponents utf16.Str
var hexStart utf16.Str

func init() {
	numbers = utf16.NewStr("0123456789")
	dot = utf16.NewStr(".")
	exponents = utf16.NewStr("eE")
	hexStart = utf16.NewStr("xX")
}