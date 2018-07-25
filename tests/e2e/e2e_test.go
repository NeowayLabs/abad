//+build e2e

// Package e2e_test has all our end to end tests that validates abad
// against google v8 engine
package e2e_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/NeowayLabs/abad/tests/fixture"
)

func TestE2E(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	testsamplesdir := filepath.Join(wd, "testdata")
	t.Parallel()
	fixture.Run(t, testsamplesdir)
}

func TestE2EConsoleLog(t *testing.T) {
	// WHY: https://stackoverflow.com/questions/22298452/v8-console-log-implementation
	// Our console.log is more rich like the browser, not d8
	fixture.RunCases(t, fixture.NewAbad(t), []fixture.TestCase{
		{
			Name: "Formatting",
			Code: `console.log("%s,%s,%s,%s,%s,%d,%i,%f", "hi",true,false,null,undefined,666,0xFF,1.1);`,
			Want: fixture.Result{
				Stdout: "hi,true,false,null,undefined,666,255,1.1\n",
			},
		},
	})
}
