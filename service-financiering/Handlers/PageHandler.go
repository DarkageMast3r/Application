package handlers

import (
	m "Financiering/Models"
	"fmt"
	"net/http"
	"strconv"
)

func HomePageHandler(wr http.ResponseWriter, rq *http.Request) {
	LoadTemplate(wr, "Templates/Home.gohtml", m.GetAllDossiers())
}

func AddorRemovePageHandler(wr http.ResponseWriter, rq *http.Request) {
	LoadTemplate(wr, "Templates/AddorRemove.gohtml", nil)
}

func DossierPageHandler(wr http.ResponseWriter, rq *http.Request) {
	DossierID, err := strconv.Atoi(rq.PathValue("dossierID"))
	if err != nil {
		fmt.Println("DossierPageHandler", err)
		return
	}
	LoadTemplate(wr, "Templates/DetailDossier.gohtml", m.GetDossierbyID(DossierID))
}
