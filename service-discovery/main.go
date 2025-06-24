package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Service_discovery_root string `json:"service_discovery_root"`
	Service_discovery_port int    `json:"service_discovery_port"`
	Allow_insecure         bool   `json:"allow_insecure"`
}

type Service struct {
	Hosts      []string
	LastServed int
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
	services := make(map[string]Service)
	config := readConfig("../config.json")

	http.HandleFunc("/Register/{service}/{port}", func(w http.ResponseWriter, r *http.Request) {
		serviceName := r.PathValue("service")
		servicePort := r.PathValue("port")
		idx := strings.LastIndex(r.RemoteAddr, ":")
		serviceUri := r.RemoteAddr[:idx] + ":" + servicePort

		service, exists := services[serviceName]
		if !exists {
			service.LastServed = 0
			service.Hosts = make([]string, 0)
		}
		for _, current := range service.Hosts {
			if current == serviceUri {
				return
			}
		}
		service.Hosts = append(service.Hosts, serviceUri)
		services[serviceName] = service
	})

	http.HandleFunc("/Get/{service}", func(w http.ResponseWriter, r *http.Request) {
		serviceName := r.PathValue("service")
		service, exists := services[serviceName]
		if !exists {
			http.NotFound(w, r)
			return
		}
		service.LastServed = (service.LastServed + 1) % len(service.Hosts)
		services[serviceName] = service
		io.WriteString(w, service.Hosts[service.LastServed])
	})
	fmt.Printf("Listening on %s\n", ":"+strconv.Itoa(config.Service_discovery_port))
	err := http.ListenAndServeTLS(
		":"+strconv.Itoa(config.Service_discovery_port),
		"../server.crt",
		"../server.key",
		nil,
	)
	if err != nil {
		log.Println(err)
	}
}
