package types_test

import (
	"testing"

	"github.com/NeowayLabs/abad/internal/utf16"
	"github.com/NeowayLabs/abad/types"
	"github.com/madlambda/spells/assert"
)

var (
	S         = utf16.S
	protoAttr = S("prototype")
)

func TestBaseObjectExtendsNull(t *testing.T) {
	obj := types.NewBaseDataObject()
	proto, err := obj.Get(protoAttr)
	assert.NoError(t, err, "failed getting proto")

	if !types.StrictEqual(proto, types.Null) {
		t.Fatalf("Raw Object extends Null type")
	}
}

func TestNewObjectExtendsProto(t *testing.T) {
	proto := types.NewBaseDataObject()
	obj := types.NewDataObject(proto)

	gotproto, err := obj.Get(protoAttr)
	assert.NoError(t, err, "prototypes differs")

	if gotproto.Kind() != types.KindObject {
		t.Fatalf("got type %s", gotproto.Kind())
	}

	gotobj := gotproto.(*types.DataObject)
	if !types.StrictEqual(proto, gotobj) {
		t.Fatalf("%s and %s are not the same prototype",
			proto, gotobj)
	}
}

func TestObjectDefineOwnProperty(t *testing.T) {
	// new property never fails
	obj := types.NewBaseDataObject()
	propName := S("madlab")
	expected := types.True
	prop := types.NewDataPropDesc(expected, true, true, true)
	ok, err := obj.DefineOwnPropertyP(propName, prop, true)
	if !ok {
		t.Fatal(err)
	}

	gotval, err := obj.Get(propName)
	assert.NoError(t, err, "get failed")
	if !types.StrictEqual(expected, gotval) {
		t.Fatalf("got wrog value: %s", gotval)
	}

	gotprop := obj.GetOwnProperty(propName)
	if gotprop.Kind() != types.KindObject {
		t.Fatalf("got property of wrong type")
	}

	got := gotprop.(*types.DataObject)
	gotdesc := got.ToPropertyDescriptor()
	if !types.IsSameDescriptor(gotdesc, prop) {
		t.Fatalf("Property descriptors differ")
	}
}