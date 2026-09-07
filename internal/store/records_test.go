package store

import (
	"testing"

	"github.com/miekg/dns"
)

func TestReplaceAndSnapshot(t *testing.T) {
	initial := []dns.RR{mustRR("example.com. 60 IN A 127.0.0.1")}
	store := New(initial)
	initial[0] = mustRR("example.com. 60 IN A 192.0.2.1")

	snapshot := store.Snapshot()
	if len(snapshot) != 1 || snapshot[0].String() != "example.com.\t60\tIN\tA\t127.0.0.1" {
		t.Fatalf("Snapshot() = %v, want an independent copy", snapshot)
	}

	replacement := []dns.RR{mustRR("www.example.com. 60 IN A 127.0.0.1")}
	store.Replace(replacement)
	if _, ok := store.Get("example.com."); ok {
		t.Fatal("Get() found a record from the previous snapshot")
	}
	if record, ok := store.Get("www.example.com."); !ok || record.String() != replacement[0].String() {
		t.Fatalf("Get() = %v, %v", record, ok)
	}
}

func mustRR(text string) dns.RR {
	record, err := dns.NewRR(text)
	if err != nil {
		panic(err)
	}
	return record
}
