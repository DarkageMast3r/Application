package main

import (
	"io"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "Hello world!\n")
	})
	http.HandleFunc("/Test", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "Hello world, test!\n")
	})
	if err := http.ListenAndServeTLS(":443", "../server.crt", "../server.key", nil); err != nil {
		log.Println(err)
	}
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Println(err)
	}
	log.Println("Server started")
}
