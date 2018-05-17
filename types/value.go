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

func (k Kind) String() string {
	switch k {
	case KindUndefined: return "undefined"
	case KindNull: return "null"
	case KindNumber: return "number"
	case KindString: return "string"
	case KindBool: return "bool"
	case KindObject: return "object"
	}

	panic("unrecognized type")
	return ""
}

func StrictEqual(a, b Value) bool {
	akind := a.Kind()
	bkind := b.Kind()

	if akind != bkind {
		return false
	}

	if akind == KindUndefined || 
		akind == KindNull {
		return true
	}

	if akind == KindNumber {
		an := a.(Number)
		bn := b.(Number)

		return an.Equal(bn)
	}

	if akind == KindString {
		as := a.(String)
		bs := b.(String)
		return as.Equal(bs)
	}

	if akind == KindBool {
		ab := a.(Bool)
		bb := b.(Bool)
		return ab.Equal(bb)
	}

	if akind == KindObject {	
		aobj := a.(*Object)
		bobj := b.(*Object)
		return aobj == bobj // pointer comparison
	}

	panic("strict equal not implemented")

	return false
}

func IsPrimitive(val Value) bool {
	switch val.Kind() {
	case KindUndefined, KindNull, KindNumber, KindString, KindBool:
		return true
	}

	return false
}