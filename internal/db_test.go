package internal

import (
	"os"
	"testing"

	"github.com/miekg/dns"
)

func TestDatabaseRoundTrip(t *testing.T) {
	oldWD, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	tmpDir := t.TempDir()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := os.Chdir(oldWD); err != nil {
			t.Fatal(err)
		}
	}()

	db, _ := Open()
	if db == nil {
		t.Fatal("Open() returned nil database")
	}
	defer db.Close()

	var tableCount int
	if err := db.QueryRow("SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name='records'").Scan(&tableCount); err != nil {
		t.Fatalf("schema check failed: %v", err)
	}
	if tableCount != 1 {
		t.Fatalf("expected records table to exist, found %d", tableCount)
	}

	rr, err := dns.NewRR("example.com. 60 IN A 127.0.0.1")
	if err != nil {
		t.Fatalf("failed to create RR: %v", err)
	}
	if err := InsertRecords(db, []dns.RR{rr}); err != nil {
		t.Fatalf("InsertRecords() failed: %v", err)
	}

	var count int
	if err := db.QueryRow("SELECT COUNT(*) FROM records WHERE name = ? AND type = ?", "example.com.", "A").Scan(&count); err != nil {
		t.Fatalf("query failed: %v", err)
	}
	if count != 1 {
		t.Fatalf("expected 1 record after insert, found %d", count)
	}

	if err := DeleteRecords(db, []dns.RR{rr}); err != nil {
		t.Fatalf("DeleteRecords() failed: %v", err)
	}
	if err := db.QueryRow("SELECT COUNT(*) FROM records WHERE name = ? AND type = ?", "example.com.", "A").Scan(&count); err != nil {
		t.Fatalf("query after delete failed: %v", err)
	}
	if count != 0 {
		t.Fatalf("expected 0 records after delete, found %d", count)
	}
}
