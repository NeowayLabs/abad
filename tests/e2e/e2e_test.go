//+build e2e

// Package e2e_test has all our end to end tests that validates abad
// against google v8 engine
package e2e_test

import (
	"os"
	"testing"
	"path/filepath"
	
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