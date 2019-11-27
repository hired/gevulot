package cli

import (
	"bytes"
	"regexp"
	"strings"
	"testing"

	"gotest.tools/assert"
	"gotest.tools/assert/cmp"
)

var (
	mockedStdout = &bytes.Buffer{}
	mockedStderr = &bytes.Buffer{}
)

func init() {
	// Override STDOUT/STDERR
	stdout = mockedStdout
	stderr = mockedStderr

	// Set build info
	version = "1.0"
	commitHash = "deadbeef"
	buildDate = "10/29/1987"
}

func cleanup() {
	mockedStdout.Reset()
	mockedStderr.Reset()
}

func TestRun(t *testing.T) {
	defer cleanup()

	// Regexp that represents an empty output
	none := "^$"

	testCases := []struct {
		input          string
		expectedStdout cmp.RegexOrPattern
		expectedStderr cmp.RegexOrPattern
	}{
		{"--version", none, regexp.QuoteMeta("gevulot version 1.0 (deadbeef) built on 10/29/1987")},
		{"--help", none, regexp.QuoteMeta("usage: gevulot")},
	}

	// NB: THIS MUST BE SEQUENTIAL!
	for _, tc := range testCases {
		err := Run(strings.Split(tc.input, " "))

		assert.NilError(t, err)
		assert.Assert(t, cmp.Regexp(tc.expectedStdout, mockedStdout.String()))
		assert.Assert(t, cmp.Regexp(tc.expectedStderr, mockedStderr.String()))

		cleanup()
	}
}
