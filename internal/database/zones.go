package database

import (
	"database/sql"
	"fmt"
)

func InsertZones(db *sql.DB, zones []string) error {
	if db == nil {
		return fmt.Errorf("database is nil")
	}
	for _, zone := range zones {
		if zone == "" {
			return fmt.Errorf("zone name is empty")
		}
		if _, err := db.Exec(`INSERT OR IGNORE INTO zones (name) VALUES (?)`, zone); err != nil {
			return fmt.Errorf("insert zone %q: %w", zone, err)
		}
	}
	return nil
}

func GetZones(db *sql.DB) ([]string, error) {
	if db == nil {
		return nil, fmt.Errorf("database is nil")
	}
	rows, err := db.Query(`SELECT name FROM zones ORDER BY name`)
	if err != nil {
		return nil, fmt.Errorf("error querying DB for zones: %w", err)
	}
	defer rows.Close()

	zones := []string{}
	for rows.Next() {
		var zone string
		if err := rows.Scan(&zone); err != nil {
			return nil, fmt.Errorf("failed to scan zone row: %w", err)
		}
		zones = append(zones, zone)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error encountered in zone iteration: %w", err)
	}
	return zones, nil
}

func DeleteZones(db *sql.DB, zones []string) error {
	if db == nil {
		return fmt.Errorf("database is nil")
	}
	for _, zone := range zones {
		if _, err := db.Exec(`DELETE FROM zones WHERE name = ?`, zone); err != nil {
			return fmt.Errorf("delete zone %q: %w", zone, err)
		}
	}
	return nil
}
