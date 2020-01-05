package pg

import (
	"fmt"
	"net/url"
	"os/user"
	"strings"
)

// ConnectionSettings is a map containing PostgreSQL connection parameters.
type ConnectionSettings map[string]string

// ParseDatabaseURL parses given PostgreSQL connection string in URI format and returns connection params as a map.
// See the documentation: https://www.postgresql.org/docs/current/libpq-connect.html#LIBPQ-CONNSTRING
func ParseDatabaseURI(connString string) (ConnectionSettings, error) {
	url, err := url.Parse(connString)

	if err != nil {
		return nil, err
	}

	// Check that scheme is actually "postgresql"
	if url.Scheme != "postgres" && url.Scheme != "postgresql" {
		return nil, fmt.Errorf("invalid PostgreSQL database URI %s", url)
	}

	settings := defaultConnectionSettings()

	// Username and password
	if url.User != nil {
		settings["user"] = url.User.Username()

		if password, present := url.User.Password(); present {
			settings["password"] = password
		}
	}

	// Host
	if host := url.Hostname(); host != "" {
		settings["host"] = host
	}

	// Port
	if port := url.Port(); port != "" {
		settings["port"] = port
	}

	// Database name
	if database := strings.TrimLeft(url.Path, "/"); database != "" {
		settings["database"] = database
	} else {
		// Default is the same as database user
		settings["database"] = settings["user"]
	}

	for k, v := range url.Query() {
		settings[k] = v[0]
	}

	return settings, nil
}

// defaultConnectionSettings returns PostgreSQL implicit default connection settings.
func defaultConnectionSettings() ConnectionSettings {
	settings := ConnectionSettings{
		"host": "localhost",
		"port": "5432",
	}

	// Use local OS username as both default DB user and DB name
	if user, err := user.Current(); err == nil {
		settings["user"] = user.Username
		settings["database"] = user.Username
	}

	return settings
}
