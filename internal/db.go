package internal

import (
	"log"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

func Open() {
	db, err := sql.Open("sqlite", "./dns.db")
	if err != nil {
		log.Fatal("failed to open database: ", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal("failed to ping database: ", err)
	}

	log.Printf("Connected to and pinged database")
}
