package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"test_service_2/service"
)

func main() {
	service.Init()

	port := service.Register("test-service-2")
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		service.Queue_Write("test", []byte("Hello from test-service-2"), "text/plain")
	})

	fmt.Printf("Listening on %s\n", ":"+strconv.Itoa(port))
	err := http.ListenAndServeTLS(
		":"+strconv.Itoa(port),
		"../server.crt",
		"../server.key",
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}
}
