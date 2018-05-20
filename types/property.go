package types

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
	// The code above is the same as:
	//   Object.defineOwnProperty(this.name, {
	//       value: "value",
	//       writable: true,
	//       enumerable: true,
	//       configurable: true
	//   });
	PropertyDescriptor struct {
		attrs map[string]Value
	}
)

// Default attribute values
// Table 7 in https://es5.github.io/#x8.6.1
var (
	DefValue        = Undefined
	DefWritable     = False
	DefGet          = Undefined
	DefSet          = Undefined
	DefEnumerable   = False
	DefConfigurable = False
)

// NewGenericPropDesc creates a new generic (empty) property
// descriptor.
// https://es5.github.io/#x8.10.3
func NewGenericPropDesc() *PropertyDescriptor {
	return &PropertyDescriptor{
		attrs: make(map[string]Value),
	}
}

// NewDataPropDesc creates a new Data Property Descriptor.
// https://es5.github.io/#x8.10.2
func NewDataPropDesc(value Value, wrt, enum, cfg bool) *PropertyDescriptor {
	p := NewGenericPropDesc()

	p.put("value", value)
	p.put("writable", Bool(wrt))
	p.put("enumerable", Bool(enum))
	p.put("configurable", Bool(cfg))

	return p
}

// NewAcessorPropDesc creates a new Acessor Property Descriptor.
// https://es5.github.io/#x8.10.1
func NewAcessorPropDesc(get, set Value, enum, cfg bool) *PropertyDescriptor {
	p := NewGenericPropDesc()
	p.put("get", get)
	p.put("set", set)
	p.put("enumerable", Bool(enum))
	p.put("configurable", Bool(cfg))
	return p
}

// https://es5.github.io/#x8.6.1
func DefaultDataPropDesc() *PropertyDescriptor {
	return NewDataPropDesc(
		Undefined, false, false, false,
	)
}

// https://es5.github.io/#x8.6.1
func DefaultAcessorPropDesc() *PropertyDescriptor {
	return NewAcessorPropDesc(
		Undefined, Undefined, false, false,
	)
}

func (p *PropertyDescriptor) put(name string, value Value) {
	p.attrs[name] = value
}

func (p *PropertyDescriptor) get(name string) Value {
	if p == nil || p.attrs == nil {
		panic(p)
	}

	v, ok := p.attrs[name]
	if ok {
		return v
	}

	return Undefined
}

// has checks if name is a existent property
func (p *PropertyDescriptor) has(name string) bool {
	_, ok := p.attrs[name]
	return ok
}

func (p *PropertyDescriptor) HasValue() bool {
	return p.has("value")
}

func (p *PropertyDescriptor) HasWritable() bool {
	return p.has("writable")
}

func (p *PropertyDescriptor) HasGet() bool {
	return p.has("get")
}
func (p *PropertyDescriptor) HasSet() bool {
	return p.has("set")
}

func (p *PropertyDescriptor) HasCfg() bool {
	return p.has("configurable")
}

func (p *PropertyDescriptor) HasEnum() bool {
	return p.has("enumerable")
}

func (p *PropertyDescriptor) SetValue(v Value) {
	p.put("value", v)
}

func (p *PropertyDescriptor) SetWritable(b Bool) {
	p.put("writable", b)
}

func (p *PropertyDescriptor) SetGet(get Value) {
	p.put("get", get)
}

func (p *PropertyDescriptor) SetSet(set Value) {
	p.put("set", set)
}

func (p *PropertyDescriptor) SetEnum(b Bool) {
	p.put("enumerable", b)
}

func (p *PropertyDescriptor) SetCfg(b Bool) {
	p.put("configurable", b)
}

// Value returns the [[value]] property
func (p *PropertyDescriptor) Value() Value {
	return p.get("value")
}

// Writable returns the [[writable]] property
func (p *PropertyDescriptor) Writable() Value {
	return p.get("writable")
}

// Get returns the [[get]] property
func (p *PropertyDescriptor) Get() Value {
	return p.get("get")
}

// Set returns the [[set]] property
func (p *PropertyDescriptor) Set() Value {
	return p.get("set")
}

// Enum returns the [[enumerable]] property
func (p *PropertyDescriptor) Enum() Value {
	return p.get("enumerable")
}

// Cfg returns the [[configurable]] property
func (p *PropertyDescriptor) Cfg() Value {
	return p.get("configurable")
}

func (p *PropertyDescriptor) IsGenericDescriptor() bool {
	return !p.IsDataDescriptor() && !p.IsAcessorDescriptor()
}

func (p *PropertyDescriptor) IsDataDescriptor() bool {
	return p.has("value") || p.has("writable")
}

func (p *PropertyDescriptor) IsAcessorDescriptor() bool {
	return p.has("get") || p.has("set")
}

func (p *PropertyDescriptor) IsAbsentDescriptor() bool {
	return !p.IsDataDescriptor() && !p.IsAcessorDescriptor() &&
		!p.HasEnum() && !p.HasCfg()
}

func (p *PropertyDescriptor) ToObject() *DataObject {
	obj := NewDataObject(Null)

	if p.IsDataDescriptor() {
		if p.has("value") {
			obj.DefineOwnPropertyP(valueAttr, NewDataPropDesc(
				p.Value(), true, true, true,
			), false)
		}

		if p.has("writable") {
			obj.DefineOwnPropertyP(writableAttr, NewDataPropDesc(
				p.Writable(), true, true, true,
			), false)
		}
	} else if p.IsAcessorDescriptor() {
		if p.has("get") {
			obj.DefineOwnPropertyP(getAttr, NewDataPropDesc(
				p.Get(), true, true, true,
			), false)
		}

		if p.has("set") {
			obj.DefineOwnPropertyP(setAttr, NewDataPropDesc(
				p.Set(), true, true, true,
			), false)
		}
	}

	if p.has("enumerable") {
		obj.DefineOwnPropertyP(enumAttr, NewDataPropDesc(
			p.Enum(), true, true, true,
		), false)
	}

	if p.has("configurable") {
		obj.DefineOwnPropertyP(cfgAttr, NewDataPropDesc(
			p.Cfg(), true, true, true,
		), false)
	}

	return obj
}

func IsSameDescriptor(a, b *PropertyDescriptor) bool {
	var ok = true

	if a.IsDataDescriptor() {
		ok = ok && StrictEqual(a.Value(), b.Value()) &&
			StrictEqual(a.Writable(), b.Writable())
	} else if a.IsAcessorDescriptor() {
		ok = ok && StrictEqual(a.Get(), b.Get()) &&
			StrictEqual(a.Set(), b.Set())
	}

	ok = ok && StrictEqual(a.Enum(), b.Enum()) &&
		StrictEqual(a.Cfg(), b.Cfg())

	return ok
}

func CopyProperties(dst, src *PropertyDescriptor) {
	if src.IsDataDescriptor() {
		if src.HasValue() {
			dst.SetValue(src.Value())
		}

		if src.HasWritable() {
			dst.SetWritable(src.Writable().(Bool))
		}
	} else if src.IsAcessorDescriptor() {
		if src.HasGet() {
			dst.SetGet(src.Get())
		}

		if src.HasSet() {
			dst.SetSet(src.Set())
		}
	}

	if src.HasEnum() {
		dst.SetEnum(src.Enum().(Bool))
	}

	if src.HasCfg() {
		dst.SetCfg(src.Cfg().(Bool))
	}
}