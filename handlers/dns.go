package handlers

import (
	"fmt"

	"github.com/miekg/dns"
)

// HandleDNS takes in a map of records and returns a handling
// function for answering DNS questions.
func HandleDNS(records map[string]dns.RR) func(w dns.ResponseWriter, r *dns.Msg) {
	fmt.Printf("Quantity of Records: %d\n", len(records))

	return func(w dns.ResponseWriter, r *dns.Msg) {
		m := new(dns.Msg)
		m.SetReply(r)
		m.Authoritative = true

		for _, q := range r.Question {
			fmt.Printf("Question: %s\n", q.Name)
			rr, ok := records[q.Name]
			if ok {
				fmt.Printf("Found Answer for: %s\n", q.Name)
				m.Answer = append(m.Answer, rr)
			}
		}

		// Assume authoritative for the domain
		// Return NXDOMAIN if not found
		if len(m.Answer) == 0 {
			m.SetRcode(r, dns.RcodeNameError)
		}

		w.WriteMsg(m)
	}
}
