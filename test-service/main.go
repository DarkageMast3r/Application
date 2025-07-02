package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"test-service/service"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "Hello world from test-service - 2")
	})
	service.Init()
	port := service.Register("Test")
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
