package cli

import (
	"os"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadServerConfig(t *testing.T) {
	wd, err := os.Getwd()

	if err != nil {
		panic(err)
	}

	configFileFixtureRelativePath := filepath.Join("testdata", "example.toml")
	configFileFixtureAbsolutePath := filepath.Join(wd, configFileFixtureRelativePath)

	t.Run("parses toml and returns config", func(t *testing.T) {
		config, err := readServerConfig(configFileFixtureAbsolutePath)

		assert.NoError(t, err)
		assert.Equal(t, config.Listen, "0.0.0.0:4242")
		assert.Equal(t, config.DatabaseURL, "postgres://localhost/hired_dev")
	})

	t.Run("it resolves config path relative to cwd", func(t *testing.T) {
		_, err := readServerConfig(configFileFixtureRelativePath)
		assert.NoError(t, err)
	})

	t.Run("returns error if file doesn't exist", func(t *testing.T) {
		_, err := readServerConfig("nonexistent file")

		assert.Error(t, err)
		assert.Regexp(t, regexp.QuoteMeta("no such file or directory"), err.Error())
	})
}
