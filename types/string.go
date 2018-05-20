package types

import (
	"math"
	"strconv"

	"github.com/NeowayLabs/abad/internal/utf16"
)

type (
	String utf16.Str
)

func NewString(str string) String {
	return String(utf16.Encode(str))
}

func (a String) ToPrimitive(hint Kind) (Value, error) { return a, nil }

func (a String) ToObject() (Object, error) {
	panic("not implemented yet")
}

func (a String) Length() int {
	return len(a)
}

func (a String) IsTrue() bool {
	return !a.IsFalse()
}

func (a String) IsFalse() bool {
	return a.Length() == 0
}

func (a String) ToBool() Bool {
	return Bool(a.IsTrue())
}

func (a String) ToNumber() Number {
	n, err := strconv.ParseFloat(utf16.Str(a).String(), 64)
	if err != nil {
		return NewNumber(math.NaN())
	}
	return NewNumber(n)
}

func (a String) ToString() String {
	return a
}

func (a String) String() string {
	return utf16.Str(a).String()
}

func (a String) Kind() Kind {
	return KindString
}

func (a String) Equal(b String) bool {
	return utf16.Str(a).Equal(utf16.Str(b))
}