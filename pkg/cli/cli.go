package cli

import (
	"fmt"
	"io"

	"github.com/hired/gevulot/pkg/server"
	log "github.com/sirupsen/logrus"
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

type cli struct {
	// Standard output stream for console messages.
	stdout io.Writer

	// Standard error stream for console messages.
	stderr io.Writer

	// runServer starts the Gevulot server.
	runServer func(configChan <-chan *server.Config) error
}

// cliArgs contains user provided arguments and flags.
type cliArgs struct {
	// True when user asked for the CLI help
	isHelp bool

	// True when user asked for the app version
	isVersion bool

	// Path to the config
	configPath string
}

// parseArgs parses CLI arguments.
func (c *cli) parseArgs(args []string) (*cliArgs, error) {
	parsedArgs := &cliArgs{}

	app := kingpin.New(appName, "")

	// Do not call os.Exit
	app.Terminate(nil)

	// Write output to stderr
	app.Writer(c.stderr)

	// Add --version flag with to display build info
	app.Version(fmt.Sprintf("%s version %s (%s) built on %s", appName, version, commitHash, buildDate))

	// Add --config flag to specify path to the config
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

func (c *cli) configureLogger() {
	log.SetOutput(c.stdout)
}

func (c *cli) prepareConfigChan(configPath string) (<-chan *server.Config, error) {
	config, err := readServerConfig(configPath)

	if err != nil {
		return nil, err
	}

	configChan := make(chan *server.Config, 1)
	configChan <- config

	watcher := newFileWatcher(configPath)
	watcher.OnWrite = func() {
		updatedConfig, err := readServerConfig(configPath)

		if err != nil {
			log.Errorf("error loading config file %s: %v", configPath, err)
			return
		}

		configChan <- updatedConfig
	}

	err = watcher.Watch()

	if err != nil {
		return nil, err
	}

	return configChan, nil
}

// Run handles CLI for Gevulot server and returns exit code.
func (c *cli) Run(args []string) (exitCode int) {
	// Parse CLI args and flags
	flags, err := c.parseArgs(args)

	if err != nil {
		fmt.Fprintf(c.stderr, "%v\n", err)
		exitCode = 1
	}

	// Exit immediately if user asked for a help or an app version
	if flags.isHelp || flags.isVersion {
		exitCode = 0
		return
	}

	// Setup logrus
	c.configureLogger()

	// Load config
	configChan, err := c.prepareConfigChan(flags.configPath)

	if err != nil {
		fmt.Fprintf(c.stderr, "failed to load config: %v\n", err)
		exitCode = 1
		return
	}

	// Run the server (this is blocking call)
	err = c.runServer(configChan)

	if err != nil {
		fmt.Fprintf(c.stderr, "server error: %v\n", err)
		exitCode = 1
		return
	}

	return
}
