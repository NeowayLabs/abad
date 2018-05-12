package types

type (
	BuiltinFunc func(args Value) (Value, error)
)