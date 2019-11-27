package cli

import (
	"fmt"
	"io"
	"os"

	"gopkg.in/alecthomas/kingpin.v2"
)

// Application name to display in help
const appName = "gevulot"

// Build information; provisioned by ldflags
var (
	version    string = "<unknown>"
	commitHash string = "<unknown>"
	buildDate  string = "<unknown>"
)

// Default STDOUT/STDERR for console messages; we override this in tests
var stdout io.Writer = os.Stdout
var stderr io.Writer = os.Stderr

// cliArgs contains user provided arguments and flags.
type cliArgs struct {
	// Path to the config
	configPath string
}

// parseArgs parses CLI arguments.
func parseArgs(args []string) (*cliArgs, error) {
	parsedArgs := &cliArgs{}

	app := kingpin.New(appName, "")

	// Do not call os.Exit
	app.Terminate(nil)

	// Write output to stderr
	app.Writer(stderr)

	// Add --version flag with to display build info
	app.Version(fmt.Sprintf("%s version %s (%s) built on %s", appName, version, commitHash, buildDate))

	app.Flag("config", "Set the configuration file path").
		Short('c').
		PlaceHolder("PATH").
		StringVar(&parsedArgs.configPath)

	_, err := app.Parse(args)

	if err != nil {
		return nil, err
	}

	return parsedArgs, nil
}

// Run executes gevulot using given CLI args.
func Run(args []string) error {
	_, err := parseArgs(args)

	if err != nil {
		return err
	}

	return nil
}
