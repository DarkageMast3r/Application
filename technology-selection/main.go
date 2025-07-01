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
	http.HandleFunc("/Category/{id}/Update", handlers.Category_Update)
	http.HandleFunc("/Category/{id}/Delete", handlers.Category_Delete)
	http.HandleFunc("/Category/Create", handlers.Category_Create)
	http.HandleFunc("/View/Category", handlers.Category_View)
	http.HandleFunc("/View/Category/{id}/Update", handlers.Category_View_Update)
	http.HandleFunc("/View/Category/Create", handlers.Category_View_Create)

	http.HandleFunc("/Need", handlers.Need_Get_All)
	http.HandleFunc("/Need/{id}", handlers.Need_Get_By_Id)
	http.HandleFunc("/Need/{id}/Update", handlers.Need_Update)
	http.HandleFunc("/Need/{id}/Delete", handlers.Need_Delete)
	http.HandleFunc("/Need/Create", handlers.Need_Create)
	http.HandleFunc("/View/Need", handlers.Need_View)
	http.HandleFunc("/View/Need/{id}/Update", handlers.Need_View_Update)
	http.HandleFunc("/View/Need/Create", handlers.Need_View_Create)

	http.HandleFunc("/Tech", handlers.Tech_Get_All)
	http.HandleFunc("/Tech/{id}", handlers.Tech_Get_By_Id)
	http.HandleFunc("/Tech/{id}/Update", handlers.Tech_Update)
	http.HandleFunc("/Tech/{id}/Delete", handlers.Tech_Delete)
	http.HandleFunc("/Tech/Create", handlers.Tech_Create)
	http.HandleFunc("/View/Tech", handlers.Tech_View)
	http.HandleFunc("/View/Tech/{id}/Update", handlers.Tech_View_Update)
	http.HandleFunc("/View/Tech/Create", handlers.Tech_View_Create)

	http.HandleFunc("/Tech/{id}/Shortlist", handlers.Tech_Shortlist)
	http.HandleFunc("/TechChoice/{id}/Choose", handlers.TechChoice_Choose)
	http.HandleFunc("/TechChoice/{id}/Reject", handlers.TechChoice_Reject)

	http.HandleFunc("/TechChoice", handlers.TechChoice_Get_All)
	http.HandleFunc("/TechChoice/{id}", handlers.TechChoice_Get_By_Id)

	http.HandleFunc("/", handlers.Start_View)
	http.HandleFunc("/Select", handlers.Selection_View)
	http.HandleFunc("/Shortlist", handlers.Shortlist_View)

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
