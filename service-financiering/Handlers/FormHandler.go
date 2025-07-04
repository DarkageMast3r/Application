package handlers

import (
	r "Financiering/Repositories"
	"fmt"
	"net/http"
	"strconv"
)

func AddDossier(wr http.ResponseWriter, rq *http.Request) {
	rq.ParseForm()
	clientid , err := strconv.Atoi(rq.Form.Get("ClientID"))
	if err != nil {
		fmt.Println("Failed to stringConvert: ", err)
	}
	zorgtechid, err := strconv.Atoi(rq.Form.Get("ZorgTechID"))
	if err != nil {
		fmt.Println("Failed to stringConvert: ", err)
	}

	err = r.InsertDossier(clientid, zorgtechid)
	if err != nil {
		fmt.Println("AddDossier: ", err)
		wr.WriteHeader(http.StatusInternalServerError)
	}
	HomePageHandler(wr, rq)
}
