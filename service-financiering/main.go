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

	http.HandleFunc("GET /Finance", h.HomePageHandler)
	http.HandleFunc("GET /Finance/Add", h.AddorRemovePageHandler)
	http.HandleFunc("Get /Finance/{dossierID}", h.HomePageHandler) //non existed handler(for now)

	http.HandleFunc("POST /Finance/Add", h.AddDossier)

	fmt.Println("o7")
	s.ListenAndServe()
}
