package types

import (
	"github.com/NeowayLabs/abad/ast"
)

func NewFunctionProto() (*Object, error) {
	proto := NewObject(Null)
	return proto, nil
}

func NewFunction(
	params []ast.Ident, body *ast.Program, scope *Object, strict bool,
) (*Object, error) {
	proto, err := NewFunctionProto()
	if err != nil {
		return nil, err
	}

	o := NewObject(proto)
	o.Class = "Function"
	o.Scope = scope
	o.Params = params
	o.Code = body
	return o, nil
}