package main

import (
	"fmt"
	"log"
	"net/http"
	"service/handlers"
	"service/service"
	"strconv"
)

func main() {
	service.Init()

	http.HandleFunc("/Category", handlers.Category_Get_All)
	http.HandleFunc("/Category/{id}", handlers.Category_Get_By_Id)

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
