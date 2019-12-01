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

func init() {
	// Set build info
	version = "1.0"
	commitHash = "deadbeef"
	buildDate = "10/29/1987"
}

func mockedCli(stdout, stderr io.Writer) *cli {
	return &cli{
		stdout:    stdout,
		stderr:    stderr,
		runServer: func(_ <-chan *server.Config) error { return nil },
	}
}

// NB: THIS MUST BE SEQUENTIAL!
func TestCliRun(t *testing.T) {
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
			mockedStdout := &bytes.Buffer{}
			mockedStderr := &bytes.Buffer{}

			cli := mockedCli(mockedStdout, mockedStderr)
			exitCode := cli.Run(args)

			assert.Equal(t, exitCode, tc.expectedExitCode)
			assert.Regexp(t, tc.expectedStdout, mockedStdout.String())
			assert.Regexp(t, tc.expectedStderr, mockedStderr.String())
		})
	}

	t.Run("starts the server", func(t *testing.T) {
		handlerCalled := false

		cli := mockedCli(nil, nil)
		cli.runServer = func(_ <-chan *server.Config) error {
			handlerCalled = true
			return nil
		}

		cli.Run([]string{"--config=testdata/example.toml"})

		assert.Equal(t, handlerCalled, true, "Run starts the server")
	})

	t.Run("exit code when there is a server error", func(t *testing.T) {
		mockedStderr := &bytes.Buffer{}

		cli := mockedCli(nil, mockedStderr)
		cli.runServer = func(_ <-chan *server.Config) error {
			return io.EOF
		}

		exitCode := cli.Run([]string{"--config=testdata/example.toml"})

		assert.Equal(t, mockedStderr.String(), "server error: EOF\n", "Run prints server error")
		assert.Equal(t, exitCode, 1, "Run returns non-zero exit code")
	})

	t.Run("exit code when server exited without error", func(t *testing.T) {
		cli := mockedCli(nil, nil)
		cli.runServer = func(_ <-chan *server.Config) error {
			return nil
		}

		exitCode := cli.Run([]string{"--config=testdata/example.toml"})

		assert.Equal(t, exitCode, 0, "Run returns zero exit code")
	})
}
