package types

import "fmt"

type (
	// PropertyDescriptor describes an Object property.
	// An Object is a collection of properties and each
	// property is either a `Named Data Property`, a
	// `Named Acessor Property` or an internal property.
	// PropertyDescriptor describes all of them, why?
	// By ECMA 8.6, the property type should be interpreted
	// by their attributes set (see IsAcessorDescriptor in 8.10.1
	// and IsDataDescriptor in 8.10.2).
	// A Named Data Property associates a name with an ECMAScript
	// value. Eg.:
	//   this.name = "value";
	// The code above is the same:
	//   Object.defineOwnProperty(this, {
	//       value: "value",
	//       writable: true,
	//       enumerable: true,
	//       configurable: true
	//   });
	PropertyDescriptor *Object

	Object struct {
		Class      string
		Extensible bool

		props map[string]*Object
	}
)

var (
	valueAttr    = utf16.S("value")
	writableAttr = utf16.S("writable")
	getAttr      = utf16.S("get")
	setAttr      = utf16.S("set")
	enumAttr     = utf16.S("enumerable")
	cfgAttr      = utf16.S("configurable")

	protoAttr = utf16.S("prototype")
)

func NewDataPropDesc(value Value, wrt, enum, cfg bool) *Object {
	o := &Object{
		props: make(map[string]*Object),
	}

	o.props["value"] = value
	o.props["writable"] = Bool(wrt)
	o.props["enumerable"] = Bool(enum)
	o.props["configurable"] = Bool(cfg)
	return o
}

func NewObject(proto Value) *Object {
	o := &Object{
		props: make(map[string]*Object),
	}

	err = o.DefineOwnProperty(protoAttr, 
			NewDataPropDesc(proto, false, true,false), true)
	return o
}

func (o *Object) Get(name utf16.Str) (Value, error) {
	if o.Class == "Function" {
		return o.getfn(name)
	}
	return o.get(name)
}

// get is the default [[Get]] implementation for objects.
// https://es5.github.io/#x8.12.3
func (o *Object) get(name utf16.Str) (Value, error) {
	desc := o.GetProperty(name)
	if desc.Kind() == KindUndefined {
		return Undef, nil
	}

	if isDataDescriptor(desc) {
		return desc.Get(valueAttr)
	}

	if !isAcessorDescriptor(desc) {
		panic("descriptor is not data not acessor")
	}

	getter, _ := desc.Get(getAttr)
	if getter.Kind() == KindUndefined {
		return Undef, nil
	}

	if getter.Class == "Function" {
		panic(fmt.Sprintf("object %s is not callable", getter))
	}

	return getter.Call(o), nil
}

func (o *Object) getfn(name utf16.Str) (Value, error) {
	v, err := o.get(name)
	if err != nil {
		return nil, err
	}

	if name.Equal(utf16.S("caller")) {
		// TODO(i4k): throw TypeError
		return nil, fmt.Errorf("TypeError exception")
	}

	return v, nil
}

func (o *Object) Put(name utf16.Str, val Value, failure bool) error {
	return nil
}

func (o *Object) CanPut(name utf16.Str) bool {
	return false
}

func (o *Object) GetOwnProperty(name utf16.Str) Value {
	return Undef
}

func (o *Object) GetProperty(name utf16.Str) Value {
	return Undef
}

func (o *Object) DefineOwnProperty(
	name utf16.Str, desc Value, failure bool,
) error {
	return nil
}

func (o *Object) HasProperty(name utf16.Str) bool {
	return false
}

func (o *Object) Delete(name utf16.Str, failure bool) error {
	return nil
}

func (o *Object) DefaultValue(hint utf16.Str) Value {
	return Undef
}