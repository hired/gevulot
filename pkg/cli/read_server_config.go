package cli

import (
	"path/filepath"

	"github.com/BurntSushi/toml"

	"github.com/hired/gevulot/pkg/server"
)

// readServerConfig unmarshals server config at the given file path.
func readServerConfig(path string) (*server.Config, error) {
	// Convert to absolute path
	absPath, err := filepath.Abs(path)

	if err != nil {
		return nil, err
	}

	// We use TOML
	config := &server.Config{}
	_, err = toml.DecodeFile(absPath, config)

	if err != nil {
		return nil, err
	}

	return config, nil
}
