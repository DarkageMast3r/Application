package main

import (
	"fmt"
	"net/http"
	"service/global"
	"service/handlers"
	"service/repository"
	"service/service"
	"strconv"
)

func is_authorised(r *http.Request) bool {
	api := global.Config.Service_discovery_root + ":" + strconv.Itoa(global.Config.Service_discovery_port) + "/Auth/api/v1/auth"
	service.LogInfo("API", api, "is not yet implemented, freely permitting authorisation")
	// TODO: Authorisation is not yet reliably implemented, make API call once possible.
	return true
}

func authorize(handler func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if !is_authorised(r) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		handler(w, r)
	}
}

func main() {
	service.Init()
	global.Init()

	repository.Database_Get()

	http.HandleFunc("/Case/Create", handlers.Case_Create)

	http.HandleFunc("/Category", authorize(handlers.Category_Get_All))
	http.HandleFunc("/Category/{id}", authorize(handlers.Category_Get_By_Id))
	http.HandleFunc("/Category/{id}/Update", authorize(handlers.Category_Update))
	http.HandleFunc("/Category/{id}/Delete", authorize(handlers.Category_Delete))
	http.HandleFunc("/Category/Create", authorize(handlers.Category_Create))
	http.HandleFunc("/View/Category", authorize(handlers.Category_View))
	http.HandleFunc("/View/Category/{id}/Update", authorize(handlers.Category_View_Update))
	http.HandleFunc("/View/Category/Create", authorize(handlers.Category_View_Create))

	http.HandleFunc("/Need", authorize(handlers.Need_Get_All))
	http.HandleFunc("/Need/{id}", authorize(handlers.Need_Get_By_Id))
	http.HandleFunc("/Need/{id}/Update", authorize(handlers.Need_Update))
	http.HandleFunc("/Need/{id}/Delete", authorize(handlers.Need_Delete))
	http.HandleFunc("/Need/Create", authorize(handlers.Need_Create))
	http.HandleFunc("/View/Need", authorize(handlers.Need_View))
	http.HandleFunc("/View/Need/{id}/Update", authorize(handlers.Need_View_Update))
	http.HandleFunc("/View/Need/Create", authorize(handlers.Need_View_Create))

	http.HandleFunc("/Tech", authorize(handlers.Tech_Get_All))
	http.HandleFunc("/Tech/{id}", authorize(handlers.Tech_Get_By_Id))
	http.HandleFunc("/Tech/{id}/Update", authorize(handlers.Tech_Update))
	http.HandleFunc("/Tech/{id}/Delete", authorize(handlers.Tech_Delete))
	http.HandleFunc("/Tech/Create", authorize(handlers.Tech_Create))
	http.HandleFunc("/View/Tech", authorize(handlers.Tech_View))
	http.HandleFunc("/View/Tech/{id}/Update", authorize(handlers.Tech_View_Update))
	http.HandleFunc("/View/Tech/Create", authorize(handlers.Tech_View_Create))

	http.HandleFunc("/Tech/{id}/Shortlist", authorize(handlers.Tech_Shortlist))
	http.HandleFunc("/TechChoice/{id}/Choose", authorize(handlers.TechChoice_Choose))
	http.HandleFunc("/TechChoice/{id}/Reject", authorize(handlers.TechChoice_Reject))

	http.HandleFunc("/TechChoice", authorize(handlers.TechChoice_Get_All))
	http.HandleFunc("/TechChoice/{id}", authorize(handlers.TechChoice_Get_All_By_Case_Id))

	http.HandleFunc("/", authorize(handlers.Start_View))
	http.HandleFunc("/Select", authorize(handlers.Selection_View))
	http.HandleFunc("/Shortlist", authorize(handlers.Shortlist_View))

	service.Register("selection", http.DefaultServeMux.ServeHTTP)

	fmt.Println("Service technology selection started.")
	forever := make(chan struct{})
	<-forever
}
