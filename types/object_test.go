package types_test

import (
	"testing"

	"github.com/NeowayLabs/abad/types"
)

func TestObjectNew(t *testing.T) {
	obj := types.NewRawObject()
	proto := obj.GetProperty(utf16.S("prototype"))
	if !types.StrictEqual(proto, types.Null) {
		t.Fatalf("Raw Object extends Null type")
	}
}