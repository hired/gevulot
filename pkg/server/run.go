package server

// Run starts the PG proxy server.
func Run(log Logger, configChan <-chan *Config) error {
	for {
		select {
		case config := <-configChan:
			log.Infof("server config: %#v", config)
		}
	}

	return nil
}
