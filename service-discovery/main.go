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

func hostIsAlive(uri string) bool {
	result, err := http.Get("https://" + uri + "/Category")
	if err != nil {
		log.Fatal(err)
		return false
	}
	if result != nil {
		result.Body.Close()
	}
	return true
}

func main() {
	services := make(map[string]Service)
	config := readConfig("../config.json")

	go func() {
		for {
			for serviceName, service := range services {
				for i, host := range service.Hosts {
					if hostIsAlive(host) {
						continue
					}
					fmt.Printf("Disconnecting host %s\n", host)
					service.Hosts[i] = service.Hosts[len(service.Hosts)-1]
					service.Hosts = service.Hosts[:len(service.Hosts)-1]
				}
				services[serviceName] = service
			}
		}
	}()

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
		fmt.Printf("Registered service %s at %s\n", serviceName, serviceUri)
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
	http.HandleFunc("/Get", func(w http.ResponseWriter, r *http.Request) {
		serviceNames := make([]string, len(services))
		i := 0
		for name := range services {
			serviceNames[i] = name
			i++
		}
		result, err := json.Marshal(serviceNames)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(result)
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
