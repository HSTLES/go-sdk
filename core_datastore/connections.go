package core_datastore

import (
	"database/sql"
	"fmt"

	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

// ConnectToInstance returns an *sql.DB for any Turso instance URL + token.
func ConnectToInstance(instanceURL, authToken string) (*sql.DB, error) {
	full := instanceURL + "?authToken=" + authToken
	db, err := sql.Open("libsql", full)
	if err != nil {
		return nil, fmt.Errorf("open instance: %w", err)
	}
	if _, err := db.Exec("PRAGMA foreign_keys = ON;"); err != nil {
		return nil, fmt.Errorf("enable fkeys: %w", err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("ping instance: %w", err)
	}
	return db, nil
}
