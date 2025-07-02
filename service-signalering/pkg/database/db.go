package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func Init() {
	var err error
	connStr := "host=localhost port=5432 user=signalering_user password=signalering_pass dbname=signalering sslmode=disable"

	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Gefaald:", err)
	}

	// Test connectie
	if err = DB.Ping(); err != nil {
		log.Fatal("Database is kapot want:", err)
	}

	fmt.Println("Success!")
}
