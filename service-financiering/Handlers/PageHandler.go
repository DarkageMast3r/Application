package handlers

import (
	r "Financiering/Repositories"
	"net/http"
	"strconv"
	"fmt"
)

func HomePageHandler(wr http.ResponseWriter, rq *http.Request) {
	LoadTemplate(wr, "Templates/Home.gohtml", r.GetDossiers())
	// maybe load some extra data, for now it works
}

func AddorRemovePageHandler(wr http.ResponseWriter, rq *http.Request) {
	LoadTemplate(wr, "Templates/AddorRemove.gohtml", nil)
	// maybe load some extra data, for now it works
}

func DossierPageHandler(wr http.ResponseWriter, rq *http.Request) {
	DossierID, err := strconv.Atoi(rq.PathValue("dossierID"))
	if err != nil {
		fmt.Println("DossierPageHandler", err)
		return
	}
	LoadTemplate(wr, "Templates/DetailDossier.gohtml", r.GetDossierbyID(DossierID))
}
