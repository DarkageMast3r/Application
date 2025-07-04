package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

var DB *sql.DB

type config struct {
	ConnectionString string `json:"connection_string"`
}

func Init() {
	var config config
	jsonFile, err := os.Open("config.json")
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()
	configBody, err := io.ReadAll(jsonFile)
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(configBody, &config)
	if err != nil {
		log.Fatal(err)
	}
	DB, err = sql.Open("mysql", config.ConnectionString)
	if err != nil {
		log.Fatal("Gefaald:", err)
	}

	// Test connectie
	if err = DB.Ping(); err != nil {
		log.Fatal("Database is kapot want:", err)
	}

	fmt.Println("Success!")
}
