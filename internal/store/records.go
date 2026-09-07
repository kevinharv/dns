package store

import (
	"sync/atomic"

	"github.com/miekg/dns"
)

// Store provides lock-free snapshots and atomic replacements of DNS records.
type Store struct {
	records atomic.Value
}

func New(records []dns.RR) *Store {
	store := &Store{}
	store.records.Store(copyRecords(records))
	return store
}

func (s *Store) Replace(records []dns.RR) {
	s.records.Store(copyRecords(records))
}

func (s *Store) Snapshot() []dns.RR {
	return copyRecords(s.records.Load().([]dns.RR))
}

func (s *Store) Get(name string) (dns.RR, bool) {
	for _, record := range s.records.Load().([]dns.RR) {
		if record.Header().Name == name {
			return record, true
		}
	}
	return nil, false
}

func copyRecords(records []dns.RR) []dns.RR {
	return append([]dns.RR(nil), records...)
}
