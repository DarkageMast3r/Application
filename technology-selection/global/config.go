package global

import (
	"encoding/json"
	"io"
	"os"
	"service/service"
)

type Configuration struct {
	Service_discovery_root string `json:"service_discovery_root"`
	Service_discovery_port int    `json:"service_discovery_port"`
}

var Config Configuration

func readConfig(path string) Configuration {
	jsonFile, err := os.Open(path)
	if err != nil {
		service.LogError(err)
	}
	defer jsonFile.Close()

	byteValue, _ := io.ReadAll(jsonFile)
	var config Configuration
	json.Unmarshal(byteValue, &config)
	return config
}

func Init() {
	Config = readConfig("config.json")
}
