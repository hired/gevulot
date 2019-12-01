package cli

import (
	"bytes"
	"io"
	"regexp"
	"strings"
	"testing"

	"github.com/hired/gevulot/pkg/server"
	"github.com/stretchr/testify/assert"
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
		runServer: func(_ server.Logger, _ <-chan *server.Config) error { return nil },
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
		expectedStdout   string
		expectedStderr   string
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
			assert.Regexp(t, tc.expectedStdout, mockedStdout.String())
			assert.Regexp(t, tc.expectedStderr, mockedStderr.String())
		})
	}

	t.Run("starts the server", func(t *testing.T) {
		defer cleanup()

		handlerCalled := false

		currentContext.runServer = func(_ server.Logger, _ <-chan *server.Config) error {
			handlerCalled = true
			return nil
		}

		Run([]string{"--config=testdata/example.toml"})

		assert.Equal(t, handlerCalled, true, "Run starts the server")
	})

	t.Run("exit code when there is a server error", func(t *testing.T) {
		defer cleanup()

		currentContext.runServer = func(_ server.Logger, _ <-chan *server.Config) error {
			return io.EOF
		}

		exitCode := Run([]string{"--config=testdata/example.toml"})

		assert.Equal(t, mockedStderr.String(), "server error: EOF\n", "Run prints server error")
		assert.Equal(t, exitCode, 1, "Run returns non-zero exit code")
	})

	t.Run("exit code when server exited without error", func(t *testing.T) {
		defer cleanup()

		currentContext.runServer = func(_ server.Logger, _ <-chan *server.Config) error {
			return nil
		}

		exitCode := Run([]string{"--config=testdata/example.toml"})

		assert.Equal(t, exitCode, 0, "Run returns zero exit code")
	})

	t.Run("provides logger instance", func(t *testing.T) {
		defer cleanup()

		currentContext.runServer = func(log server.Logger, _ <-chan *server.Config) error {
			log.Info("logger is working!")
			return nil
		}

		Run([]string{"--config=testdata/example.toml"})

		assert.Contains(t, mockedStdout.String(), "logger is working!", "Run prints log message")
	})
}
