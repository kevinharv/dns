package api

import (
	"fmt"

	"github.com/kevinharv/dns/internal/config"
	"github.com/kevinharv/dns/internal/postgres"
	"github.com/kevinharv/dns/internal/store"
	"github.com/miekg/dns"
)

var seedRecords = []dns.RR{
	makeRR("www.kevharv.com.", "1.2.3.4", "A", 60),
	makeRR("lab.kevharv.com.", "9.9.9.9", "A", 60),
	makeRR("jellyfin.kevharv.com.", "lab.kevharv.com.", "CNAME", 60),
}

func LoadRecords(configuration *config.DNSServerConfiguration) (*store.Store, error) {
	db, err := postgres.Open(configuration)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	if err := postgres.InsertRecords(db, seedRecords); err != nil {
		return nil, err
	}
	stored, err := postgres.GetRecords(db)
	if err != nil {
		return nil, err
	}
	return store.New(stored), nil
}

func makeRR(name, value, recordType string, ttl int) dns.RR {
	record, _ := dns.NewRR(fmt.Sprintf("%s %d IN %s %s", name, ttl, recordType, value))
	return record
}
