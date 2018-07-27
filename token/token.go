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
	Mul
	Rem
	Quo
	Inc
	Dec
	Dot
	LParen
	RParen
	Comma
	SemiColon
	LBrace
	RBrace
	LBrack
	RBrack
	Less
	Greater
	LessEq
	GreaterEq
	Equal
	NotEqual
	TEqual // stands for triple equal x_x
	NotTEqual
	LShift
	RShift
	RShiftZero // Right Shift with zero filling, yay JS
	And
	Or
	Xor
	Not
	LNot
	LAnd
	LOr
	Colon
	Assign
	AddAssign
	SubAssign
	MulAssign
	RemAssign
	QuoAssign
	LShiftAssign
	RShiftAssign
	RShiftZeroAssign
	AndAssign
	OrAssign
	XorAssign
	Ternary

	Ident

	Null
	Undefined

	Break
	Case
	Catch
	Continue
	Debugger
	Default
	Delete
	Do
	Else
	Finally
	For
	Function
	If
	In
	InstanceOf
	New
	Return
	Switch
	This
	Throw
	Try
	TypeOf
	Var
	Void
	While
	With

	EOF
)

var names = map[Type]string{
	Illegal:          "Illegal",
	Decimal:          "Decimal",
	Hexadecimal:      "Hexadecimal",
	Octal:            "Octal",
	String:           "String",
	Bool:             "Bool",
	Minus:            "-",
	Plus:             "+",
	Mul:              "*",
	Rem:              "%",
	Quo:              "/",
	Inc:              "++",
	Dec:              "--",
	LBrace:           "{",
	RBrace:           "}",
	LBrack:           "[",
	RBrack:           "[",
	Less:             "<",
	Greater:          ">",
	LessEq:           "<=",
	GreaterEq:        ">=",
	Equal:            "==",
	NotEqual:         "!=",
	TEqual:           "===",
	NotTEqual:        "!==",
	LShift:           "<<",
	RShift:           ">>",
	RShiftZero:       ">>>",
	And:              "&",
	Or:               "|",
	Xor:              "^",
	Not:              "~",
	LNot:             "!",
	LAnd:             "&&",
	LOr:              "||",
	Colon:            ":",
	Assign:           "=",
	AddAssign:        "+=",
	SubAssign:        "-=",
	MulAssign:        "*=",
	RemAssign:        "%=",
	QuoAssign:        "/=",
	LShiftAssign:     "<<=",
	RShiftAssign:     ">>=",
	RShiftZeroAssign: ">>>=",
	AndAssign:        "&=",
	OrAssign:         "|=",
	XorAssign:        "^=",
	Ternary: "?",
	Dot:              ".",
	LParen:           "(",
	RParen:           ")",
	Comma:            ",",
	Ident:            "Ident",
	SemiColon:        "SemiColon",
	Null:             "Null",
	Undefined:        "Undefined",
	Break:            "Break",
	Case:             "Case",
	Catch:            "Catch",
	Continue:         "Continue",
	Debugger:         "Debugger",
	Default:          "Default",
	Delete:           "Delete",
	Do:               "Do",
	Else:             "Else",
	Finally:          "Finally",
	For:              "For",
	Function:         "Function",
	If:               "If",
	In:               "In",
	InstanceOf:       "InstanceOf",
	New:              "New",
	Return:           "Return",
	Switch:           "Switch",
	This:             "This",
	Throw:            "Throw",
	Try:              "Try",
	TypeOf:           "TypeOf",
	Var:              "Var",
	Void:             "Void",
	While:            "While",
	With:             "With",
	EOF:              "EOF",
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
