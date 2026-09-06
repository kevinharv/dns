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
	// Setup logging
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Printf("Starting DNS Nameserver")

	// Get configuration from CLI arguments
	port := flag.Int("port", 8053, "DNS server port")
	flag.Parse()
	log.Printf("DNS Port: %d", *port)


	// Load DNS records from database into memory. Pass
	// to handler function factory.
	log.Printf("Loading records...")
	loadedRecords := handlers.LoadRecords()
	log.Printf("Loaded %d records", len(loadedRecords))

	// Create server and setup handling
	srv := handlers.Setup()
	dns.HandleFunc(".", srv.HandleDNS(loadedRecords))

	// Start UDP listener
	go func() {
		srv := &dns.Server{Addr: fmt.Sprintf(":%d", *port), Net: "udp"}
		log.Printf("Starting UDP Server")
		if err := srv.ListenAndServe(); err != nil {
			log.Fatalf("Failed to start UDP listener: %s\n", err.Error())
		}
	}()

	// Start TCP listener
	go func() {
		srv := &dns.Server{Addr: fmt.Sprintf(":%d", *port), Net: "tcp"}
		log.Printf("Starting TCP Server")
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
