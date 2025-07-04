package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"service/models"
	"service/repository"
	"service/service"
	"strconv"
)

func Case_Create(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		service.LogWarning("Read body:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	clientCase := models.Case{}
	err = json.Unmarshal(body, &clientCase)
	if err != nil {
		service.LogWarning("Case decode:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = repository.Case_Save(&clientCase)
	if err != nil {
		service.LogWarning("Case save:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(strconv.Itoa(clientCase.Id)))
}
