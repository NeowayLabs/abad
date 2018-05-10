package abad_test

import (
	"testing"

	"github.com/NeowayLabs/abad"
	"github.com/NeowayLabs/abad/types"
	"github.com/madlambda/spells/assert"
)

func TestNumberEval(t *testing.T) {
	for _, tc := range []struct {
		code string
		obj  abad.Obj
		err  error
	}{
		{
			code: "0",
			obj:  types.Number(0.0),
		},
		{
			code: ".0",
			obj:  types.Number(0.0),
		},
		{
			code: "+0",
			obj:  types.Number(0.0),
		},
		{
			code: "-0",
			obj:  types.Number(-0.0),
		},
		{
			code: "1.0e10",
			obj:  types.Number(1.0e10),
		},
	} {
		js := abad.NewAbad("<anonymous>")
		obj, err := js.Eval(tc.code)
		assert.EqualErrs(t, tc.err, err, "errors differ")

		got, ok := obj.(types.Number)
		if !ok {
			t.Fatalf("got value other than number: %s", obj)
		}

		want := tc.obj.(types.Number)
		assert.EqualFloats(t, float64(want), float64(got), 
			"number differs")
	}
}