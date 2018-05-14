package types

import "github.com/NeowayLabs/abad/envrec"

func NewFunctionProto() *Object {
	proto := NewObject(Null)
	proto.Extensible = true
	proto.Length = 0
	return proto
}

func NewFunction(
	params []utf16.S, body *ast.Body, scope *envrec.Decl, strict bool,
) *Object {
	o := NewObject(NewFunctionProto())
	o.Class = "Function"
	o.Scope = scope
	o.Params = params
	o.Code = body
	o.Extensible = true
}