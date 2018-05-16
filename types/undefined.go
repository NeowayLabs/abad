package types

import (
	"math"

	"github.com/NeowayLabs/abad/internal/utf16"
)

type undefined utf16.Str // yeah, science!

var Undefined = undefined(utf16.S("undefined"))

func (u undefined) IsTrue() bool {
	return false
}

func (u undefined) IsFalse() bool {
	return true
}

func (u undefined) ToBool() Bool {
	return Bool(u.IsTrue())
}

func (u undefined) ToString() String {
	return NewString(utf16.Str(u).String())
}

func (u undefined) ToNumber() Number {
	return Number(math.NaN())
}

func (_ undefined) Kind() Kind {
	return KindUndefined
}

func (_ undefined) ToPrimitive(hint Kind) (Value, error) {
	return Undefined, nil
}