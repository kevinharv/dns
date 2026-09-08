package database

import (
	"testing"

	"github.com/kevinharv/dns/internal/config"
	"github.com/miekg/dns"
)

func TestRecordRoundTrip(t *testing.T) {
	db, err := Open(&config.DNSServerConfiguration{SQLiteDBPath: ":memory:"})
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	record, err := dns.NewRR("example.com. 60 IN A 127.0.0.1")
	if err != nil {
		t.Fatal(err)
	}
	if err := InsertRecords(db, []dns.RR{record}); err != nil {
		t.Fatalf("InsertRecords() failed: %v", err)
	}
	records, err := GetRecords(db)
	if err != nil || len(records) != 1 || records[0].String() != record.String() {
		t.Fatalf("GetRecords() = %v, %v", records, err)
	}
	if err := DeleteRecords(db, []dns.RR{record}); err != nil {
		t.Fatalf("DeleteRecords() failed: %v", err)
	}
	records, err = GetRecords(db)
	if err != nil || len(records) != 0 {
		t.Fatalf("GetRecords() after delete = %v, %v", records, err)
	}
}
