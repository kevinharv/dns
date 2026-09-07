package postgres

import (
	"testing"

	"github.com/kevinharv/dns/internal/config"
)

func TestZoneRoundTrip(t *testing.T) {
	db, err := Open(&config.DNSServerConfiguration{SQLiteDBPath: ":memory:"})
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	zones := []string{"example.com.", "internal.example.com."}
	if err := InsertZones(db, zones); err != nil {
		t.Fatalf("InsertZones() failed: %v", err)
	}
	loaded, err := GetZones(db)
	if err != nil || len(loaded) != len(zones) {
		t.Fatalf("GetZones() = %v, %v", loaded, err)
	}
	if err := DeleteZones(db, []string{zones[0]}); err != nil {
		t.Fatalf("DeleteZones() failed: %v", err)
	}
	loaded, err = GetZones(db)
	if err != nil || len(loaded) != 1 || loaded[0] != zones[1] {
		t.Fatalf("GetZones() after delete = %v, %v", loaded, err)
	}
}
