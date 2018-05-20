package builtins_test

import (
	"testing"

	"github.com/NeowayLabs/abad/builtins"
	"github.com/madlambda/spells/assert"
)

func TestConsoleToString(t *testing.T) {
	console, err := builtins.NewConsole()
	assert.NoError(t, err, "console creation")
	assert.EqualStrings(t, console.String(),
		"[object Object]", "console toString")
}