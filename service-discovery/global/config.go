package global

import (
	"encoding/json"
	"io"
	"log"
	"os"
)

type Configuration struct {
	Service_discovery_root string   `json:"service_discovery_root"`
	Service_discovery_port int      `json:"service_discovery_port"`
	Allow_insecure         bool     `json:"allow_insecure"`
	Queues                 []string `json:"queues"`
}

var Config Configuration = Configuration{}

func ReadConfig(path string) Configuration {
	jsonFile, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer jsonFile.Close()

	byteValue, _ := io.ReadAll(jsonFile)
	var config Configuration
	json.Unmarshal(byteValue, &config)
	Config = config
	return config
}
