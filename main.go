package main

import (
	"os"

	"github.com/hired/gevulot/pkg/cli"
)

// main is the entrypoint of Gevulot.
func main() {
	// All logic lives in the cli package
	exitCode := cli.Run(os.Args[1:])

	// NB: this is THE ONLY PLACE where we exit from the program
	// DO NOT USE os.Exit ANYWHERE ELSE!
	os.Exit(exitCode)
}
