package config

import (
	"encoding/json"
	"log/slog"
	"os"
)

type DNSServerConfiguration struct {
	DisplayName        string   `json:"displayName"`
	FQDN               string   `json:"fqdn"`
	Port               int      `json:"port"`
	TLSEnabled         bool     `json:"tlsEnabled"`
	TLSCertificatePath string   `json:"tlsCertificatePath"`
	TLSPrivateKeyPath  string   `json:"tlsPrivateKeyPath"`
	DNSUpstreamServers []string `json:"dnsUpstreamServers"`
	SQLiteDBPath       string   `json:"sqliteDBPath"`
}

func (c *DNSServerConfiguration) LoadFromFile(path string) error {
	content, err := os.ReadFile(path)
	if err != nil {
		slog.Error("Failed to open DNS server configuration", "path", path)
		return err
	}

	if err := json.Unmarshal(content, c); err != nil {
		slog.Error("Failed to read DNS server configuration", "path", path)
		return err
	}

	slog.Info("Loaded DNS server configuration", "path", path)
	return nil
}
