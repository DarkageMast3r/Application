package handlers

import (
	// m "Financiering/Models"
	r "Financiering/Repositories"
	"fmt"
	"net/http"
	// "strconv"
)

func AddDossier(wr http.ResponseWriter, rq *http.Request) {
	db := r.Database_Get()

	rq.ParseForm()
	clientid := rq.Form.Get("ClientID")
	zorgtechid := rq.Form.Get("ZorgTechID")

	_, err := db.Query("INSERT INTO financieringsdossier(ClientID, ZorgTechID) VALUES(?,?)", clientid, zorgtechid)
	if err != nil {
		fmt.Println("AddDossier2: ", err)
		return
	}
	//return to homepage
	HomePageHandler(wr, rq)
}
