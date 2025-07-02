package main

import (
	h "Financiering/Handlers"
	r "Financiering/Repositories"
	"fmt"
	"net/http"
)

func main() {
	s := &http.Server{
		Addr:    ":8080",
		Handler: http.DefaultServeMux,
	}
	r.Database_Get()

	http.HandleFunc("GET /Finance", h.HomeHandler)

	fmt.Println("o7")
	s.ListenAndServe()
}
