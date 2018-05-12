package types

type (
	// Value is a heavy interface.
	// In JS type coercion is crazy... almost everything
	// could be coerced to something else. Then, better to
	// encapsulate the question regarding what a value can be
	// in the same interface than dealing with concrete types
	// every time.
	// Every type implements the Value interface.
	Value interface {
		IsTrue() bool
		IsFalse() bool

		ToBool() Bool
		ToNumber() Number
		ToString() String
	}
)