package server

// Config contains configuration parameters for the server package.
// cli package use this to unmarshall the gevulot.toml.
type Config struct {
	// Local IP address and port on which Gevolut will listen for client connections.
	Listen string

	// Database connection string for the proxied PostgreSQL server.
	DatabaseURL string
}
