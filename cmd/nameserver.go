package main

import (
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/kevinharv/dns/internal/api"
	"github.com/kevinharv/dns/internal/config"
	"github.com/kevinharv/dns/internal/server"
	"github.com/miekg/dns"
)

func main() {	
	// Get configuration from file
	dnsConfig := config.DNSServerConfiguration{}
	err := dnsConfig.LoadFromFile("scripts/dns_config.json")
	if err != nil {
		slog.Error("Configuration load failed", "error", err.Error())
	}

	// Setup logging
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})
	logger := slog.New(handler)
	slog.SetDefault(logger)
	slog.Info("Starting DNS Nameserver")

	// Load DNS records from database into memory. Pass
	// to handler function factory.
	slog.Debug("Loading records...")
	loadedRecords, err := api.LoadRecords(&dnsConfig)
	if err != nil {
		slog.Error("Record load failed", "error", err)
		return
	}
	slog.Debug("Loaded records", "count", len(loadedRecords.Snapshot()))

	// Create server and setup handling
	srv := server.New(dnsConfig)
	dns.HandleFunc(".", srv.HandleDNS(loadedRecords))

	// Start UDP listener
	go func() {
		srv := &dns.Server{Addr: fmt.Sprintf(":%d", dnsConfig.Port), Net: "udp"}
		slog.Info("Starting UDP Server")
		if err := srv.ListenAndServe(); err != nil {
			slog.Error("Failed to start UDP listener", "error", err)
		}
	}()

	// Start TCP listener
	go func() {
		srv := &dns.Server{Addr: fmt.Sprintf(":%d", dnsConfig.Port), Net: "tcp"}
		slog.Info("Starting TCP Server")
		if err := srv.ListenAndServe(); err != nil {
			slog.Error("Failed to start TCP listener", "error", err)
		}
	}()

	// Gracefully shutdown on interrupt
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	s := <-sig
	slog.Info("Received stop signal. Stopping server...", "signal", s.String())
}
