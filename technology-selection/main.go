package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"main/models"
	"net"
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

func getLocalPort() int {
	conn, err := net.Dial("udp", "0.0.0.0:80")
	if err != nil {
		log.Fatal(err)
	}
	return conn.LocalAddr().(*net.UDPAddr).Port
}

func registerService(config Config, name string, port int) {
	host_root := config.Service_discovery_root + ":" + strconv.Itoa(config.Service_discovery_port)
	http.Get(host_root + "/Register/" + name + "/" + strconv.Itoa(port))
}

func main() {
	config := readConfig("../config.json")
	if config.Allow_insecure {
		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}

	var test models.Category
	test.Id = 5
	test.Id += 3
	models.Category_Init()

	http.HandleFunc("/Category", func(w http.ResponseWriter, r *http.Request) {
		result, err := json.Marshal(*models.Category_Get_All())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		w.Write(result)
	})
	http.HandleFunc("/Category/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		category := models.Category_Get_By_Id(id)
		if category == nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("{}"))
		}

		result, err := json.Marshal(category)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		w.Write(result)
	})

	port := getLocalPort()
	registerService(config, "technology-selection", port)

	fmt.Printf("Listening on %s\n", ":"+strconv.Itoa(port))
	err := http.ListenAndServeTLS(
		":"+strconv.Itoa(port),
		"../server.crt",
		"../server.key",
		nil,
	)
	if err != nil {
		log.Println(err)
	}
}
