package handlers

import (
	r "Financiering/Repositories"
	"net/http"
)

func HomeHandler(wr http.ResponseWriter, rq *http.Request) {
	LoadTemplate(wr, "Templates/Home.gohtml", r.ReadTable())
	// maybe load some extra data, for now it works
}
