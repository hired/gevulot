package cli

import (
	"fmt"
	"io"
	"os"

	"github.com/hired/gevulot/pkg/server"
	"github.com/sirupsen/logrus"
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

// cliArgs contains user provided arguments and flags.
type cliArgs struct {
	// True when user asked for the CLI help
	isHelp bool

	// True when user asked for the app version
	isVersion bool

	// Path to the config
	configPath string
}

// cliContext contains the global state for CLI.
type cliContext struct {
	// Standard output stream for console messages.
	stdout io.Writer

	// Standard error stream for console messages.
	stderr io.Writer

	// runServer starts the Gevulot.
	runServer func(log server.Logger, configChan <-chan server.Config) error
}

// Default execution context.
var defaultContext = &cliContext{
	stdout:    os.Stdout,
	stderr:    os.Stderr,
	runServer: server.Run,
}

// Current execution context; we override this in tests.
var currentContext = defaultContext

// parseArgs parses CLI arguments.
func parseArgs(args []string) (*cliArgs, error) {
	parsedArgs := &cliArgs{}

	app := kingpin.New(appName, "")

	// Do not call os.Exit
	app.Terminate(nil)

	// Write output to stderr
	app.Writer(currentContext.stderr)

	// Add --version flag with to display build info
	app.Version(fmt.Sprintf("%s version %s (%s) built on %s", appName, version, commitHash, buildDate))

	app.Flag("config", "Set the configuration file path").
		Short('c').
		PlaceHolder("PATH").
		Default("gevulot.toml").
		StringVar(&parsedArgs.configPath)

	// Expose --help and --version flags to our struct
	app.HelpFlag.BoolVar(&parsedArgs.isHelp)
	app.VersionFlag.BoolVar(&parsedArgs.isHelp)

	_, err := app.Parse(args)

	if err != nil {
		return nil, err
	}

	return parsedArgs, nil
}

func configureLogger() *logrus.Logger {
	logger := logrus.New()
	logger.Out = currentContext.stdout

	return logger
}

// Run executes gevulot using given CLI args. The function returns program exit code.
func Run(args []string) (exitCode int) {
	flags, err := parseArgs(args)

	if err != nil {
		fmt.Fprintf(currentContext.stderr, "%v\n", err)
		exitCode = 1
	}

	if flags.isHelp || flags.isVersion {
		exitCode = 0
		return
	}

	log := configureLogger()

	configChan := make(chan server.Config, 1)
	err = currentContext.runServer(log, configChan)

	if err != nil {
		fmt.Fprintf(currentContext.stderr, "server error: %v\n", err)
		exitCode = 1
		return
	}

	return
}
