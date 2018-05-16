package types

import (
	"fmt"
	"math"

	"github.com/NeowayLabs/abad/ast"
	"github.com/NeowayLabs/abad/internal/utf16"
)

type (
	// Object is a collection of named objects.
	Object struct {
		// Class is the kind of object
		Class      string
		Extensible bool

		Scope  *Object
		Params []ast.Ident
		Code   *ast.Program

		props map[string]*PropertyDescriptor
	}

	callable interface {
		Call(this *Object, args []Value) Value
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

	S = utf16.S
)

func DefaultPrototypeDesc() *PropertyDescriptor {
	return NewDataPropDesc(Null, true, false, false)
}

// NewObject creates a new Object using proto as
// prototype.
func NewObject(proto Value) *Object {
	o := NewRawObject()

	// error ignored because it does not fail if
	// there's no previous properties.
	o.DefineOwnPropertyP(protoAttr,
		NewDataPropDesc(proto, false, true, false), true)
	return o
}

// NewRawObject creates a prototypeless object.
func NewRawObject() *Object {
	obj := &Object{
		props: make(map[string]*PropertyDescriptor),
	}

	obj.props["prototype"] = DefaultPrototypeDesc()
	return obj
}

func (o *Object) IsFalse() bool { return false }
func (o *Object) IsTrue() bool  { return true }
func (_ *Object) Kind() Kind    { return KindObject }
func (_ *Object) ToBool() Bool  { return True }
func (o *Object) ToNumber() Number {
	primVal, err := o.ToPrimitive(KindNumber)
	if err != nil {
		return NewNumber(math.NaN())
	}

	return primVal.ToNumber()
}

func (o *Object) ToString() String {
	primVal, err := o.ToPrimitive(KindString)
	if err != nil {
		return NewString("")
	}

	return primVal.ToString()
}

func (o *Object) ToPrimitive(hint Kind) (Value, error) {
	return o.DefaultValue(hint)
}

func (o *Object) ToPropertyDescriptor() *PropertyDescriptor {
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

func (o *Object) Get(name utf16.Str) (Value, error) {
	if o.Class == "Function" {
		return o.functionGet(name)
	}
	return o.genericGet(name)
}

// genericGet is the default [[Get]] implementation for objects.
// https://es5.github.io/#x8.12.3
func (o *Object) genericGet(name utf16.Str) (Value, error) {
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

	getter := value.(*Object)
	if getter.Class == "Function" {
		panic(fmt.Sprintf("object %s is not callable", getter))
	}

	return getter.Call(o, []Value{}), nil
}

// functionGet implements [[Get]] for Function.
func (o *Object) functionGet(name utf16.Str) (Value, error) {
	v, err := o.genericGet(name)
	if err != nil {
		return nil, err
	}

	if name.Equal(utf16.S("caller")) {
		// TODO(i4k): throw TypeError
		return nil, NewTypeErrorS("property caller is unaceptable")
	}

	return v, nil
}

func (o *Object) Put(name utf16.Str, val Value, failure bool) error {
	return nil
}

func (o *Object) get(name utf16.Str) (*PropertyDescriptor, bool) {
	v, ok := o.props[name.String()]
	return v, ok
}

func (o *Object) put(name utf16.Str, val *PropertyDescriptor) {
	o.props[name.String()] = val
}

func (o *Object) CanPut(name utf16.Str) bool {
	return false
}

func (o *Object) Call(this *Object, args []Value) Value {
	return Undefined
}



func (o *Object) getOwnProperty(name utf16.Str) (*PropertyDescriptor, bool) {
	prop, ok := o.get(name)
	if !ok {
		return nil, false
	}
	return prop, true
}

func (o *Object) GetOwnProperty(name utf16.Str) Value {
	prop, ok := o.get(name)
	if !ok {
		return Undefined
	}

	return prop.ToObject()
}

func (o *Object) getProperty(name utf16.Str) (*PropertyDescriptor, bool) {
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

	obj := protoval.(*Object)

	return obj.getProperty(name)
}

func (o *Object) GetProperty(name utf16.Str) Value {
	prop := o.GetOwnProperty(name)
	if !StrictEqual(prop, Undefined) {
		return prop
	}

	proto := o.GetOwnProperty(protoAttr)
	if proto.Kind() != KindObject {
		return Undefined
	}

	obj := proto.(*Object)
	return obj.GetProperty(name)
}

func (o *Object) DefineOwnProperty(
	name utf16.Str, desc Value, throw bool,
) (bool, error) {
	if desc.Kind() != KindObject {
		if throw {
			return false, NewTypeErrorS(
				"DefineOwnProperty expects a PropertyDescriptor object",
			)
		}

		return false, nil
	}

	descobj := desc.(*Object)

	return o.DefineOwnPropertyP(name, descobj.ToPropertyDescriptor(), throw)
}

// https://es5.github.io/#x8.12.9
func (o *Object) DefineOwnPropertyP(
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

	extensible := o.Extensible
	current, ok := o.getOwnProperty(name)
	if !ok {
		if !extensible {
			return retOrThrow(NewTypeErrorS("Object %s is not extensible",
				o.Class))
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
			return retOrThrow(NewTypeErrorS("configurable is false"))
		}

		if descEnum != curEnum {
			return retOrThrow(
				NewTypeErrorS("enumerable dont match for configuration disabled"),
			)
		}
	}

	if current.IsDataDescriptor() != desc.IsDataDescriptor() {
		if !curCfg {
			return retOrThrow(NewTypeErrorS("configurable is false, cannot" +
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
					NewTypeErrorS("configurable is false and writable mismatch"),
				)
			}

			if !curWr {
				if desc.HasValue() &&
					//TODO(i4k): SameValue() ?
					!StrictEqual(current.Value(), desc.Value()) {
					return retOrThrow(NewTypeErrorS("writable is false"))
				}
			}
		}
	}

	err := copyProperties(current, desc)
	if err != nil {
		return throw(NewTypeError(err.Error()))
	}

	return throw(o.Put(name, current, throw))
}

// setOwnProperty just sets the property. Calls from ECMAScript
// must invoke DefineOwnProperty that does the correct validations.
func (o *Object) setOwnProperty(name utf16.Str, desc *PropertyDescriptor, throw bool) (bool, error) {
	retOrThrow := func(err error) (bool, error) {
		if err != nil {
			if throw {
				return false, err
			}
			return false, nil
		}

		return true, nil
	}

	if desc.IsGenericDescriptor() ||
		desc.IsDataDescriptor() {
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
	prop := o.GetProperty()
	return !StrictEqual(prop, Undefined)
}

/*
func (o *Object) Delete(name utf16.Str, failure bool) error {
	return nil
}

*/

// https://es5.github.io/#x8.12.8
func (o *Object) DefaultValue(hint Kind) (Value, error) {
	if hint == KindString ||
		hint == KindDate {
		return o.defaultString()
	}

	return o.defaultNumber()
}

func (o *Object) defaultString() (Value, error) {
	toString := o.Get(toStringAttr)
	if stringify, ok := toString.(callable); ok {
		str := stringify.Call(o, []Value{})
		if IsPrimitive(str) {
			return str, nil
		}
	}

	valueOf := o.Get(valueOfAttr)
	if valueFunc, ok := valueOf.(callable); ok {
		val := valueFunc.Call(o, []Value{})
		if IsPrimitive(val) {
			return val, nil
		}
	}

	return nil, NewTypeError("Object has no defaultValue")
}

func (o *Object) defaultNumber() (Value, error) {
	valueOf := o.Get(valueofAttr)
	if valuefunc, ok := valueOf.(callable); ok {
		val := valuefunc.Call(o, []Value{})
		if IsPrimitive(val) {
			return val, nil
		}
	}

	tostring := o.Get(tostringAttr)
	if stringify, ok := tostring.(callable); ok {
		str := stringify.Call(o, []Value{})
		if IsPrimitive(str) {
			return str
		}
	}

	return nil, NewTypeError("Object has no defaultValue")
}
