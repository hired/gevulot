package pg

import (
	"os/user"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseDatabaseURI(t *testing.T) {
	user, err := user.Current()
	assert.NoError(t, err)

	defaultUser := user.Username

	//nolint:lll
	validTestCases := map[string]ConnectionSettings{
		"postgresql://":                       {"host": "localhost", "port": "5432", "user": defaultUser, "database": defaultUser},
		"postgres://":                         {"host": "localhost", "port": "5432", "user": defaultUser, "database": defaultUser},
		"postgresql://hired.com":              {"host": "hired.com", "port": "5432", "user": defaultUser, "database": defaultUser},
		"postgresql://hired.com:5433":         {"host": "hired.com", "port": "5433", "user": defaultUser, "database": defaultUser},
		"postgresql://localhost/mydb":         {"host": "localhost", "port": "5432", "user": defaultUser, "database": "mydb"},
		"postgresql://alice@example.com":      {"host": "example.com", "port": "5432", "user": "alice", "database": "alice"},
		"postgresql://bob:secret@example.com": {"host": "example.com", "port": "5432", "user": "bob", "password": "secret", "database": "bob"},
		"postgresql://charlie@hired.com/otherdb?connect_timeout=10&application_name=myapp": {"host": "hired.com", "port": "5432", "user": "charlie", "database": "otherdb", "connect_timeout": "10", "application_name": "myapp"},
	}

	for uri, expected := range validTestCases {
		settings, err := ParseDatabaseURI(uri)

		assert.NoError(t, err)
		assert.Equalf(t, expected, settings, "uri: %s", uri)
	}

	invalidTestCases := []string{
		"mysql://localhost",
		"postgresql://%$&^",
	}

	for _, uri := range invalidTestCases {
		_, err := ParseDatabaseURI(uri)
		assert.Errorf(t, err, "uri: %s", uri)
	}
}
