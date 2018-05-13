package types

import (
	"math"

	"github.com/NeowayLabs/abad/internal/utf16"
)

type Undefined utf16.Str // yeah, science!

var Undef = Undefined(utf16.S("undefined"))

func (u Undefined) IsTrue() bool {
	return false
}

func (u Undefined) IsFalse() bool {
	return true
}

func (u Undefined) ToBool() Bool {
	return Bool(u.IsTrue())
}

func (u Undefined) ToString() String {
	return NewString(utf16.Str(u).String())
}

func (u Undefined) ToNumber() Number {
	return Number(math.NaN())
}

func (_ Undefined) Kind() Kind {
	return KindUndefined
}