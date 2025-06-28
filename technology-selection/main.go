package main

import (
	"fmt"
	"log"
	"net/http"
	"service/handlers"
	"service/repository"
	"service/service"
	"strconv"
)

func main() {
	service.Init()
	repository.Database_Get()
	http.HandleFunc("/Category", handlers.Category_Get_All)
	http.HandleFunc("/Category/{id}", handlers.Category_Get_By_Id)
	http.HandleFunc("/Category/Create", handlers.Category_Create)

	port := service.Register("technology-selection")

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
