package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	_ "github.com/microsoft/go-mssqldb"
)

type Config struct {
	Connection_String string `json:"connection_string"`
}

func readConfig(path string) Config {
	jsonFile, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()

	byteValue, _ := io.ReadAll(jsonFile)
	var config Config
	json.Unmarshal(byteValue, &config)
	return config
}

var database *sql.DB

func Database_Get() *sql.DB {
	if database != nil {
		return database
	}

	config := readConfig("config.json")

	db, err := sql.Open("mssql", config.Connection_String)
	if err != nil {
		log.Fatal(err)
	}
	database = db
	return db
}
