package postgres

import (
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/kevinharv/dns/internal/config"
	_ "github.com/mattn/go-sqlite3"
)

const defaultDBPath = "./dns.db"

// Open initializes the management-plane database. The current implementation
// uses SQLite; the package boundary allows the backend to change independently.
func Open(configurations ...*config.DNSServerConfiguration) (*sql.DB, error) {
	dbPath := defaultDBPath
	if len(configurations) > 0 && configurations[0] != nil && configurations[0].SQLiteDBPath != "" {
		dbPath = configurations[0].SQLiteDBPath
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		slog.Error("failed to open database", "error", err.Error())
		return nil, err
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}
	if err := WriteSchema(db); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to initialize schema: %w", err)
	}

	slog.Info("Connected to and pinged database")
	return db, nil
}

func WriteSchema(db *sql.DB) error {
	if db == nil {
		return fmt.Errorf("database is nil")
	}

	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS zones (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL UNIQUE
		);

		CREATE TABLE IF NOT EXISTS records (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			type TEXT NOT NULL,
			class TEXT NOT NULL,
			ttl INTEGER NOT NULL,
			data TEXT NOT NULL,
			UNIQUE(name, type, class, ttl, data)
		);`)
	return err
}
