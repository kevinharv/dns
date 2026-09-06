package handlers

import (
	"log"
	"math/rand/v2"
	"time"

	"github.com/miekg/dns"
)

type DNSServer struct {
	client   *dns.Client
	upstream []string
}

func Setup() *DNSServer {
	server := DNSServer{}
	server.client = &dns.Client{
		Net:     "udp",
		Timeout: 5 * time.Second,
	}
	server.upstream = []string{
		"1.1.1.1:53",
		"8.8.8.8:53",
	}

	return &server
}

// HandleDNS takes in a map of records and returns a handling
// function for answering DNS questions.
func (s *DNSServer) HandleDNS(records map[string]dns.RR) func(w dns.ResponseWriter, r *dns.Msg) {
	return func(w dns.ResponseWriter, r *dns.Msg) {
		m := new(dns.Msg)
		m.SetReply(r)
		m.Authoritative = true

		for _, q := range r.Question {
			rr, ok := records[q.Name]
			if ok {
				log.Printf("Question: %s | Answer: %s", q.Name, rr.String())
				m.Answer = append(m.Answer, rr)
			}
		}

		// If answered, return early
		if len(m.Answer) > 0 {
			w.WriteMsg(m)
			return
		}

		// If not answered, round-robin forward to upstream
		// TODO - check authority, return NXDOMAIN if authoritative AND not found
		if len(m.Answer) == 0 {
			resp, _, err := s.client.Exchange(r, s.upstream[rand.IntN(len(s.upstream))])
			if err != nil {
				log.Printf("ERROR: Failed to exchange with upstream: %s", err.Error())
				m.SetRcode(r, dns.RcodeServerFailure)
				w.WriteMsg(m)
				return
			}

			w.WriteMsg(resp)
		}
	}
}
