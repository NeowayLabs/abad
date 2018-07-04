package envrec_test

import (
	"fmt"
	"testing"

	"github.com/NeowayLabs/abad/envrec"
	"github.com/NeowayLabs/abad/internal/utf16"
	"github.com/NeowayLabs/abad/types"
	"github.com/madlambda/spells/assert"
)

type testcase struct {
	ident string
	value types.Value
	err   error
}

var S = utf16.S
var E = fmt.Errorf

func TestEnvDecl(t *testing.T) {
	for _, tc := range []testcase{
		{
			ident: "console",
			value: types.NewNumber(1),
		},
		{
			ident: "window",
			value: types.NewNumber(666.0),
		},
		{
			ident: "",
			err:   E("empty binding name"),
		},
		{
			ident: "_",
			value: types.Undefined,
		},
		{
			ident: "$",
			value: types.NewString("jquery"),
		},
	} {
		testEnvRecDecl(t, tc)
	}
}

func testEnvRecDecl(t *testing.T, tc testcase) {
	ident := S(tc.ident)
	env := envrec.NewDeclEnv()
	err := env.New(ident, true)
	assert.EqualErrs(t, tc.err, err, "errs dont match")

	if err != nil {
		return
	}

	if !env.Has(ident) {
		t.Fatalf("binding not created")
	}

	err = env.Set(ident, tc.value, true)
	assert.NoError(t, err, "EnvDecl Set binding failed")

	got, err := env.Get(ident, true)
	assert.NoError(t, err, "DeclEnv Get binding failed")

	if !types.StrictEqual(got, tc.value) {
		t.Fatalf("conditional === failed. Got '%s' but expected '%s'",
			got, tc.value)
	}

	if !env.Del(ident) {
		t.Fatalf("Failed to delete DeclEnv mutable property")
	}

	if env.Has(ident) {
		t.Fatalf("DeclEnv still have a deleted binding")
	}
}
