package internal

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
	"github.com/miekg/dns"
)

func Open() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./dns.db")
	if err != nil {
		log.Fatal("failed to open database: ", err)
		return nil, err
	}

	if err := db.Ping(); err != nil {
		db.Close()
		log.Fatal("failed to ping database: ", err)
		return nil, &dns.Error{}
	}

	if err := WriteSchema(db); err != nil {
		db.Close()
		log.Fatal("failed to initialize schema: ", err)
		return nil, err
	}

	log.Printf("Connected to and pinged database")
	return db, nil
}

func WriteSchema(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS records (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		type TEXT NOT NULL,
		class TEXT NOT NULL,
		ttl INTEGER NOT NULL,
		data TEXT NOT NULL,
		UNIQUE(name, type, class, ttl, data)
	);`

	_, err := db.Exec(query)
	return err
}

func InsertRecords(db *sql.DB, records []dns.RR) error {
	if db == nil {
		return fmt.Errorf("database is nil")
	}

	for _, rr := range records {
		hdr := rr.Header()
		data := rr.String()
		if data == "" {
			data = fmt.Sprintf("%s %d IN %s %s", hdr.Name, hdr.Ttl, dns.TypeToString[hdr.Rrtype], rr.String())
		}

		_, err := db.Exec(
			`INSERT OR IGNORE INTO records (name, type, class, ttl, data) VALUES (?, ?, ?, ?, ?)`,
			hdr.Name,
			dns.TypeToString[hdr.Rrtype],
			dns.ClassToString[hdr.Class],
			hdr.Ttl,
			data,
		)
		if err != nil {
			return fmt.Errorf("insert record %q: %w", hdr.Name, err)
		}
	}
	return nil
}

func GetRecords(db *sql.DB) ([]dns.RR, error) {
	if db == nil {
		return nil, fmt.Errorf("database is nil")
	}

	records := []dns.RR{}

	res, err := db.Query(`SELECT name, type, class, ttl, data FROM records`)
	if err != nil {
		return nil, fmt.Errorf("error querying DB for records: %s", err.Error())
	}
	defer res.Close()

	for res.Next() {
		var name, rtype, class string
		var ttl int
		var data string

		if err := res.Scan(&name, &rtype, &class, &ttl, &data); err != nil {
			return nil, fmt.Errorf("failed to scan record row: %w", err)
		}

		rec, err := dns.NewRR(data)
		if err != nil {
			return nil, fmt.Errorf("failed to parse record %q from DB: %w", data, err)
		}

		records = append(records, rec)
	}

	if err := res.Err(); err != nil {
		return nil, fmt.Errorf("error encountered in iteration: %s", err.Error())
	}

	return records, nil
}

func DeleteRecords(db *sql.DB, records []dns.RR) error {
	if db == nil {
		return fmt.Errorf("database is nil")
	}

	for _, rr := range records {
		hdr := rr.Header()
		_, err := db.Exec(
			`DELETE FROM records WHERE name = ? AND type = ? AND class = ? AND ttl = ?`,
			hdr.Name,
			dns.TypeToString[hdr.Rrtype],
			dns.ClassToString[hdr.Class],
			hdr.Ttl,
		)
		if err != nil {
			return fmt.Errorf("delete record %q: %w", hdr.Name, err)
		}
	}
	return nil
}
