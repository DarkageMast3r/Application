package handlers

import (
	r "Financiering/Repositories"
	"net/http"
)

func HomePageHandler(wr http.ResponseWriter, rq *http.Request) {
	LoadTemplate(wr, "Templates/Home.gohtml", r.GetDossiers())
	// maybe load some extra data, for now it works
}

func AddorRemovePageHandler(wr http.ResponseWriter, rq *http.Request) {
	LoadTemplate(wr, "Templates/AddorRemove.gohtml", nil)
	// maybe load some extra data, for now it works
}
