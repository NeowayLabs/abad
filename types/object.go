package types

import (
	"fmt"
	"math"

	"github.com/NeowayLabs/abad/internal/utf16"
)

type (
	// DataObject is a collection of named values.
	DataObject struct {
		// Class is the kind of object
		class         string
		notExtensible bool
		props         map[string]*PropertyDescriptor
	}

	callable interface {
		Call(this *DataObject, args []Value) Value
	}
)

var (
	// S is an alias to make utf-16 strings.
	S            = utf16.S
	valueAttr    = S("value")
	writableAttr = S("writable")
	getAttr      = S("get")
	setAttr      = S("set")
	enumAttr     = S("enumerable")
	cfgAttr      = S("configurable")

	protoAttr    = S("prototype")
	toStringAttr = S("toString")
	valueOfAttr  = S("valueOf")
)

// DefaultPrototypeDesc is a base prototype object that extends Null.
// The root of the prototype-based type hierarchy.
func DefaultPrototypeDesc() *PropertyDescriptor {
	return NewDataPropDesc(Null, true, false, false)
}

// NewDataObject creates a new DataObject using proto as
// prototype.
func NewDataObject(proto Value) *DataObject {
	obj := NewBaseDataObject()

	// is must extend proto
	delete(obj.props, "prototype")

	// error ignored because it does not fail if
	// there's no previous properties.
	ok, err := obj.DefineOwnPropertyP(protoAttr,
		NewDataPropDesc(proto, false, true, false), true)

	if !ok {
		// should never occurs
		panic(err)
	}

	return obj
}

func NewBaseDataObject() *DataObject {
	return NewDataObjectP(DefaultPrototypeDesc())
}

// NewDataObjectP creates an object with PropertyDescriptor
// proto as prototype attribute.
// The *everything is an object* concept makes
// impossible to defineOwnProperties work by
// passing an object to define an object (recursive
// definition).
func NewDataObjectP(proto *PropertyDescriptor) *DataObject {
	obj := &DataObject{
		class: "object",
		props: make(map[string]*PropertyDescriptor),
	}

	obj.props["prototype"] = proto
	return obj
}

// Class returns the object class
func (o *DataObject) Class() string       { return o.class }
func (o *DataObject) NotExtensible() bool { return o.notExtensible }

// Value interface implementations

func (o *DataObject) IsFalse() bool { return false }
func (o *DataObject) IsTrue() bool  { return true }
func (_ *DataObject) Kind() Kind    { return KindObject }
func (_ *DataObject) ToBool() Bool  { return True }
func (o *DataObject) ToNumber() Number {
	primVal, err := o.ToPrimitive(KindNumber)
	if err != nil {
		return NewNumber(math.NaN())
	}

	return primVal.ToNumber()
}

func (o *DataObject) ToString() String {
	primVal, err := o.ToPrimitive(KindString)
	if err != nil {
		return NewString("")
	}

	return primVal.ToString()
}

func (o *DataObject) ToPrimitive(hint Kind) (Value, error) {
	return o.DefaultValue(hint)
}

func (o *DataObject) ToObject() (Object, error) {
	return o, nil
}

func (o *DataObject) ToPropertyDescriptor() *PropertyDescriptor {
	var (
		value, get, set     Value
		writable, enum, cfg Value
	)

	if o.HasProperty(valueAttr) {
		value, _ = o.Get(valueAttr)
	}

	if o.HasProperty(writableAttr) {
		b, _ := o.Get(writableAttr)
		writable = b.ToBool()
	}

	if o.HasProperty(getAttr) {
		get, _ = o.Get(getAttr)
	}

	if o.HasProperty(setAttr) {
		set, _ = o.Get(getAttr)
	}

	if o.HasProperty(enumAttr) {
		b, _ := o.Get(enumAttr)
		enum = b.ToBool()
	}

	if o.HasProperty(cfgAttr) {
		b, _ := o.Get(cfgAttr)
		cfg = b.ToBool()
	}

	if enum == nil {
		enum = DefEnumerable
	}

	if cfg == nil {
		cfg = DefConfigurable
	}

	if value != nil || writable != nil {
		if writable == nil {
			writable = DefWritable
		}

		if value == nil {
			value = DefValue
		}

		return NewDataPropDesc(value, writable.IsTrue(), enum.IsTrue(), cfg.IsTrue())
	} else if get != nil || set != nil {
		if get == nil {
			get = DefGet
		}

		if set == nil {
			set = DefSet
		}

		return NewAcessorPropDesc(get, set, enum.IsTrue(), cfg.IsTrue())
	}

	prop := NewGenericPropDesc()
	prop.put("enumerable", enum)
	prop.put("configurable", cfg)
	return prop
}

// Get is the default [[Get]] implementation for objects.
// https://es5.github.io/#x8.12.3
func (o *DataObject) Get(name utf16.Str) (Value, error) {
	desc, ok := o.getProperty(name)
	if !ok {
		return Undefined, nil
	}

	if desc.IsDataDescriptor() {
		return desc.Value(), nil
	}

	if !desc.IsAcessorDescriptor() {
		panic("descriptor is not data not acessor")
	}

	value := desc.Get()
	if StrictEqual(value, Undefined) {
		return Undefined, nil
	}

	getter, ok := value.(callable)
	if !ok {
		panic(fmt.Sprintf("object %s is not callable", getter))
	}

	return getter.Call(o, []Value{}), nil
}

func (o *DataObject) Put(name utf16.Str, val Value, throw bool) error {
	if !o.CanPut(name) {
		if throw {
			return NewTypeError("can not put data on this object")
		}

		return nil
	}

	ownDesc, ok := o.getOwnProperty(name)
	if ok && ownDesc.IsDataDescriptor() {
		// TODO(i4k): I'm not sure of this parameters
		desc := NewDataPropDesc(ownDesc.Value(), true, true, true)
		_, err := o.DefineOwnPropertyP(name, desc, throw)
		return err
	}

	desc, ok := o.getProperty(name)
	if !ok {
		desc := NewDataPropDesc(val, true, true, true)
		_, err := o.DefineOwnPropertyP(name, desc, throw)
		return err
	}

	if desc.IsDataDescriptor() {
		valueDesc := NewDataPropDesc(val, true, true, true)
		_, err := o.DefineOwnPropertyP(name, valueDesc, throw)
		return err
	}

	if desc.IsAcessorDescriptor() {
		set := desc.Set()
		if StrictEqual(set, Undefined) {
			panic("setter is undefined for acessor property")
		}

		setter, ok := set.(callable)
		if !ok {
			panic("setter is not a Function")
		}

		_ = setter.Call(o, []Value{val})
		return nil
	}

	panic("TODO(i4k): property is not an acessor nor data. Is this a problem?")
	return nil
}

func (o *DataObject) get(name utf16.Str) (*PropertyDescriptor, bool) {
	v, ok := o.props[name.String()]
	return v, ok
}

func (o *DataObject) put(name utf16.Str, val *PropertyDescriptor) {
	o.props[name.String()] = val
}

func (o *DataObject) CanPut(name utf16.Str) bool {
	desc, ok := o.getOwnProperty(name)
	if ok {
		if desc.IsAcessorDescriptor() {
			return !StrictEqual(desc.Set(), Undefined)
		} else if desc.IsDataDescriptor() {
			return desc.Writable().IsTrue()
		}

		panic("property is acessor nor data descriptor")
		return false
	}

	protodesc, ok := o.getOwnProperty(protoAttr)
	if !ok {
		panic("prototype not found")
	}

	proto := protodesc.Value()
	if StrictEqual(proto, Null) {
		return !o.NotExtensible()
	}

	if proto.Kind() != KindObject {
		panic(fmt.Sprintf("unexpected prototype value: %s", proto))
	}

	oproto := proto.(Object)
	inherited, ok := oproto.getProperty(name)
	if !ok {
		return !o.NotExtensible()
	}

	if inherited.IsAcessorDescriptor() {
		return !StrictEqual(inherited.Set(), Undefined)
	} else if inherited.IsDataDescriptor() {
		// TODO: Perhaps a o.IsExtensible here would be more clear
		return !o.NotExtensible()
	}

	panic("inherited isn't acessor not data descriptor")
	return false
}

func (o *DataObject) getOwnProperty(name utf16.Str) (*PropertyDescriptor, bool) {
	prop, ok := o.get(name)
	if !ok {
		return nil, false
	}
	return prop, true
}

func (o *DataObject) GetOwnProperty(name utf16.Str) Value {
	prop, ok := o.get(name)
	if !ok {
		return Undefined
	}

	return prop.ToObject()
}

func (o *DataObject) getProperty(name utf16.Str) (*PropertyDescriptor, bool) {
	prop, ok := o.getOwnProperty(name)
	if ok {
		return prop, true
	}

	protodesc, ok := o.getOwnProperty(protoAttr)
	if !ok {
		return nil, false
	}

	if !protodesc.HasValue() {
		return nil, false
	}

	protoval := protodesc.Value()

	if protoval.Kind() != KindObject {
		return nil, false
	}

	obj := protoval.(Object)
	return obj.getProperty(name)
}

func (o *DataObject) GetProperty(name utf16.Str) Value {
	prop, ok := o.getProperty(name)
	if ok {
		return prop.ToObject()
	}
	return Undefined
}

func (o *DataObject) DefineOwnProperty(
	name utf16.Str, desc Value, throw bool,
) (bool, error) {
	if desc.Kind() != KindObject {
		if throw {
			return false, NewTypeError(
				"DefineOwnProperty expects a PropertyDescriptor object",
			)
		}

		return false, nil
	}

	descobj := desc.(*DataObject)

	return o.DefineOwnPropertyP(name, descobj.ToPropertyDescriptor(), throw)
}

// https://es5.github.io/#x8.12.9
func (o *DataObject) DefineOwnPropertyP(
	name utf16.Str, desc *PropertyDescriptor, throw bool,
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

	notExtensible := o.notExtensible
	current, ok := o.getOwnProperty(name)
	if !ok {
		if notExtensible {
			return retOrThrow(NewTypeError("DataObject %s is not extensible",
				o.Class()))
		}

		return o.setOwnProperty(name, desc, throw)
	}

	if desc.IsAbsentDescriptor() {
		return true, nil
	}

	// uses internal SameValue(x, y)
	// TODO
	if IsSameDescriptor(desc, current) {
		return true, nil
	}

	curCfg := current.Cfg().ToBool()
	curEnum := current.Enum().ToBool()
	curWr := current.Writable().ToBool()

	descCfg := desc.Cfg().ToBool()
	descEnum := desc.Enum().ToBool()
	descWr := desc.Writable().ToBool()

	if !curCfg {
		if descCfg {
			return retOrThrow(NewTypeError("configurable is false"))
		}

		if descEnum != curEnum {
			return retOrThrow(
				NewTypeError("enumerable dont match for configuration disabled"),
			)
		}
	}

	if current.IsDataDescriptor() != desc.IsDataDescriptor() {
		if !curCfg {
			return retOrThrow(NewTypeError("configurable is false, cannot" +
				" change from data descriptor to acessor, and vice-versa"))
		}

		var newdesc *PropertyDescriptor

		if current.IsDataDescriptor() {
			newdesc = DefaultAcessorPropDesc()
		} else {
			newdesc = DefaultDataPropDesc()
		}

		newdesc.SetEnum(curEnum)
		newdesc.SetCfg(curCfg)

		current = newdesc
	} else if current.IsDataDescriptor() && desc.IsDataDescriptor() {
		if !curCfg {
			if !curWr && descWr {
				return retOrThrow(
					NewTypeError("configurable is false and writable mismatch"),
				)
			}

			if !curWr {
				if desc.HasValue() &&
					//TODO(i4k): SameValue() ?
					!StrictEqual(current.Value(), desc.Value()) {
					return retOrThrow(NewTypeError("writable is false"))
				}
			}
		}
	}

	CopyProperties(current, desc)
	o.put(name, current)
	return true, nil
}

// setOwnProperty just sets the property. Calls from ECMAScript
// must invoke DefineOwnProperty that does the correct validations.
func (o *DataObject) setOwnProperty(name utf16.Str, desc *PropertyDescriptor, throw bool) (bool, error) {
	if desc.IsGenericDescriptor() ||
		desc.IsDataDescriptor() {
		newdesc := DefaultDataPropDesc()
		CopyProperties(newdesc, desc)
		o.put(name, newdesc)
		return true, nil
	}

	if !desc.IsAcessorDescriptor() {
		panic("descriptor must be generic, data or acessor")
	}

	newdesc := DefaultAcessorPropDesc()
	CopyProperties(newdesc, desc)
	o.put(name, newdesc)
	return true, nil
}

func (o *DataObject) HasProperty(name utf16.Str) bool {
	prop := o.GetProperty(name)
	return !StrictEqual(prop, Undefined)
}

// https://es5.github.io/#x8.12.8
func (o *DataObject) DefaultValue(hint Kind) (Value, error) {
	if hint == KindString {
		//TODO(i4k): || hint == KindDate {
		return o.defaultString()
	}

	return o.defaultNumber()
}

func (o *DataObject) defaultString() (Value, error) {
	toString, _ := o.Get(toStringAttr)
	if stringify, ok := toString.(Function); ok {
		str := stringify.Call(o, []Value{})
		if IsPrimitive(str) {
			return str, nil
		}
	}

	valueOf, _ := o.Get(valueOfAttr)
	if valueFunc, ok := valueOf.(callable); ok {
		val := valueFunc.Call(o, []Value{})
		if IsPrimitive(val) {
			return val, nil
		}
	}

	return nil, NewTypeError("DataObject has no defaultValue")
}

func (o *DataObject) defaultNumber() (Value, error) {
	valueOf, _ := o.Get(valueOfAttr)
	if valuefunc, ok := valueOf.(callable); ok {
		val := valuefunc.Call(o, []Value{})
		if IsPrimitive(val) {
			return val, nil
		}
	}

	tostring, _ := o.Get(toStringAttr)
	if stringify, ok := tostring.(callable); ok {
		str := stringify.Call(o, []Value{})
		if IsPrimitive(str) {
			return str, nil
		}
	}

	return nil, NewTypeError("DataObject has no defaultValue")
}

func (o *DataObject) String() string {
	v, err := o.defaultString()
	if err != nil {
		panic(err)
	}

	return v.ToString().String()
}
