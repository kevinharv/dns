package server

import (
	"log/slog"
	"math/rand/v2"
	"time"

	"github.com/kevinharv/dns/internal/config"
	"github.com/kevinharv/dns/internal/store"
	"github.com/miekg/dns"
)

type DNSServer struct {
	client   *dns.Client
	upstream []string
}

func New(c config.DNSServerConfiguration) *DNSServer {
	server := DNSServer{
		client: &dns.Client{
			Net:     "udp",
			Timeout: 5 * time.Second,
		},
		upstream: c.DNSUpstreamServers}
	return &server
}

func (s *DNSServer) HandleDNS(records *store.Store) func(w dns.ResponseWriter, r *dns.Msg) {
	return func(w dns.ResponseWriter, r *dns.Msg) {
		message := new(dns.Msg)
		message.SetReply(r)
		message.Authoritative = true

		for _, question := range r.Question {
			record, ok := records.Get(question.Name)
			if !ok {
				continue
			}
			slog.Debug("DNS answer", "question", question.Name, "answer", record.String())
			message.Answer = append(message.Answer, record)
		}

		if len(message.Answer) > 0 {
			_ = w.WriteMsg(message)
			return
		}

		response, _, err := s.client.Exchange(r, s.upstream[rand.IntN(len(s.upstream))])
		if err != nil {
			slog.Error("failed to exchange with upstream", "error", err)
			message.SetRcode(r, dns.RcodeServerFailure)
			_ = w.WriteMsg(message)
			return
		}

		w.WriteMsg(response)
	}
}
