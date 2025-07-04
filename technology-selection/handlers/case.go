package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"service/models"
	"service/repository"
)

func Case_Create(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println("Read body:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	clientCase := models.Case{}
	err = json.Unmarshal(body, &clientCase)
	if err != nil {
		fmt.Println("Case decode:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = repository.Case_Save(&clientCase)
	if err != nil {
		fmt.Println("Case save:", err)
	}
	w.WriteHeader(http.StatusOK)
}
