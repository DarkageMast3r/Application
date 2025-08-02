package handlers

import (
	m "Financiering/Models"
	r "Financiering/Repositories"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func AddDossier(wr http.ResponseWriter, rq *http.Request) {
	rq.ParseForm()
	clientID, err := strconv.Atoi(rq.Form.Get("ClientID"))
	if err != nil {
		fmt.Println("Failed to stringConvert: ", err)
	}
	zorgTechID, err := strconv.Atoi(rq.Form.Get("ZorgTechID"))
	if err != nil {
		fmt.Println("Failed to stringConvert: ", err)
	}
	
	_, err = r.InsertDossier(clientID, zorgTechID)
	if err != nil {
		log.Println("AddDossier: ", err)
		wr.WriteHeader(http.StatusInternalServerError)
	}
	HomePageHandler(wr, rq)
}

func RemoveDossier(wr http.ResponseWriter, rq *http.Request) {
	DossierID, err := strconv.Atoi(rq.PathValue("dossierID"))
	if err != nil {
		log.Println("RemoveDossier: ", err)
		wr.WriteHeader(http.StatusInternalServerError)
	}
	err = r.RemoveDossier(DossierID)
	if err != nil {
		log.Println("RemoveDossier: ", err)
		wr.WriteHeader(http.StatusInternalServerError)
	}
	HomePageHandler(wr, rq)
}

func AddBudget(wr http.ResponseWriter, rq *http.Request) {
	var Budget m.Budget
	rq.ParseForm()
	maxBedrag, err := strconv.Atoi(rq.Form.Get("MaxBedrag"))
	if err != nil {
		log.Println("AddBudget(strconv): ", err)
		wr.WriteHeader(http.StatusInternalServerError)
	}
	Budget.MaxBedrag = float64(maxBedrag)
	DossierID, err := strconv.Atoi(rq.PathValue("dossierID"))
	if err != nil {
		log.Println("AddBudget(strconv)2: ", err)
		wr.WriteHeader(http.StatusInternalServerError)
	}

	Budget.NewBudget()
	err = r.ConnectDossier(Budget.ID, DossierID)
	if err != nil {
		log.Println("ConnectDossier failed: ", err)
	}
	HomePageHandler(wr, rq)
}