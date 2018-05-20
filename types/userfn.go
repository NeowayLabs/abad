package types

import (
	"github.com/NeowayLabs/abad/ast"
	"github.com/NeowayLabs/abad/internal/utf16"
)

type (
	UserFunction struct {
		*DataObject

		isFnPrototype bool

		params []utf16.Str
		body   *ast.FnBody
		scope  interface{}
	}
)

func NewUserFunctionPrototype() *UserFunction {
	return &UserFunction{
		isFnPrototype: true,
		DataObject:    NewBaseDataObject(),
	}
}

func NewUserFunction(
	params []utf16.Str, body *ast.FnBody, scope interface{}, strict bool,
) *UserFunction {
	return &UserFunction{
		params: params,
		body:   body,
		scope:  scope,
		DataObject: NewDataObject(NewUserFunctionPrototype()),
	}
}

func (f *UserFunction) Call(this *Object, params []Value) Value {
	if f.isFnPrototype {
		return Undefined
	}

	return True // todo
}