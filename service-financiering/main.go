package main

import (
	h "Financiering/Handlers"
	"fmt"
	"net/http"
)

func main() {
	s := &http.Server{
		Addr:    ":8080",
		Handler: http.DefaultServeMux,
	}

	http.HandleFunc("GET /", h.HomeHandler)

	fmt.Println("o7")
	s.ListenAndServe()
}
