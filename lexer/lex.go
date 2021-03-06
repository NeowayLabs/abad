package lexer

import (
	"fmt"
	"unicode"

	"github.com/NeowayLabs/abad/internal/utf16"
	"github.com/NeowayLabs/abad/token"
)

type Tokval struct {
	Type   token.Type
	Value  utf16.Str
	Line   uint
	Column uint
}

// EOF is the End of File token.
var EOF = Tokval{Type: token.EOF, Value: utf16.S("EOF")}

// Equal tells if token is the same as other.
func (t Tokval) Equal(other Tokval) bool {
	return t.Type == other.Type && t.Value.Equal(other.Value)
}

func (t Tokval) EqualPos(other Tokval) bool {
	return t.Line == other.Line && t.Column == other.Column
}

func (t Tokval) String() string {
	return fmt.Sprintf(
		"token:type[%s],value[%s],line[%d],column[%d]", t.Type, t.Value, t.Line, t.Column)
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
	code     []rune
	position uint
	line     uint
	column   uint

	puncStates map[rune]lexerState
}

type match struct {
	str   string
	token token.Type
}

type lexerState func() (Tokval, lexerState)

func newLexer(code []rune) *lexer {
	l := &lexer{code: code, line: 1, column: 1}
	l.initPuncStates()
	return l
}

func (l *lexer) initialState() (Tokval, lexerState) {

	l.skipSpaces()

	if l.isEOF() {
		return EOF, nil
	}

	if l.isInvalidRune() {
		return l.illegalToken()
	}

	if l.isNumber() {
		l.fwd()
		return l.numberState()
	}

	if l.isDoubleQuote() {
		l.fwd()
		return l.stringState()
	}

	if l.isPunctuator() {
		return l.punctuator()
	}

	return l.identifierState()
}

func (l *lexer) initPuncStates() {
	// http://es5.github.io/#x7.7

	state := func(t token.Type) lexerState {
		return func() (Tokval, lexerState) {
			return l.token(t), l.initialState
		}
	}

	l.puncStates = map[rune]lexerState{
		dot:        l.dotState,
		comma:      state(token.Comma),
		semiColon:  state(token.SemiColon),
		leftParen:  state(token.LParen),
		rightParen: state(token.RParen),
		rune('~'):  state(token.Not),
		rune('?'):  state(token.Ternary),
		rune(':'):  state(token.Colon),
		rune('['):  state(token.LBrack),
		rune(']'):  state(token.RBrack),
		rune('{'):  state(token.LBrace),
		rune('}'):  state(token.RBrace),
		rune('*'): l.acceptFirst([]match{
			{str: "*=", token: token.MulAssign},
			{str: "*", token: token.Mul},
		}),
		rune('/'): l.acceptFirst([]match{
			{str: "/=", token: token.QuoAssign},
			{str: "/", token: token.Quo},
		}),
		rune('%'): l.acceptFirst([]match{
			{str: "%=", token: token.RemAssign},
			{str: "%", token: token.Rem},
		}),
		rune('<'): l.acceptFirst([]match{
			{str: "<<=", token: token.LShiftAssign},
			{str: "<<", token: token.LShift},
			{str: "<=", token: token.LessEq},
			{str: "<", token: token.Less},
		}),
		rune('>'): l.acceptFirst([]match{
			{str: ">>>=", token: token.RShiftZeroAssign},
			{str: ">>>", token: token.RShiftZero},
			{str: ">>=", token: token.RShiftAssign},
			{str: ">>", token: token.RShift},
			{str: ">=", token: token.GreaterEq},
			{str: ">", token: token.Greater},
		}),
		rune('&'): l.acceptFirst([]match{
			{str: "&&", token: token.LAnd},
			{str: "&=", token: token.AndAssign},
			{str: "&", token: token.And},
		}),
		rune('|'): l.acceptFirst([]match{
			{str: "||", token: token.LOr},
			{str: "|=", token: token.OrAssign},
			{str: "|", token: token.Or},
		}),
		rune('^'): l.acceptFirst([]match{
			{str: "^=", token: token.XorAssign},
			{str: "^", token: token.Xor},
		}),
		rune('!'): l.acceptFirst([]match{
			{str: "!==", token: token.NotTEqual},
			{str: "!=", token: token.NotEqual},
			{str: "!", token: token.LNot},
		}),
		assign: l.acceptFirst([]match{
			{str: "===", token: token.TEqual},
			{str: "==", token: token.Equal},
			{str: "=", token: token.Assign},
		}),
		minusSign: l.acceptFirst([]match{
			{str: "--", token: token.Dec},
			{str: "-=", token: token.SubAssign},
			{str: "-", token: token.Minus},
		}),
		plusSign: l.acceptFirst([]match{
			{str: "++", token: token.Inc},
			{str: "+=", token: token.AddAssign},
			{str: "+", token: token.Plus},
		}),
	}
}

// acceptFirst takes a list of matches and returns the
// first matched token, if no match is found it is considered
// as an error and a illegal token is produced.
//
// This function is useful when multiple well know tokens
// starts with the same char, so you can search for the better
// match given a common start rune.
func (l *lexer) acceptFirst(matches []match) lexerState {
	return func() (Tokval, lexerState) {
		for _, m := range matches {
			tok, ok := l.accept(m)
			if ok {
				return tok, l.initialState
			}
		}
		return l.illegalToken()
	}
}

func (l *lexer) accept(m match) (Tokval, bool) {
	want := []rune(m.str)
	code := l.code[l.position:]

	if len(code) < len(want) {
		return Tokval{}, false
	}

	for i, r := range want {
		if r != code[i] {
			return Tokval{}, false
		}
	}

	l.position += uint(len(want) - 1)
	return l.token(m.token), true
}

func (l *lexer) dotState() (Tokval, lexerState) {
	l.fwd()
	if l.isTokenEnd() {
		return l.illegalToken()
	}
	allowExponent := true
	allowDot := false
	return l.decimalState(allowExponent, allowDot)
}

func (l *lexer) punctuator() (Tokval, lexerState) {
	return l.puncStates[l.cur()]()
}

func (l *lexer) isPunctuator() bool {
	_, ok := l.puncStates[l.cur()]
	return ok
}

func (l *lexer) updateLine() {
	l.line += 1
	l.column = 1
}

func (l *lexer) stringState() (Tokval, lexerState) {

	for !l.isEOF() && !l.isDoubleQuote() {
		if l.isNewline() {
			return l.illegalToken()
		}
		l.fwd()
	}

	if l.isEOF() {
		return l.illegalToken()
	}

	return l.stringToken(), l.initialState
}

func (l *lexer) numberState() (Tokval, lexerState) {

	if l.isEOF() {
		return l.token(token.Decimal), l.initialState
	}

	if l.isHexStart() {
		l.fwd()

		if l.isEOF() || l.isNewline() {
			return l.illegalToken()
		}

		return l.hexadecimalState()
	}

	allowExponent := true
	allowDot := true
	return l.decimalState(allowExponent, allowDot)
}

func (l *lexer) illegalToken() (Tokval, lexerState) {
	return Tokval{
		Type:  token.Illegal,
		Value: newStr(l.code),
	}, nil
}

func (l *lexer) identifierState() (Tokval, lexerState) {

	// TODO: handle keywords followed by dot and ( ? like null() ? or leave the parser to handle it ?

	for !l.isEOF() {

		if l.isDot() {
			l.bwd()
			return l.token(token.Ident), l.accessMemberState
		}

		if l.isLeftParen() || l.isTokenEnd() {
			l.bwd()
			return l.identOrKeywordToken(), l.initialState
		}

		l.fwd()
	}

	return l.identOrKeywordToken(), l.initialState
}

func (l *lexer) identOrKeywordToken() Tokval {
	val := l.curValue()
	keywordType, isKeyword := keywords[string(val)]
	if isKeyword {
		return l.token(keywordType)
	}
	return l.token(token.Ident)
}

func (l *lexer) startIdentifierState() (Tokval, lexerState) {

	if l.isEOF() {
		return EOF, nil
	}

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
		if l.isTokenEnd() {
			l.bwd()
			return l.token(token.Hexadecimal), l.initialState
		}
		if !l.isHexadecimal() {
			return l.illegalToken()
		}
		l.fwd()
	}

	return l.token(token.Hexadecimal), l.initialState
}

func (l *lexer) decimalState(allowExponent bool, allowDot bool) (Tokval, lexerState) {

	for !l.isEOF() {
		if l.isExponentPartStart() {
			if !allowExponent {
				return l.illegalToken()
			}
			l.fwd()
			return l.exponentPartState()
		}

		if l.isDot() {
			if !allowDot {
				return l.illegalToken()
			}
			l.fwd()
			return l.decimalState(allowExponent, false)
		}

		if l.isTokenEnd() {
			l.bwd()
			return l.token(token.Decimal), l.initialState
		}

		if !l.isNumber() {
			return l.illegalToken()
		}

		l.fwd()
	}

	return l.token(token.Decimal), l.initialState
}

func (l *lexer) exponentPartState() (Tokval, lexerState) {

	if l.isTokenEnd() {
		return l.illegalToken()
	}

	if l.isMinusSign() || l.isPlusSign() {
		// TODO: test 1.0e- and 1.0e+
		l.fwd()
	}

	allowExponent := false
	allowDot := true
	return l.decimalState(allowExponent, allowDot)
}

func (l *lexer) skipSpaces() {
	for l.isNewline() || l.isWhiteSpace() {
		if l.isNewline() {
			l.updateLine()
		} else {
			l.updateColumn()
		}
		l.consume()
	}
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

func (l *lexer) isNewline() bool {
	if l.isEOF() {
		return false
	}
	return containsRune(lineTerminators, l.cur())
}

func (l *lexer) isWhiteSpace() bool {
	if l.isEOF() {
		return false
	}
	return containsRune(whiteSpaces, l.cur())
}

func (l *lexer) isHexadecimal() bool {
	return containsRune(hexnumbers, l.cur())
}

func (l *lexer) isExponentPartStart() bool {
	return containsRune(exponentPartStart, l.cur())
}

func (l *lexer) isComma() bool {
	return l.cur() == comma
}

func (l *lexer) isDoubleQuote() bool {
	return l.cur() == doubleQuote
}

func (l *lexer) isSemiColon() bool {
	return l.cur() == semiColon
}

// tokenEnd tries to capture the most common causes of a token ending
func (l *lexer) isTokenEnd() bool {
	if l.isEOF() {
		return true
	}
	return l.isRightParen() || l.isComma() || l.isNewline() || l.isSemiColon() || l.isWhiteSpace()
}

func (l *lexer) fwd() {
	l.position += 1
}

func (l *lexer) bwd() {
	l.position -= 1
}

// curValue returns all the runes that compose the current token
// being analised instead of just the current rune.
// It has no side effects.
func (l *lexer) curValue() []rune {

	if l.isEOF() {
		return l.code
	}
	return l.code[:l.position+1]
}

func (l *lexer) consume() {
	if l.isEOF() {
		l.code = nil
	} else {
		l.code = l.code[l.position+1:]
	}
	l.position = 0
}

// token will generate a token consuming all the code
// until the current position. After calling this method
// the token will not be available anymore (it has been consumed)
// on the code and the position will be reset to zero.
func (l *lexer) token(t token.Type) Tokval {

	val := l.curValue()
	column := l.updateColumn()
	l.consume()

	return Tokval{Type: t, Value: newStr(val), Line: l.line, Column: column}
}

func (l *lexer) updateColumn() uint {
	column := l.column
	l.column += l.position + 1
	return column
}

func (l *lexer) stringToken() Tokval {
	// WHY: we need to remove the double quotes
	// around the string.

	val := l.code[1:l.position]

	column := l.updateColumn()
	l.consume()

	return Tokval{
		Type:   token.String,
		Value:  newStr(val),
		Line:   l.line,
		Column: column,
	}
}

var numbers []rune
var hexnumbers []rune
var lineTerminators []rune
var whiteSpaces []rune
var linefeed rune
var carriageRet rune
var semiColon rune
var dot rune
var minusSign rune
var plusSign rune
var leftParen rune
var rightParen rune
var comma rune
var doubleQuote rune
var assign rune
var hexStart []rune
var exponentPartStart []rune
var keywords map[string]token.Type

func init() {
	numbers = []rune("0123456789")
	hexnumbers = append(numbers, []rune("abcdefABCDEF")...)
	linefeed = rune('\u000A')
	carriageRet = rune('\u000D')
	lineSep := rune('\u2028')
	paragraphSep := rune('\u2029')
	lineTerminators = []rune{linefeed, carriageRet, lineSep, paragraphSep}
	dot = rune('.')
	minusSign = rune('-')
	plusSign = rune('+')
	leftParen = rune('(')
	rightParen = rune(')')
	comma = rune(',')
	doubleQuote = rune('"')
	semiColon = rune(';')
	hexStart = []rune("xX")
	exponentPartStart = []rune("eE")
	assign = rune('=')
	keywords = newKeywords()
	whiteSpaces = newWhiteSpaces()
}

func newKeywords() map[string]token.Type {
	return map[string]token.Type{
		"null":       token.Null,
		"undefined":  token.Undefined,
		"false":      token.Bool,
		"true":       token.Bool,
		"break":      token.Break,
		"case":       token.Case,
		"continue":   token.Continue,
		"debugger":   token.Debugger,
		"default":    token.Default,
		"delete":     token.Delete,
		"do":         token.Do,
		"else":       token.Else,
		"finally":    token.Finally,
		"for":        token.For,
		"function":   token.Function,
		"if":         token.If,
		"in":         token.In,
		"instanceof": token.InstanceOf,
		"new":        token.New,
		"return":     token.Return,
		"switch":     token.Switch,
		"this":       token.This,
		"throw":      token.Throw,
		"try":        token.Try,
		"typeof":     token.TypeOf,
		"var":        token.Var,
		"void":       token.Void,
		"while":      token.While,
		"with":       token.With,
	}
}

func newWhiteSpaces() []rune {

	tab := rune('\u0009')
	verticalTab := rune('\u000B')
	formFeed := rune('\u000C')
	space := rune('\u0020')
	noBreakSpace := rune('\u00A0')
	byteOrderMark := rune('\uFEFF')

	return []rune{tab, verticalTab, formFeed, space, noBreakSpace, byteOrderMark}
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
