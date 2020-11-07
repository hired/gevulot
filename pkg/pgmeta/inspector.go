package pgmeta

import (
	"database/sql"
	"fmt"

	// This package is the only place where we use PG driver
	_ "github.com/lib/pq"
	"github.com/lib/pq/oid"
)

// Table represents fully qualified table name.
type Table struct {
	Schema string
	Name   string
}

// Inspector allows getting meta-information about PostgreSQL database.
type Inspector struct {
	db *sql.DB
}

// Inspect connects to the database with the given DSN (connection string) and
// initializes a new Inspector instance for it.
func Inspect(dsn string) (*Inspector, error) {
	db, err := sql.Open("postgres", dsn)

	if err != nil {
		return nil, err
	}

	// sql.Open doesn't actually connect to a DB. We need to ping it to check that the connection
	// string is valid.
	err = db.Ping()

	if err != nil {
		db.Close()
		return nil, fmt.Errorf("inspector: error connecting to the database: %w", err)
	}

	return &Inspector{db: db}, nil
}

// DatabaseName returns database name (ie. gevulot_test).
func (i *Inspector) DatabaseName() (string, error) {
	var name string

	row := i.db.QueryRow(`SELECT current_database();`)
	err := row.Scan(&name)

	return name, err
}

// OIDTableMapping returns a list of all database tables and its respective OIDs.
func (i *Inspector) OIDTableMapping() (map[oid.Oid]Table, error) {
	rows, err := i.db.Query(`
      SELECT c.oid AS oid
           , c.relname AS table
           , n.nspname AS schema
        FROM pg_class c
        JOIN pg_catalog.pg_namespace n ON n.oid = c.relnamespace
       WHERE n.nspname NOT IN ('information_schema')
         AND n.nspname NOT LIKE 'pg_%'
		 AND c.relkind IN ('r', 'm', 'v');`,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	mapping := make(map[oid.Oid]Table)

	for rows.Next() {
		var tableOid oid.Oid

		var table Table

		err := rows.Scan(&tableOid, &table.Name, &table.Schema)

		if err != nil {
			return nil, err
		}

		mapping[tableOid] = table
	}

	return mapping, nil
}

// Close closes database connection.
func (i *Inspector) Close() error {
	return i.db.Close()
}
