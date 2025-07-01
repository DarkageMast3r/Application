package handlers

import (
	"net/http"
)

func HomeHandler(wr http.ResponseWriter, rq *http.Request) {
	LoadTemplate(wr, "Templates/Home.gohtml", "Training data, pls change")
	// maybe load some extra data, for now it works
}
