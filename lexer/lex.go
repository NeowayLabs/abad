package lexer

import (
	"unicode"
	
	"github.com/NeowayLabs/abad/internal/utf16"
	"github.com/NeowayLabs/abad/token"
)

type Tokval struct {
	Type  token.Type
	Value utf16.Str
	Line uint
	Column uint
}

var EOF Tokval = Tokval{ Type: token.EOF }

func (t Tokval) Equal(other Tokval) bool {
	return t.Type == other.Type && t.Value.Equal(other.Value)
}

func (t Tokval) EqualPos(other Tokval) bool {
	return t.Line == other.Line && t.Column == other.Column
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
		currentState := newLexer(decodedCode).initialState
		
		for currentState != nil {
			token, newState := currentState()
			tokens <- token
			currentState = newState
		}
		
		close(tokens)
	}()

	return tokens
}

type lexer struct {
	code []rune
	position uint
	line uint
	column uint
}

type lexerState func() (Tokval, lexerState)

func newLexer(code []rune) *lexer {
	return &lexer {code:code, line: 1, column: 1}
}

func (l *lexer) initialState() (Tokval, lexerState) {

	if l.isEOF() {
		return EOF, nil
	}
		
	if l.isInvalidRune() {
		return l.illegalToken()
	}
	
	if l.isPlusSign() {
		// TODO: handle ++
		return l.token(token.Plus), l.initialState
	}
	
	if l.isMinusSign() {
		// TODO: handle --
		return l.token(token.Minus), l.initialState
	}
		
	if l.isNumber() {
		l.fwd()
		return l.numberState()
	}
		
	if l.isDot() {
		l.fwd()
		return l.realDecimalState()
	}
	
	if l.isRightParen() {
		return l.token(token.RParen), l.initialState
	}
		
	return l.identifierState()
}

func (l *lexer) numberState() (Tokval, lexerState) {

	if l.isEOF() {
		return l.token(token.Decimal), l.initialState
	}
	
	if l.isDot() {
		l.fwd()
		return l.realDecimalState()
	}
	
	if l.isNumber() {
		l.fwd()
		return l.decimalState(true)
	}
	
	if l.isHexStart() {
		l.fwd()
		
		if l.isEOF() {
			return l.illegalToken()
		}
		
		return l.hexadecimalState()
	}
	
	if l.isExponentPartStart() {
		l.fwd()
		return l.exponentPartState()
	}
		
	return l.illegalToken()
}

func (l *lexer) illegalToken() (Tokval, lexerState) {
	return Tokval{
		Type: token.Illegal,
		Value: newStr(l.code),
	}, nil
}

func (l *lexer) identifierState() (Tokval, lexerState) {

	for !l.isEOF() {
		if l.isDot() {
			l.bwd()
			return l.token(token.Ident), l.accessMemberState
		}
		// TODO: test and code right parenthesis on this state
		if l.isLeftParen() {
			l.bwd()
			return l.token(token.Ident), l.leftParenState
		}
		l.fwd()
	}
		
	return l.token(token.Ident), l.initialState
}

func (l *lexer) leftParenState() (Tokval, lexerState) {
	return l.token(token.LParen), l.initialState
}

func (l *lexer) startIdentifierState() (Tokval, lexerState) {
	if l.isNumber() {
		return l.illegalToken()
	}
	if l.isDot() {
		return l.illegalToken()
	}
	return l.identifierState()
}

func (l *lexer) accessMemberState() (Tokval, lexerState) {
	return l.token(token.Dot), l.startIdentifierState
}

func (l *lexer) hexadecimalState() (Tokval, lexerState) {

	for !l.isEOF() {
		if !l.isHexadecimal() {
			return l.illegalToken()
		}
		l.fwd()
	}
		
	return l.token(token.Hexadecimal), l.initialState
}

func (l *lexer) realDecimalState() (Tokval, lexerState) {

	for !l.isEOF() {
		if l.isExponentPartStart() {
			l.fwd()
			return l.exponentPartState()
		}
		
		if !l.isNumber() {
			return l.illegalToken()
		}
		
		l.fwd()
	}
	
	return l.token(token.Decimal), l.initialState
}

func (l *lexer) decimalState(allowExponent bool) (Tokval, lexerState) {

	for !l.isEOF() {
		if l.isExponentPartStart() {
			if !allowExponent {
				return l.illegalToken()
			}
			l.fwd()
			return l.exponentPartState()
		}
		
		if l.isDot() {
			l.fwd()
			return l.realDecimalState()
		}
		
		if !l.isNumber() {
			return l.illegalToken()
		}
		
		l.fwd()
	}
	
	return l.token(token.Decimal), l.initialState
}

func (l *lexer) exponentPartState() (Tokval, lexerState) {
	// TODO: can exponent be like: 1.0e ?
	
	if l.isMinusSign() || l.isPlusSign() {
		// TODO: test 1.0e- and 1.0e+
		l.fwd()
	}
	
	return l.decimalState(false)
}

func (l *lexer) cur() rune {
	return l.code[l.position]
}

func (l *lexer) isNumber() bool {
	return containsRune(numbers, l.cur())
}

func (l *lexer) isEOF() bool {
	return l.position >= uint(len(l.code))
}

func (l *lexer) isDot() bool {
	return l.cur() == dot
}

func (l *lexer) isHexStart() bool {
	return containsRune(hexStart, l.cur())
}

func (l *lexer) isInvalidRune() bool {
	return unicode.ReplacementChar == l.cur()
}

func (l *lexer) isMinusSign() bool {
	return l.cur() == minusSign
}

func (l *lexer) isPlusSign() bool {
	return l.cur() == plusSign
}

func (l *lexer) isLeftParen() bool {
	return l.cur() == leftParen
}

func (l *lexer) isRightParen() bool {
	return l.cur() == rightParen
}

func (l *lexer) isHexadecimal() bool {
	return containsRune(hexnumbers, l.cur())
}

func (l *lexer) isExponentPartStart() bool {
	return containsRune(exponentPartStart, l.cur())
}

func (l *lexer) fwd() {
	l.position += 1
}

func (l *lexer) bwd() {
	l.position -= 1
}

// token will generate a token consuming all the code
// until the current position. After calling this method
// the token will not be available anymore (it has been consumed)
// on the code and the position will be reset to zero.
func (l *lexer) token(t token.Type) Tokval {
	var val []rune
	
	if l.isEOF() {
		val = l.code
		l.code = nil
	} else {
		val = l.code[:l.position + 1]
		l.code = l.code[l.position + 1:]
	}
	
	// FIXME: should use position, but for now this works for the lack of tests
	c := l.column
	l.column += 1
	
	l.position = 0	
	return Tokval{Type:t, Value: newStr(val), Line: l.line, Column: c}
}

var numbers []rune
var hexnumbers []rune
var dot rune
var minusSign rune
var plusSign rune
var leftParen rune
var rightParen rune
var hexStart []rune
var exponentPartStart []rune

func init() {
	numbers = []rune("0123456789")
	hexnumbers = append(numbers, []rune("abcdefABCDEF")...)
	dot = rune('.')
	minusSign = rune('-')
	plusSign = rune('+')
	leftParen = rune('(')
	rightParen = rune(')')
	hexStart = []rune("xX")
	exponentPartStart = []rune("eE")
}

func containsRune(runes []rune, r rune) bool {
	for _, n := range runes {
		if r == n {
			return true
		}
	}
	return false	
}

func newStr(r []rune) utf16.Str {
	return utf16.NewFromRunes(r)
}