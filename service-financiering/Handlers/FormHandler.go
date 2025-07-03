package handlers

import (
	r "Financiering/Repositories"
	"fmt"
	"net/http"
)

func AddDossier(wr http.ResponseWriter, rq *http.Request) {
	db := r.Database_Get()

	rq.ParseForm()
	clientid := rq.Form.Get("ClientID")
	zorgtechid := rq.Form.Get("ZorgTechID")

	_, err := db.Query("INSERT INTO financieringsdossier(ClientID, ZorgTechID) VALUES(?,?)", clientid, zorgtechid)
	if err != nil {
		fmt.Println("AddDossier: ", err)
		wr.WriteHeader(http.StatusBadRequest)
	}
	//return to homepage
	HomePageHandler(wr, rq)
}
