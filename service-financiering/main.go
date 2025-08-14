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
	http.HandleFunc("Get /Remove/{dossierID}", h.RemoveDossier)
	http.HandleFunc("Get /GetBudget/{ClientID}", h.ViewBudgetPage)

	http.HandleFunc("POST /Add", h.AddDossier)
	http.HandleFunc("POST /AddBudget/{dossierID}", h.AddBudget)

	service.Register("financing", http.DefaultServeMux.ServeHTTP)

	fmt.Println("Service financing started")
	forever := make(chan struct{})
	<-forever
}
