package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"main/handlers"
	"main/models"
	"net/http"
	"os"
	"strconv"
)

type Config struct {
	Service_discovery_root string `json:"service_discovery_root"`
	Service_discovery_port int    `json:"service_discovery_port"`
	Allow_insecure         bool   `json:"allow_insecure"`
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

func main() {
	config := readConfig("../config.json")
	var port_str = strconv.Itoa(config.Service_discovery_port)
	models.Service_Init()

	http.HandleFunc("/Register/{service}/{port}", handlers.Service_Register)
	http.HandleFunc("/", handlers.Service_Get_Names)
	http.HandleFunc("/Service", handlers.Service_Get_Names)
	http.HandleFunc("/Service/{service}", handlers.Service_Get)

	fmt.Printf("Listening on %s\n", ":"+port_str)
	err := http.ListenAndServeTLS(
		":"+port_str,
		"../server.crt",
		"../server.key",
		nil,
	)
	if err != nil {
		log.Println(err)
	}
}
