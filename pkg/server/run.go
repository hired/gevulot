package server

import (
	log "github.com/sirupsen/logrus"
)

// Run starts the PG proxy server.
func Run(configChan <-chan *Config) error {
	for {
		select {
		case config := <-configChan:
			log.Infof("server config: %#v", config)
		}
	}

	return nil
}
