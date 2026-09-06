# DNS

*DNS nameserver and resolver for the modern enterprise.*

## Run

This repository includes a minimal DNS server scaffold using `github.com/miekg/dns`.

- Build and run:

```bash
go run ./cmd/nameserver.go
```

- The server listens on port `8053` (UDP and TCP). Query it with `dig`:

```bash
dig @127.0.0.1 -p 8053 example.local A
```

You should see an A record pointing to `127.0.0.1` for `example.local.`

To run on port 53, start the binary as root or with capabilities to bind privileged ports.
