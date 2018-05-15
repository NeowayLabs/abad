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

	// Object is a collection of named objects.
	Object struct {
		// Class is the kind of object
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
	o := newRawObject()
	o.put(valueAttr, value)
	o.put(writableAttr, Bool(wrt))
	o.put(enumAttr, Bool(enum))
	o.put(cfgAttr, Bool(cfg))
	return o
}

func NewAcessorPropDesc(get, set Value, enum, cfg bool) *Object {
	o := newRawObject()
	o.put(getAttr, get)
	o.put(setAttr, set)
	o.put(enumAttr, Bool(enum))
	o.put(cfgAttr, Bool(cfg))
	return o
}

// https://es5.github.io/#x8.6.1
func DefaultDataPropDesc() *Object {
	return NewDataPropDesc(Undef, false, false, false)
}

// https://es5.github.io/#x8.6.1
func DefaultAcessorPropDesc() *Object {
	return NewAcessorPropDesc(Undef, Undef, false, false)
}

// NewObject creates a new Object using proto as
// prototype.
func NewObject(proto Value) (*Object, error) {
	o := newRawObject()
	err = o.DefineOwnProperty(protoAttr,
		NewDataPropDesc(proto, false, true, false), true)
	if err != nil {
		return nil, err
	}
	return o, nil
}

// newRawObject creates a prototypeless object.
// Cannot be exposed to ECMAScript.
func newRawObject() *Object {
	return &Object{
		props: make(map[string]*Object),
	}
}

func (o *Object) Get(name utf16.Str) (Value, error) {
	if o.Class == "Function" {
		return o.functionGet(name)
	}
	return o.genericGet(name)
}

// genericGet is the default [[Get]] implementation for objects.
// https://es5.github.io/#x8.12.3
func (o *Object) genericGet(name utf16.Str) (Value, error) {
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

// functionGet implements [[Get]] for Function.
func (o *Object) functionGet(name utf16.Str) (Value, error) {
	v, err := o.genericGet(name)
	if err != nil {
		return nil, err
	}

	if name.Equal(utf16.S("caller")) {
		// TODO(i4k): throw TypeError
		return nil, NewTypeError("property caller is unaceptable")
	}

	return v, nil
}

func (o *Object) Put(name utf16.Str, val Value, failure bool) error {
	return nil
}

func (o *Object) put(name utf16.Str, val Value) {
	o.props[name.String()] = val
}

func (o *Object) CanPut(name utf16.Str) bool {
	return false
}

func (o *Object) GetOwnProperty(name utf16.Str) Value {
	prop := o.Get(name)
	if prop == Undef {
		return Undef
	}

	desc := newRawObject()
	if IsDataDescriptor(prop) {
		desc.put(valueAttr, prop.Get(valueAttr))
	}
}

func (o *Object) GetProperty(name utf16.Str) Value {
	return Undef
}

// https://es5.github.io/#x8.12.9
func (o *Object) DefineOwnProperty(
	name utf16.Str, desc Value, throw bool,
) (bool, error) {
	// throw exception if requested, otherwise quietly returns
	retOrThrow := func(err error) (bool, error) {
		if err != nil {
			if throw {
				return false, err
			}

			return false, nil
		}

		return true, nil
	}

	if desc.Kind() != KindObject {
		return retOrThrow(NewTypeError(
			"DefineOwnProperty expects a PropertyDescriptor object",
		))
	}

	objdesc := desc.(*Object)

	current := o.GetOwnProperty(name)
	if current.Kind() != KindObject {
		panic("internal error: property is not a descriptor")
	}

	extensible := o.Extensible
	if StrictEqual(current, Undef) {
		if !extensible {
			return throw(NewTypeError("Object %s is not extensible",
				o.Class))
		}

		return o.setOwnProperty()
	}

	if isAbsentDescriptor(objdesc) {
		return true, nil
	}

	// uses internal SameValue(x, y)
	// TODO
	if isSameDescriptor(objdesc, current) {
		return true, nil
	}

	curCfg := current.Get(cfgAttr).ToBool()
	curEnum := current.Get(enumAttr).ToBool()
	curWr := current.Get(writableAttr).ToBool()

	descCfg := objdesc.Get(cfgAttr).ToBool()
	descEnum := objdesc.Get(enumAttr).ToBool()
	descWr := objdesc.Get(writableAttr).ToBool()

	if !curCfg {
		if descCfg {
			return throw(NewTypeError("configurable is false"))
		}

		if descEnum != curEnum {
			return throw(
				NewTypeError("enumerable dont match for configuration disabled"),
			)
		}
	}

	if IsDataDescriptor(current) != IsDataDescriptor(objdesc) {
		if !curCfg {
			return throw(NewTypeError("configurable is false, cannot" +
				" change from data descriptor to acessor, and vice-versa"))
		}

		var newdesc *Object

		if IsDataDescriptor(current) {
			newdesc = DefaultAcessorPropDesc()
		} else {
			newdesc = DefaultDataPropDesc()
		}

		err := newdesc.Put(enumAttr, curEnum, shouldThrow)
		if err != nil {
			return throw(err)
		}

		err = newdesc.Put(cfgAttr, curCfg, shouldThrow)
		if err != nil {
			return throw(err)
		}

		current = newdesc
	} else if IsDataDescriptor(current) && IsDataDescriptor(objdesc) {
		if !curCfg {
			if !curWr && descWr {
				return throw(
					NewTypeError("configurable is false and writable mismatch"),
				)
			}

			if !curWr {
				if objdesc.HasProperty(valueAttr) &&
					!SameValue(current.Get(valueAttr), objdesc.Get(valueAttr)) {
					return throw(NewTypeError("writable is false"))
				}
			}
		}
	}

	err := copyProperties(current, objdesc)
	if err != nil {
		return throw(NewTypeError(err.Error()))
	}

	return throw(o.Put(name, current, throw))
}

// setOwnProperty just sets the property. Calls from ECMAScript
// must invoke DefineOwnProperty that does the correct validations.
func (o *Object) setOwnProperty(name utf16.Str, desc *Object, throw bool) (bool, error) {
	retOrThrow := func(err error) (bool, error) {
		if err != nil {
			if throw {
				return false, err
			}
			return false, nil
		}

		return true, nil
	}

	if IsGenericDescriptor(desc) ||
		IsDataDescriptor(desc) {
		newdesc := DefaultDataDescProp()
		err := copyDataDesc(newdesc, desc, throw)
		if err != nil {
			return retOrThrow(err)
		}

		return retOrThrow(o.Put(name, newdesc, throw))
	}

	if !IsAcessorDescriptor(desc) {
		panic("descriptor must be generic, data or acessor")
	}

	newdesc := DefaultAcessorPropDesc()
	err := copyAcessorDesc(newdesc, desc, throw)
	if err != nil {
		return retOrThrow(err)
	}

	return retOrThrow(o.Put(name, newdesc, throw))
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

func copyDataDesc(dst, src *Object, throw bool) error {
	values := []utf16.Str{
		valueAttr,
		writableAttr,
		enumAttr,
		cfgAttr,
	}

	for _, attr := range values {
		v := src.Get(attr)
		if !StrictEqual(v, Undef) {
			err := dst.Put(attr, v, throw)
			if err != nil && throw {
				return false, err
			}
		}
	}

	return nil
}