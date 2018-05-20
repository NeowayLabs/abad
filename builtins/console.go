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

	logfn := types.NewBuiltinfn(log)
	err := console.Put(logAttr, logfn, true)
	if err != nil {
		return nil, err
	}

	toStrfn := types.NewBuiltinfn(toString)
	err = console.Put(toStringAttr, toStrfn, true)
	return console, nil
}

func log(_ types.Object, args []types.Value) types.Value {
	for _, v := range args {
		fmt.Printf("%s\n", v.ToString())
	}
	return types.Undefined
}

func toString(_ types.Object, args []types.Value) types.Value {
	return types.NewString("[object Object]")
}