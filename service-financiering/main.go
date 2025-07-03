package main

import (
	h "Financiering/Handlers"
	r "Financiering/Repositories"
	"Financiering/service"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func main() {
	r.Database_Get()

	http.HandleFunc("GET /Finance", h.HomePageHandler)
	http.HandleFunc("GET /Finance/Add", h.AddorRemovePageHandler)
	http.HandleFunc("Get /Finance/{dossierID}", h.HomePageHandler) //non existed handler(for now)

	http.HandleFunc("POST /Finance/Add", h.AddDossier)

	fmt.Println("o7")

	service.Init()
	port := service.Register("financing", http.DefaultServeMux.ServeHTTP)
	log.Println("Listening on port", port)
	http.ListenAndServeTLS(
		":"+strconv.Itoa(port),
		"../server.crt",
		"../server.key",
		nil,
	)
}
