package main

import (
	"os"

	"github.com/hired/gevulot/pkg/cli"
)

// main is the entrypoint of Gevulot.
func main() {
	// All logic lives in the cli package
	err := cli.Run(os.Args[1:])

	if err != nil {
		// NB: this is THE ONLY PLACE where we exit from the program abnormaly
		// DO NOT USE os.Exit!
		os.Exit(1)
	}
}
