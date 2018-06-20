package cli_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/NeowayLabs/abad/cmd/abad/cli"
	"github.com/madlambda/spells/assert"
)

var trim = strings.TrimSpace

func TestCli(t *testing.T) {
	for _, tc := range []struct {
		in  string
		out string
	}{
		{
			in:  "0",
			out: "0",
		},
		{
			in:  "0xff",
			out: "255",
		},
		{
			in:  "1.0e1",
			out: "10",
		},
		{
			in:  "1e-10",
			out: "0.0000000001",
		},
		{
			in:  "1e10",
			out: "10000000000",
		},
	} {
		var inb bytes.Buffer
		var outb bytes.Buffer
		cli := cli.NewCli("test.js", &inb, &outb)

		_, err := inb.Write([]byte(tc.in + "\n"))
		assert.NoError(t, err)
		cli.ReadEval()

		got := trim(outb.String())
		expected := "> < " + trim(tc.out)
		assert.EqualStrings(t, expected, got, "cli output")
	}
}