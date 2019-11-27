package cli

import (
	"bytes"
	"io"
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
		input            string
		expectedStdout   cmp.RegexOrPattern
		expectedStderr   cmp.RegexOrPattern
		expectedExitCode int
	}{
		{"--version", none, regexp.QuoteMeta("gevulot version 1.0 (deadbeef) built on 10/29/1987"), 0},
		{"--help", none, regexp.QuoteMeta("usage: gevulot"), 0},
	}

	for _, tc := range testCases {
		args := strings.Split(tc.input, " ")

		t.Run(tc.input, func(t *testing.T) {
			defer cleanup()

			exitCode := Run(args)

			assert.Equal(t, exitCode, tc.expectedExitCode)
			assert.Assert(t, cmp.Regexp(tc.expectedStdout, mockedStdout.String()))
			assert.Assert(t, cmp.Regexp(tc.expectedStderr, mockedStderr.String()))
		})
	}

	t.Run("starts the server", func(t *testing.T) {
		defer cleanup()

		handlerCalled := false

		currentContext.runServer = func() error {
			handlerCalled = true
			return nil
		}

		Run(nil)

		assert.Equal(t, handlerCalled, true, "Run starts the server")
	})

	t.Run("exit code when there is a server error", func(t *testing.T) {
		defer cleanup()

		currentContext.runServer = func() error {
			return io.EOF
		}

		exitCode := Run(nil)

		assert.Equal(t, mockedStderr.String(), "server error: EOF\n", "Run prints server error")
		assert.Equal(t, exitCode, 1, "Run returns non-zero exit code")
	})

	t.Run("exit code when server exited without error", func(t *testing.T) {
		defer cleanup()

		currentContext.runServer = func() error {
			return nil
		}

		exitCode := Run(nil)

		assert.Equal(t, exitCode, 0, "Run returns zero exit code")
	})
}
