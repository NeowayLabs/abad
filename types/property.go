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
const (
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
	p.put("writable", wrt)
	p.put("enumerable", enum)
	p.put("configurable", cfg)

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
	return o
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

func (p *PropertyDescriptor) HasCfg() bool {
	return p.has("configurable")
}

func (p *PropertyDescriptor) HasEnum() bool {
	return p.has("enumerable")
}

func (p *PropertyDescriptor) SetValue(v Value) {
	p.put("value", v)
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

func (p *PropertyDescriptor) ToObject() *Object {
	obj := NewObject(Null)

	if p.IsDataDescriptor() {
		if p.has("value") {
			obj.DefineOwnPropertyP("value", NewDataPropDesc(
				p.Value(), true, true, true,
			), false)
		}

		if p.has("writable") {
			obj.DefineOwnPropertyP("writable", NewDataPropDesc(
				p.Writable(), true, true, true,
			), false)
		}
	} else if p.IsAcessorDescriptor() {
		if p.has("get") {
			obj.DefineOwnPropertyP("get", NewDataPropDesc(
				p.Get(), true, true, true,
			), false)
		}

		if p.has("set") {
			obj.DefineOwnPropertyP("set", NewDataPropDesc(
				p.Set(), true, true, true,
			), false)
		}
	}

	if p.has("enumerable") {
		obj.DefineOwnProperty("enumerable", NewDataPropDesc(
			p.Enum(), true, true, true,
		), false)
	}

	if p.has("configurable") {
		obj.DefineOwnProperty("configurable", NewDataPropDesc(
			p.Cfg(), true, true, true,
		), false)
	}

	return obj
}

func IsSameDescriptor(a, b *PropertyDescriptor) bool {
	return a.Value().Equal(b.Value()) &&
		a.Writable().Equal(b.Writable()) &&
		a.Enum().Equal(b.Enum()) &&
		a.Cfg().Equal(b.Cfg()) &&
		a.Get().Equal(b.Get()) &&
		a.Set().Equal(b.Set())
}

func CopyProperties(dst, src *PropertyDescriptor) {
	if src.IsDataDescriptor() {
		if src.HasValue() {
			dst.SetValue(src.Value())
		}

		if src.HasWritable() {
			dst.SetWritable(src.Writable())
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
		dst.SetEnum(src.Enum())
	}

	if src.HasCfg() {
		dst.SetCfg(src.Cfg())
	}
}