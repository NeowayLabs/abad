package types

type (
	Execfn    func(this Object, args []Value) Value
	Builtinfn struct {
		*UserFunction

		fn Execfn
	}
)

func NewBuiltinfn(fn Execfn) *Builtinfn {
	return &Builtinfn{
		fn: fn,

		UserFunction: &UserFunction{
			DataObject: NewDataObject(NewUserFunctionPrototype()),
		},
	}
}

func (f *Builtinfn) Call(this Object, args []Value) Value {
	return f.fn(this, args)
}