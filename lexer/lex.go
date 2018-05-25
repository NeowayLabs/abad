package lexer

import (
	"unicode"
	
	"github.com/NeowayLabs/abad/internal/utf16"
	"github.com/NeowayLabs/abad/token"
)

type Tokval struct {
	Type  token.Type
	Value utf16.Str
}

var EOF Tokval = Tokval{ Type: token.EOF }

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
	
		decodedCode := code.Runes()
		currentState := initialState(decodedCode)
		
		for currentState != nil {
			token, newState := currentState()
			tokens <- token
			currentState = newState
		}
		
		close(tokens)
	}()

	return tokens
}

type lexerState func() (Tokval, lexerState)

func initialState(code []rune) lexerState {

	return func() (Tokval, lexerState) {
		if len(code) == 0 {
			return EOF, nil
		}
		
		if isInvalidRune(code[0]) {
			return illegalToken(code)
		}
		
		if isNumber(code[0]) {
			return numberState(code, 1)
		}
		
		if isDot(code[0]) {
			return realDecimalState(code, 1)
		}
		
		return EOF, nil
	}
}

func numberState(code []rune, position uint) (Tokval, lexerState) {

	if isEOF(code, position) {
		return Tokval{
			Type: token.Decimal,
			Value: newStr(code),
		}, initialState(code[position:])
	}
	
	if isDot(code[position]) {
		return realDecimalState(code, position + 1)
	}
	
	if isNumber(code[position]) {
		return decimalState(code, position + 1)
	}
	
	if isHexStart(code[position]) {
		if isEOF(code, position + 1) {
			return illegalToken(code)
		}
		return hexadecimalState(code, position + 1)
	}
	
	if isExponentPartStart(code[position]) {
		return exponentPartState(code, position + 1)
	}
		
	return illegalToken(code)
}

func illegalToken(code []rune) (Tokval, lexerState) {
	return Tokval{
		Type: token.Illegal,
		Value: newStr(code),
	}, nil
}

func hexadecimalState(code []rune, position uint) (Tokval, lexerState) {

	for !isEOF(code, position) {
		if !isHexadecimal(code[position]) {
			return illegalToken(code)
		}
		position += 1
	}
		
	return Tokval{
		Type: token.Hexadecimal,
		Value: newStr(code),
	}, initialState(code[position:])
}

func realDecimalState(code []rune, position uint) (Tokval, lexerState) {

	for !isEOF(code, position) {
		if isExponentPartStart(code[position]) {
			return exponentPartState(code, position + 1)
		}
		
		if !isNumber(code[position]) {
			return illegalToken(code)
		}
		
		position += 1
	}
	
	return Tokval{
		Type: token.Decimal,
		Value: newStr(code),
	}, initialState(code[position:])
}

func decimalState(code []rune, position uint) (Tokval, lexerState) {

	for !isEOF(code, position) {
		if isExponentPartStart(code[position]) {
			return exponentPartState(code, position + 1)
		}
		
		if isDot(code[position]) {
			return realDecimalState(code, position + 1)
		}
		
		if !isNumber(code[position]) {
			return illegalToken(code)
		}
		position += 1
	}
	
	return Tokval{
		Type: token.Decimal,
		Value: newStr(code),
	}, initialState(code[position:])
}

func exponentPartState(code []rune, position uint) (Tokval, lexerState) {
	// TODO: Need to validate malformed exponent numbers to improve this
	// eg: 1.0e123e
	// TODO: can exponent be like: 1.0e ?
	
	if isMinusSign(code[position]) || isPlusSign(code[position]) {
		// TODO: test 1.0e- and 1.0e+
		position += 1
	}
	
	return decimalState(code, position)
}

func isNumber(r rune) bool {
	return containsRune(numbers, r)
}

func containsRune(runes []rune, r rune) bool {
	for _, n := range runes {
		if r == n {
			return true
		}
	}
	return false	
}

func isEOF(code []rune, position uint) bool {
	return position >= uint(len(code))
}

func isDot(r rune) bool {
	return r == dot
}

func isHexStart(r rune) bool {
	return containsRune(hexStart, r)
}

func newStr(r []rune) utf16.Str {
	return utf16.NewFromRunes(r)
}

func isInvalidRune(r rune) bool {
	return unicode.ReplacementChar == r
}

func isMinusSign(r rune) bool {
	return r == minusSign
}

func isPlusSign(r rune) bool {
	return r == plusSign
}

func isHexadecimal(r rune) bool {
	return containsRune(hexnumbers, r)
}

func isExponentPartStart(r rune) bool {
	return containsRune(exponentPartStart, r)
}

var numbers []rune
var hexnumbers []rune
var dot rune
var minusSign rune
var plusSign rune
var hexStart []rune
var exponentPartStart []rune

func init() {
	numbers = []rune("0123456789")
	hexnumbers = append(numbers, []rune("abcdefABCDEF")...)
	dot = rune('.')
	minusSign = rune('-')
	plusSign = rune('+')
	hexStart = []rune("xX")
	exponentPartStart = []rune("eE")
}