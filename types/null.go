package types

import (
	"github.com/NeowayLabs/abad/internal/utf16"
)

type null utf16.Str

var Null = null(utf16.S("null"))

func (n null) ToPrimitive(_ Kind) (Value, error) {
	return n, nil
}

func (_ null) ToObject() (Object, error) {
	return nil, NewTypeError("cannot convert to Object")
}

func (_ null) Kind() Kind       { return KindNull }
func (_ null) IsFalse() bool    { return true }
func (_ null) IsTrue() bool     { return false }
func (_ null) ToBool() Bool     { return False }
func (_ null) ToNumber() Number { return NewNumber(+0) }
func (_ null) ToString() String { return String(Null) }
