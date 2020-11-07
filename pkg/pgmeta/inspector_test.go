package pgmeta

import (
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// nolint:gochecknoglobals
var databaseURL = os.Getenv("DATABASE_URL")

// nolint:gochecknoinits
func init() {
	if databaseURL == "" {
		panic("no DATABASE_URL specified")
	}
}

func TestInspectorDatabaseName(t *testing.T) {
	// Extract db name from the DSN
	url, err := url.Parse(databaseURL)
	require.NoError(t, err)

	databaseName := strings.TrimLeft(url.Path, "/")
	require.NotEmpty(t, databaseName)

	inspector, err := Inspect(databaseURL)
	require.NoError(t, err)

	defer inspector.Close()

	name, err := inspector.DatabaseName()
	require.NoError(t, err)

	assert.Equal(t, name, databaseName)
}

func TestInspectorOIDTableMapping(t *testing.T) {
	inspector, err := Inspect(databaseURL)
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
	inspector, err := Inspect(databaseURL)
	require.NoError(t, err)

	err = inspector.Close()
	require.NoError(t, err)

	_, err = inspector.DatabaseName()
	assert.Error(t, err)
}
