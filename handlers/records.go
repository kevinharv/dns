package handlers

import (
	"fmt"
	"github.com/miekg/dns"
)

// TODO - replace with DB read
var fakeDB = map[string]dns.RR{
	"www.kevharv.com.":      makeRR("www.kevharv.com.", "1.2.3.4", "A", 60),
	"lab.kevharv.com.":      makeRR("lab.kevharv.com.", "9.9.9.9", "A", 60),
	"jellyfin.kevharv.com.": makeRR("jellyfin.kevharv.com.", "lab.kevharv.com.", "CNAME", 60),
}

// makeRR creates a DNS resource record from typed inputs. It ignores errors, returning
// only the RR itself. Simple wrapper around dns.NewRR().
func makeRR(name string, value string, rtype string, ttl int) dns.RR {
	r, _ := dns.NewRR(fmt.Sprintf("%s %d IN %s %s", name, ttl, rtype, value))
	return r
}

// LoadRecords retrieves DNS records from the database and constructs RR objects
// from them. It returns a map of names to records.
func LoadRecords() map[string]dns.RR {
	records := map[string]dns.RR{}
	
	// TODO - open DB connection

	// PLACEHOLDER - batch records in from DB
	for name, rec := range fakeDB {
		records[name] = rec
		fmt.Printf("Loaded Record: %s\n", name)
	}
	
	return records
}
