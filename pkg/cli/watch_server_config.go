package cli

import (
	log "github.com/sirupsen/logrus"

	"github.com/hired/gevulot/pkg/server"
)

// watchServerConfig watches for the config file and sends updated config to the given channel.
func watchServerConfig(configPath string, configChan chan *server.Config) error {
	watcher := newFileWatcher(configPath)
	watcher.OnWrite = func() {
		updatedConfig, err := readServerConfig(configPath)

		if err != nil {
			log.Errorf("error loading config file %s: %v", configPath, err)
			return
		}

		configChan <- updatedConfig
	}

	return watcher.Watch()
}
