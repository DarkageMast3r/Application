package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
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

func httpGet(url string) (error, string) {
	resp, err := http.Get(url)
	if err != nil {
		return err, ""
	}
	body, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return err, ""
	}
	return nil, string(body)
}

func main() {
	config := readConfig("../config.json")
	if config.Allow_insecure {
		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}
	host_root := config.Service_discovery_root + ":" + strconv.Itoa(config.Service_discovery_port)
	fmt.Printf("Connecting to %s\n", host_root)
	err, serviceUri := httpGet(host_root + "/Get/Test")
	if err != nil {
		log.Fatal(err)
	}
	err, message := httpGet("https://" + serviceUri)
	fmt.Print(message)
}
