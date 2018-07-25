package builtins

import (
	"fmt"
	"strings"

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

	// TODO: handle error
	console.Put(toStringAttr, toStrfn, true)
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
	// This will not handle errors in formatting properly
	// But it will work for well formatted messages
	if len(args) == 0 {
		fmt.Println("")
		return types.Undefined
	}
	
	vals := []string{}
	for _, v := range args {
		vals = append(vals, v.ToString().String())
	}
	msg := ""
	if hasFormatting(vals[0]) {
		msg = sprintf(vals)
	} else {
		msg = strings.Join(vals, " ")
	}
	fmt.Println(msg)
	return types.Undefined
}

func sprintf(vals []string) string {
	msg := vals[0]
	vals = vals[1:]
	replace := func(fmtReplacer string) {
		msg = strings.Replace(msg, fmtReplacer, "%s", -1)
	}
	
	for _, fmtReplacer := range []string{ "%d", "%i", "%f", "%o", "%O" } {
		replace(fmtReplacer)
	}
	
	args := make([]interface{}, len(vals))
	for i, v := range vals {
		args[i] = v
	}
	
	return fmt.Sprintf(msg, args...)
}

func hasFormatting(a string) bool {
	return strings.Contains(a, "%")
}

func toStringer(str string) types.Execfn {
	return func(_ types.Object, args []types.Value) types.Value {
		return types.NewString(str)
	}
}
