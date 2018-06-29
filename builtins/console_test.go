package builtins_test

import (
	"testing"

	"github.com/NeowayLabs/abad/builtins"
	"github.com/NeowayLabs/abad/types"
	"github.com/madlambda/spells/assert"
	"github.com/NeowayLabs/abad/internal/utf16"
)

func TestConsoleToString(t *testing.T) {
	console, err := builtins.NewConsole()
	assert.NoError(t, err, "console creation")
	assert.EqualStrings(t, console.String(),
		"[object Object]", "console toString")

	log, err := console.Get(utf16.S("log"))
	assert.NoError(t, err, "console get log")

	logfn, ok := log.(*types.Builtinfn)
	if !ok {
		t.Fatalf("log is not a function")
	}

	call(logfn)
}


func call(fn types.Function) {
	fn.Call(nil, nil)
}