package builtins

import (
	"fmt"

	"github.com/NeowayLabs/abad/internal/utf16"
	"github.com/NeowayLabs/abad/types"
)

type (
	Console struct {
		*types.DataObject
	}
)

var (
	logAttr      = utf16.S("log")
	toStringAttr = utf16.S("toString")
)

func NewConsole() (*Console, error) {
	console := &Console{
		DataObject: types.NewBaseDataObject(),
	}

	logfn, err := newlog()
	if err != nil {
		return nil, err
	}

	err = console.Put(logAttr, logfn, true)
	if err != nil {
		return nil, err
	}

	toStrfn := types.NewBuiltinfn(
		toStringer("[object Object]"),
	)

	err = console.Put(toStringAttr, toStrfn, true)
	return console, nil
}

func newlog() (*types.Builtinfn, error) {
	logfn := types.NewBuiltinfn(log)
	toStrfn := types.NewBuiltinfn(
		toStringer("function () { [native code] }"),
	)
	err := logfn.Put(toStringAttr, toStrfn, true)
	return logfn, err
}

func log(_ types.Object, args []types.Value) types.Value {
	for _, v := range args {
		fmt.Printf("%s\n", v.ToString())
	}
	return types.Undefined
}

func toStringer(str string) types.Execfn {
	return func(_ types.Object, args []types.Value) types.Value {
		return types.NewString(str)
	}
}