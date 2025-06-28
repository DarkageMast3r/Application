package handlers

import (
	"encoding/json"
	"net/http"
	"service/repository"
	"strconv"
)

func Category_Get_All(w http.ResponseWriter, r *http.Request) {
	result, err := json.Marshal(repository.Category_Get_All())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.Write(result)
}

func Category_Get_By_Id(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	category := repository.Category_Get_By_Id(id)
	if category == nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("{}"))
	}

	result, err := json.Marshal(category)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.Write(result)
}
