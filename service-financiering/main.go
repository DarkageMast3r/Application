package main

import (
	h "Financiering/Handlers"
	r "Financiering/Repositories"
	"Financiering/service"
	"fmt"
	"net/http"
)

func main() {
	service.Init()
	r.Database_Get()

	http.HandleFunc("GET /", h.HomePageHandler)
	http.HandleFunc("GET /Add", h.AddorRemovePageHandler)
	http.HandleFunc("GET /{dossierID}", h.DossierPageHandler)
	http.HandleFunc("Get /{dossierID}/Remove", h.RemoveDossier)

	http.HandleFunc("POST /Add", h.AddDossier)

	service.Register("financing", http.DefaultServeMux.ServeHTTP)

	fmt.Println("Service financing started")
	forever := make(chan struct{})
	<-forever
}
