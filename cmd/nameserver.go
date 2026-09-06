package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/kevinharv/dns/handlers"
	"github.com/miekg/dns"
)

func main() {
	// Get configuration from CLI arguments
	port := flag.Int("port", 8053, "DNS server port")
	flag.Parse()

	// Load DNS records from database into memory. Pass
	// to handler function factory.
	loadedRecords := handlers.LoadRecords()
	dns.HandleFunc(".", handlers.HandleDNS(loadedRecords))

	// Start UDP listener
	go func() {
		srv := &dns.Server{Addr: fmt.Sprintf(":%d", *port), Net: "udp"}
		fmt.Printf("Starting UDP Server\n")
		if err := srv.ListenAndServe(); err != nil {
			log.Fatalf("Failed to start UDP listener: %s\n", err.Error())
		}
	}()

	// Start TCP listener
	go func() {
		srv := &dns.Server{Addr: fmt.Sprintf(":%d", *port), Net: "tcp"}
		fmt.Printf("Starting TCP Server\n")
		if err := srv.ListenAndServe(); err != nil {
			log.Fatalf("Failed to start TCP listener: %s\n", err.Error())
		}
	}()

	// Gracefully shutdown on interrupt
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	s := <-sig
	log.Fatalf("Signal (%v) received. Stopping server.", s)
}
