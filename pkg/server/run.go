package server

// Run starts the PG proxy server.
func Run(configChan <-chan *Config) error {
	cfg := NewConfigDistributor(configChan)
	defer cfg.Close()

	srv := NewServer(cfg)
	defer srv.Close()

	return srv.Start()
}
