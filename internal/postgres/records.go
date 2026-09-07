package postgres

import (
	"database/sql"
	"fmt"

	"github.com/miekg/dns"
)

func InsertRecords(db *sql.DB, records []dns.RR) error {
	if db == nil {
		return fmt.Errorf("database is nil")
	}
	for _, record := range records {
		if record == nil {
			return fmt.Errorf("record name and data are required")
		}
		header := record.Header()
		if _, err := db.Exec(`INSERT OR IGNORE INTO records (name, type, class, ttl, data) VALUES (?, ?, ?, ?, ?)`, header.Name, dns.TypeToString[header.Rrtype], dns.ClassToString[header.Class], header.Ttl, record.String()); err != nil {
			return fmt.Errorf("insert record %q: %w", header.Name, err)
		}
	}
	return nil
}

func GetRecords(db *sql.DB) ([]dns.RR, error) {
	if db == nil {
		return nil, fmt.Errorf("database is nil")
	}
	rows, err := db.Query(`SELECT name, type, class, ttl, data FROM records ORDER BY name, type, data`)
	if err != nil {
		return nil, fmt.Errorf("error querying DB for records: %w", err)
	}
	defer rows.Close()

	records := []dns.RR{}
	for rows.Next() {
		var name, recordType, class, data string
		var ttl uint32
		if err := rows.Scan(&name, &recordType, &class, &ttl, &data); err != nil {
			return nil, fmt.Errorf("failed to scan record row: %w", err)
		}
		record, err := dns.NewRR(data)
		if err != nil {
			return nil, fmt.Errorf("failed to parse record %q: %w", data, err)
		}
		records = append(records, record)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error encountered in record iteration: %w", err)
	}
	return records, nil
}

func DeleteRecords(db *sql.DB, records []dns.RR) error {
	if db == nil {
		return fmt.Errorf("database is nil")
	}
	for _, record := range records {
		if record == nil {
			return fmt.Errorf("record is nil")
		}
		header := record.Header()
		if _, err := db.Exec(`DELETE FROM records WHERE name = ? AND type = ? AND class = ? AND ttl = ? AND data = ?`, header.Name, dns.TypeToString[header.Rrtype], dns.ClassToString[header.Class], header.Ttl, record.String()); err != nil {
			return fmt.Errorf("delete record %q: %w", header.Name, err)
		}
	}
	return nil
}
