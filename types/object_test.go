package types_test

import (
	"testing"

	"github.com/NeowayLabs/abad/internal/utf16"
	"github.com/NeowayLabs/abad/types"
	"github.com/madlambda/spells/assert"
)

type (
	DataTestcase struct {
		val types.Value
		wrt bool
		enu bool
		cfg bool
	}
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
	assert.NoError(t, err, "failed getting prototype")

	if gotproto.Kind() != types.KindObject {
		t.Fatalf("got type %s", gotproto.Kind())
	}

	gotobj := gotproto.(*types.DataObject)
	if !types.StrictEqual(proto, gotobj) {
		t.Fatalf("%s and %s are not the same prototype",
			proto, gotobj)
	}
}

func TestObjectDefineOwnPropertyDATA(t *testing.T) {
	for _, tc := range []DataTestcase{
		{val: types.True, wrt: true, enu: true, cfg: true},
		{val: types.Null, wrt: true, enu: true, cfg: true},
		{val: types.NewNumber(1.0), wrt: false, enu: true, cfg: true},
		{val: types.NewBaseDataObject(), wrt: true, enu: true, cfg: true},
		{val: types.NewBaseDataObject(), wrt: false, enu: true, cfg: true},
	} {
		obj := types.NewBaseDataObject()
		testDataDescriptor(t, obj, "madlab", tc)
		if tc.wrt {
			tc.val = types.NewNumber(666.0)
			testDataDescriptor(t, obj, "madlab", tc)
		} else {
			tc.val = types.NewNumber(123456.7) // improbable number
			testDataDescriptorFail(t, obj, "madlab", tc)
		}
	}
}

func testDataDescriptor(t *testing.T, obj *types.DataObject, property string, tc DataTestcase) {
	// new property never fails
	propName := S(property)
	expected := tc.val
	prop := types.NewDataPropDesc(expected, tc.wrt, tc.enu, tc.cfg)
	ok, err := obj.DefineOwnPropertyP(propName, prop, true)
	if !ok {
		t.Fatal(err)
	}

	gotval, err := obj.Get(propName)
	assert.NoError(t, err, "get failed")
	if !types.StrictEqual(expected, gotval) {
		t.Fatalf("got wrong value: %s", gotval)
	}

	gotprop := obj.GetOwnProperty(propName)
	if gotprop.Kind() != types.KindObject {
		t.Fatalf("Expected KindObject property got[%s]", gotprop.Kind())
	}

	got := gotprop.(*types.DataObject)
	gotdesc := got.ToPropertyDescriptor()
	if !types.IsSameDescriptor(gotdesc, prop) {
		t.Fatalf("Property descriptors differs: %+v != %+v", gotdesc, prop)
	}
}

func testDataDescriptorFail(t *testing.T, obj *types.DataObject, property string, tc DataTestcase) {
	propName := S(property)
	prop := types.NewDataPropDesc(tc.val, tc.wrt, tc.enu, tc.cfg)
	ok, _ := obj.DefineOwnPropertyP(propName, prop, true)
	if ok {
		t.Fatal("should fail")
	}
}
