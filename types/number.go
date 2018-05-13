package types

import (
	"math"
	"strconv"
)

type (
	Number float64
)

var ε = math.Nextafter(1, 2) - 1

func NewNumber(a float64) Number {
	return Number(a)
}

func (a Number) Value() float64 { return float64(a) }

func (a Number) String() string {
	return strconv.FormatFloat(float64(a), 'f', -1, 64)
}

// https://es5.github.io/#x9.2
func (a Number) IsTrue() bool {
	return !a.IsFalse()
}

func (a Number) IsFalse() bool {
	return math.IsNaN(float64(a)) ||
		a == -0.0 ||
		a == +0.0
}

// ToBool returns a Boolean according to:
// https://es5.github.io/#x9.2
func (a Number) ToBool() Bool {
	if a.IsTrue() {
		return Bool(true)
	}

	return Bool(false)
}

// ToNumber retrieves the number.
// https://es5.github.io/#x9.3
func (a Number) ToNumber() Number {
	return a
}

// ToString converts the number to string.
// Check https://es5.github.io/#x9.8
// TODO(i4k): revisit this.
func (a Number) ToString() String {
	val := strconv.FormatFloat(float64(a), 'f', -1, 64)
	return NewString(val)
}

func (_ Number) Kind() Kind {
	return KindNumber
}

func (a Number) Equal(b Number) bool {
	if math.IsNaN(a.Value()) ||
		math.IsNaN(b.Value()) {
		return false
	}

	return equalValues(a.Value(), b.Value())
}

func equalValues(a, b float64) bool {
	return math.Abs(a-b) < ε && math.Abs(b-a) < ε
}
