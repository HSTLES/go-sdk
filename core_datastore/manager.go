package core_datastore

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/hstles/go-sdk/core_config"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

// Manager ties together your CoreDB and a Turso API client.
type Manager struct {
	CoreDB *sql.DB
	turso  *TursoClient
}

// NewManager initializes CoreDB and the Turso platform client.
// cfg should come from core_config.LoadCoreConfig().
func NewManager(cfg *core_config.CoreConfig) (*Manager, error) {
	// --- open core DB ---
	dsn := fmt.Sprintf("%s?authToken=%s", cfg.TursoCoreDBURL, cfg.TursoAuthToken)
	coreDB, err := sql.Open("libsql", dsn)
	if err != nil {
		return nil, fmt.Errorf("open coredb: %w", err)
	}

	// enable foreign keys
	if _, err := coreDB.Exec("PRAGMA foreign_keys = ON;"); err != nil {
		coreDB.Close()
		return nil, fmt.Errorf("enable foreign keys on coredb: %w", err)
	}
	if err := coreDB.Ping(); err != nil {
		coreDB.Close()
		return nil, fmt.Errorf("ping coredb: %w", err)
	}

	// --- init Turso client ---
	// allow override via TURSO_API_URL, otherwise default to the official endpoint
	tursoAPIURL := core_config.RequireEnvDefault("TURSO_API_URL", "https://api.turso.sh/v1")
	turso := NewTursoClient(tursoAPIURL, cfg.TursoPlatformToken, &http.Client{})

	return &Manager{
		CoreDB: coreDB,
		turso:  turso,
	}, nil
}

// Close cleans up the CoreDB connection.
func (m *Manager) Close() error {
	return m.CoreDB.Close()
}
