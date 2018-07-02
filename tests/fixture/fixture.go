// Package fixture contains all the test fixture required to
// run end to end tests comparing abad against google v8.
// It is used internally by abad's tests but it can be useful for
// clients to write their own end to end tests for more specific scenarios
package fixture

import (
	"testing"
)


// Run will run tests using the provided samplesdir as a 
// source of JavaScript samples. It will compare the results
// of running the code on abad with Google's V8 engine.
//
// This function assumes that both executables are installed and
// working on the system.
//
// The samplesdir will be traversed recursively and for each dir it
// will create a new hierarchy of subtests and the name of the dir is
// used as the name of the test (the filename of the sample is also used).
func Run(t *testing.T, samplesdir string) {
}


type Result struct {
	stdout string
	stderr string
}

type JsInterpreter func(codepath string) (error, Result)