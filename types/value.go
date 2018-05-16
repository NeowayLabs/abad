package types

type (
	Kind int

	// Value is a heavy interface.
	// In JS type coercion is crazy... almost everything
	// could be coerced to something else. Then, better to
	// encapsulate the question regarding what a value can be
	// in the same interface than dealing with concrete types
	// every time.
	// Every type implements the Value interface.
	Value interface {
		Kind() Kind

		IsTrue() bool
		IsFalse() bool

		ToPrimitive(hint Kind) (Value, error)
		ToBool() Bool
		ToNumber() Number
		ToString() String
	}
)

const (
	KindUndefined Kind = iota
	KindNull
	KindNumber
	KindString
	KindBool
	KindObject
)

func StrictEqual(a, b Value) bool {
	if a.Kind() != b.Kind() {
		return false
	}

	if a.Kind() == KindUndefined ||
		a.Kind() == KindNull {
		return true
	}

	if a.Kind() == KindNumber {
		an := a.(Number)
		bn := b.(Number)

		return an.Equal(bn)
	}

	if a.Kind() == KindString {
		as := a.(String)
		bs := b.(String)
		return as.Equal(bs)
	}

	if a.Kind() == KindBool {
		ab := a.(Bool)
		bb := b.(Bool)
		return ab.Equal(bb)
	}

	// TODO(i4k): implement Object comparison

	return false
}

func IsPrimitive(val Value) bool {
	switch val.Kind() {
	case KindUndefined, KindNull, KindNumber, KindString, KindBool:
		return true
	}

	return false
}