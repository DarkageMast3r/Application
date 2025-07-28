package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func AddDossier(wr http.ResponseWriter, rq *http.Request) {
	rq.ParseForm()
	clientid, err := strconv.Atoi(rq.Form.Get("ClientID"))
	if err != nil {
		fmt.Println("Failed to stringConvert: ", err)
	}
	zorgtechid, err := strconv.Atoi(rq.Form.Get("ZorgTechID"))
	if err != nil {
		fmt.Println("Failed to stringConvert: ", err)
	}

	var dossier m.financieringsdossier
	err = dossier.NieuwDossier(clientid, zorgtechid)
	if err != nil {
		log.Println("AddDossier: ", err)
		wr.WriteHeader(http.StatusInternalServerError)
	}
	HomePageHandler(wr, rq)
}

func RemoveDossier(wr http.ResponseWriter, rq *http.Request) {
	DossierID, err := strconv.Atoi(rq.PathValue("dossierID"))
	if err != nil {
		log.Println("AddDossier: ", err)
		wr.WriteHeader(http.StatusInternalServerError)
	}
	err = r.RemoveDossier(DossierID)
	if err != nil {
		log.Println("AddDossier: ", err)
		wr.WriteHeader(http.StatusInternalServerError)
	}
	AddorRemovePageHandler(wr, rq)
}
