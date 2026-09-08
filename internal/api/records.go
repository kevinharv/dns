package api

import (
	"fmt"

	"github.com/kevinharv/dns/internal/config"
	"github.com/kevinharv/dns/internal/database"
	"github.com/kevinharv/dns/internal/store"
	"github.com/miekg/dns"
)

var seedRecords = []dns.RR{
	makeRR("www.kevharv.com.", "1.2.3.4", "A", 60),
	makeRR("lab.kevharv.com.", "9.9.9.9", "A", 60),
	makeRR("jellyfin.kevharv.com.", "lab.kevharv.com.", "CNAME", 60),
}

func LoadRecords(configuration *config.DNSServerConfiguration) (*store.Store, error) {
	db, err := database.Open(configuration)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	if err := database.InsertRecords(db, seedRecords); err != nil {
		return nil, err
	}
	stored, err := database.GetRecords(db)
	if err != nil {
		return nil, err
	}
	return store.New(stored), nil
}

func makeRR(name string, value string, recordType string, ttl int) dns.RR {
	record, _ := dns.NewRR(fmt.Sprintf("%s %d IN %s %s", name, ttl, recordType, value))
	return record
}
