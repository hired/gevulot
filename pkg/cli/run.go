package cli

import (
	"os"

	"github.com/hired/gevulot/pkg/server"
)

// Run handles CLI for Gevulot server and returns exit code.
// This is the only publicly exposed entry point for the cli package.
func Run(args []string) int {
	// Initialize new instance with default STDERR/STDOUT
	cli := &cli{
		stdout:    os.Stdout,
		stderr:    os.Stderr,
		runServer: server.Run, // late binding to improve testability
	}

	return cli.Run(args)
}
