package pgmeta

import (
	"fmt"
	"os"
	"os/user"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var pgUser, pgHost, pgPort, pgDatabase, dsn string

func init() {
	// Postgresql username
	pgUser = os.Getenv("PGUSER")

	if pgUser == "" {
		osUser, err := user.Current()

		if err != nil {
			panic(err)
		}

		pgUser = osUser.Username
	}

	// Postgresql hostname
	pgHost = os.Getenv("PGHOST")

	if pgHost == "" {
		pgHost = "localhost"
	}

	// Postgresql port
	pgPort = os.Getenv("PGPORT")

	if pgPort == "" {
		pgPort = "5432"
	}

	// Name of the database we created for testing purposes (see Makefile).
	pgDatabase = os.Getenv("PGDATABASE")

	if pgDatabase == "" {
		pgDatabase = "gevolut_test"
	}

	// Format DSN for the inspector
	dsn = fmt.Sprintf("postgresql://%s@%s:%s/%s?sslmode=disable", pgUser, pgHost, pgPort, pgDatabase)
}

func TestInspectorDatabaseName(t *testing.T) {
	inspector, err := Inspect(dsn)
	require.NoError(t, err)

	defer inspector.Close()

	name, err := inspector.DatabaseName()
	require.NoError(t, err)

	assert.Equal(t, pgDatabase, name)
}

func TestInspectorOIDTableMapping(t *testing.T) {
	inspector, err := Inspect(dsn)
	require.NoError(t, err)

	defer inspector.Close()

	mapping, err := inspector.OIDTableMapping()
	require.NoError(t, err)

	//
	// XXX: there is no way to check the OID values as they are changed each time database is re-created.
	//

	tables := make([]Table, 0, len(mapping))

	for oid := range mapping {
		tables = append(tables, mapping[oid])
	}

	assert.ElementsMatch(t, tables, []Table{
		{"public", "companies"},
		{"public", "users"},
	})
}

func TestInspectorClose(t *testing.T) {
	inspector, err := Inspect(dsn)
	require.NoError(t, err)

	err = inspector.Close()
	require.NoError(t, err)

	_, err = inspector.DatabaseName()
	assert.Error(t, err)
}
