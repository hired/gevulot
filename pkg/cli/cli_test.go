package cli

import (
	"bytes"
	"errors"
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
	// Mock execution context
	mockContext()

	// Set build info
	version = "1.0"
	commitHash = "deadbeef"
	buildDate = "10/29/1987"
}

func mockContext() {
	currentContext = &cliContext{
		stdout:    mockedStdout,
		stderr:    mockedStderr,
		runServer: func() error { return nil },
	}
}

func cleanup() {
	mockedStdout.Reset()
	mockedStderr.Reset()
	mockContext()
}

// NB: THIS MUST BE SEQUENTIAL!
func TestCliRun(t *testing.T) {
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

	for _, tc := range testCases {
		args := strings.Split(tc.input, " ")

		t.Run(tc.input, func(t *testing.T) {
			defer cleanup()

			err := Run(args)

			assert.NilError(t, err)
			assert.Assert(t, cmp.Regexp(tc.expectedStdout, mockedStdout.String()))
			assert.Assert(t, cmp.Regexp(tc.expectedStderr, mockedStderr.String()))
		})
	}

	t.Run("run server", func(t *testing.T) {
		defer cleanup()

		handlerCalled := false
		serverError := errors.New("server error")

		currentContext.runServer = func() error {
			handlerCalled = true
			return serverError
		}

		err := Run(nil)

		assert.Equal(t, true, handlerCalled, "Run starts the server")
		assert.Equal(t, err, serverError, "Run propogates server error")
	})
}
