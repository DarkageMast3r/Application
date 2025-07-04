package main

import (
	h "Financiering/Handlers"
	r "Financiering/Repositories"
	"Financiering/service"
	"net/http"
)

func main() {
	r.Database_Get()

	http.HandleFunc("GET /Finance", h.HomePageHandler)
	http.HandleFunc("GET /Finance/Add", h.AddorRemovePageHandler)
	http.HandleFunc("GET /Finance/{dossierID}", h.DossierPageHandler)

	http.HandleFunc("POST /Finance/Add", h.AddDossier)

	service.Init()
	forever := make(chan struct{})
	<-forever
}
