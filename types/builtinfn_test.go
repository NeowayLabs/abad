package types_test

import (
	"testing"

	"github.com/NeowayLabs/abad/types"
)

var Str = types.NewString

func TestBuiltin(t *testing.T) {
	for _, tc := range []struct {
		input  []types.Value
		fn     types.Execfn
		output types.Value
	}{
		{
			input: []types.Value{Str("hello"), Str("world")},
			fn: func(obj types.Object, args []types.Value) types.Value {
				return types.Undefined
			},
			output: types.Undefined,
		},
		{
			input: []types.Value{Str("hello"), Str("world")},
			fn: func(obj types.Object, args []types.Value) types.Value {
				return args[0]
			},
			output: Str("hello"),
		},
		{
			input: []types.Value{Str("hello"), Str("world")},
			fn: func(obj types.Object, args []types.Value) types.Value {
				return types.NewNumber(float64(len(args)))
			},
			output: types.NewNumber(2.0),
		},
	} {
		global := types.NewBaseDataObject()
		builtin := types.NewBuiltinfn(tc.fn)
		got := builtin.Call(global, tc.input)
		if !types.StrictEqual(tc.output, got) {
			t.Fatalf("values differ: '%s' != '%s'", got, tc.output)
		}
	}
}
