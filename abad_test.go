package abad_test

import (
	"fmt"
	"testing"

	"github.com/NeowayLabs/abad"
	"github.com/NeowayLabs/abad/types"
	"github.com/madlambda/spells/assert"
)

var E = fmt.Errorf

func TestNumberEval(t *testing.T) {
	for _, tc := range []struct {
		code string
		obj  types.Value
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
		{
			code: "-+0",
			obj:  types.Number(-0.0),
		},
		{
			code: "-+-0",
			obj:  types.Number(0.0),
		},
		{
			code: "-+-+0",
			obj:  types.Number(0.0),
		},
		{
			code: "0.1.",
			err:  E("parser error: <anonymous>:1:0: invalid token: 0.1."),
		},
	} {
		js, err := abad.NewAbad()
		assert.NoError(t, err, "failed to start ecma")
		obj, err := js.EvalFile("<anonymous>", tc.code)
		assert.EqualErrs(t, tc.err, err, "errors differ")

		if err != nil {
			continue
		}

		got, ok := obj.(types.Number)
		if !ok {
			t.Fatalf("got value other than number: %s", obj)
		}

		want := tc.obj.(types.Number)
		assert.EqualFloats(t, float64(want), float64(got),
			"number differs")
	}
}

func TestIdentExprEval(t *testing.T) {
	for _, tc := range []struct {
		code string
		ret  string
		err  error
	}{
		{
			code: "console",
			ret:  "[object Object]",
		},
		{
			code: "angular",
			err:  E("[angular] is not defined"),
		},
	} {
		js, err := abad.NewAbad()
		assert.NoError(t, err, "failed to start interpreter")
		val, err := js.Eval(tc.code)
		assert.EqualErrs(t, tc.err, err, "errors differ")

		if err != nil {
			continue
		}

		obj := val.(types.Object)

		gotstr := obj.String()

		assert.EqualStrings(t, tc.ret, gotstr, "strings don't match")
	}
}

func TestStringEval(t *testing.T) {
	for _, tc := range []struct {
		name string
		code string
		want string
	}{
		{
			name: "Empty",
			code: `""`,
			want: "",
		},
		{
			name: "Gibberish",
			code: `"!@#$$%&*()_+={}^/;~]/.,"`,
			want: "!@#$$%&*()_+={}^/;~]/.,",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			js, err := abad.NewAbad()
			assert.NoError(t, err, "failed to start interpreter")
			val, err := js.Eval(tc.code)

			assert.NoError(t, err, "unexpected error evaluating code")
			str := val.(types.String)

			got := str.String()

			assert.EqualStrings(t, tc.want, got, "output does not match expectation")
		})
	}
}
