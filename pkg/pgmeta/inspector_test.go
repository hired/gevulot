package pgmeta

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Name of the database we created for testing purposes (see Makefile).
const TestDBName = "gevulot_test"

func TestInspectorDatabaseName(t *testing.T) {
	inspector, err := Inspect("postgresql:///" + TestDBName + "?sslmode=disable")
	require.NoError(t, err)

	defer inspector.Close()

	name, err := inspector.DatabaseName()
	require.NoError(t, err)

	assert.Equal(t, TestDBName, name)
}

func TestInspectorOIDTableMapping(t *testing.T) {
	inspector, err := Inspect("postgresql:///" + TestDBName + "?sslmode=disable")
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
	inspector, err := Inspect("postgresql:///" + TestDBName + "?sslmode=disable")
	require.NoError(t, err)

	err = inspector.Close()
	require.NoError(t, err)

	_, err = inspector.DatabaseName()
	assert.Error(t, err)
}
