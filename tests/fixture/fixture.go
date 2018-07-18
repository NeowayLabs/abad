// Package fixture contains all the test fixture required to
// run end to end tests comparing abad against google v8.
// It is used internally by abad's tests but it can be useful for
// clients to write their own end to end tests for more specific scenarios
package fixture

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/madlambda/spells/assert"
)

type Result struct {
	Stdout string
	Stderr string
}

type JsInterpreter func(codepath string) (error, Result)

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
	samplesdir, err := filepath.Abs(samplesdir)
	if err != nil {
		t.Fatal(err)
	}
	abadInterpreter := NewAbad(t)
	v8Interpreter := NewV8(t)

	RunWithInterpreters(t, samplesdir, v8Interpreter, abadInterpreter)
}

func NewAbad(t *testing.T) JsInterpreter {
	return newInterpreter(t, "abad")
}

func NewV8(t *testing.T) JsInterpreter {
	return newInterpreter(t, "d8")
}

func RunWithInterpreters(
	t *testing.T,
	samplesdir string,
	reference JsInterpreter,
	undertest JsInterpreter,
) {

	err := filepath.Walk(samplesdir, func(path string, info os.FileInfo, err error) error {
		assert.NoError(t, err, "error[%s] walking code samples dir[%s], path[%s]", err, samplesdir, path)

		if info.IsDir() {
			return nil
		}

		testname := strings.TrimPrefix(path, samplesdir)[1:]
		t.Run(testname, func(t *testing.T) {
			err, want := reference(path)
			assertSuccessRun(t, want, err)

			err, got := undertest(path)
			assertSuccessRun(t, got, err)

			assertEqualOutput(t, "stdout", want.Stdout, got.Stdout)
			assertEqualOutput(t, "stderr", want.Stderr, got.Stderr)
		})

		return nil
	})

	assert.NoError(t, err)
}

func assertSuccessRun(t *testing.T, r Result, err error) {
	t.Helper()

	if err == nil {
		return
	}

	t.Fatalf("\nfatal error:[%s]\n\nstdout:\n%s\nstderr:\n%s\n", err, r.Stdout, r.Stderr)
}

func assertEqualOutput(t *testing.T, outputname string, want string, got string) {
	t.Helper()

	wantedlines := strings.Split(want, "\n")
	gotlines := strings.Split(got, "\n")

	if len(wantedlines) != len(gotlines) {
		t.Errorf(
			"%s: wanted output has [%d] lines but got output has [%d]",
			outputname,
			len(wantedlines),
			len(gotlines),
		)
		t.Fatalf("%s: wanted[%s] != got[%s]", outputname, want, got)
	}

	for line, wantedline := range wantedlines {
		gotline := gotlines[line]
		if wantedline != gotline {
			t.Errorf(
				"%s: line[%d] differ: want[%s] != got[%s]",
				outputname,
				line,
				wantedline,
				gotline,
			)
		}
	}
}

func newInterpreter(t *testing.T, jsinterpreter string) JsInterpreter {
	a := exec.Command(jsinterpreter, "-help")
	err := a.Run()
	if err != nil {
		t.Fatalf(
			"unable to find the interpreter[%s] installed, got error: %s",
			jsinterpreter, err)
	}
	return func(codepath string) (error, Result) {
		js := exec.Command(jsinterpreter, codepath)
		stdout := &bytes.Buffer{}
		stderr := &bytes.Buffer{}

		js.Stdout = stdout
		js.Stderr = stderr

		err := js.Run()

		return err, Result{Stdout: stdout.String(), Stderr: stderr.String()}
	}
}
