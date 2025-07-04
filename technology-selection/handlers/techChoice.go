package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"service/models"
	"service/repository"
	"strconv"
)

func TechChoice_Get_All(w http.ResponseWriter, r *http.Request) {
	result, err := json.Marshal(repository.TechChoice_Get_All())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.Write(result)
}

func TechChoice_Get_By_Id(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	techChoice := repository.TechChoice_Get_By_Id(id)
	if techChoice == nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("{}"))
	}

	result, err := json.Marshal(techChoice)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.Write(result)
}

func TechChoice_Choose(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	techChoice := repository.TechChoice_Get_By_Id(id)
	if techChoice == nil || techChoice.Status != models.SelectionStatus_Shortlist {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	r.ParseForm()
	techChoice.Reasoning.String = r.Form.Get("reasoning")
	techChoice.Status = models.SelectionStatus_Chosen
	err = repository.TechChoice_Save(techChoice)
	if err != nil {
		fmt.Println(err)
		return
	}

	Shortlist_View(w, r)

	clientCase := repository.Case_Get_By_Id(techChoice.CaseId)
	request := make(map[string]string)
	request["clientId"] = clientCase.ClientId
	request["zorgtechId"] = strconv.Itoa(techChoice.TechId)
	jsonBody, err := json.Marshal(&request)
	if err != nil {
		fmt.Println(err)
		return
	}

	http.Post("localhost/implementation/api/v1/imlementatie/aanvraag", "text/json", bytes.NewReader(jsonBody))
}

func TechChoice_Reject(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	techChoice := repository.TechChoice_Get_By_Id(id)
	if techChoice == nil || techChoice.Status != models.SelectionStatus_Shortlist {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	r.ParseForm()
	techChoice.Reasoning.String = r.Form.Get("reasoning")
	techChoice.Status = models.SelectionStatus_Rejected
	err = repository.TechChoice_Save(techChoice)
	if err != nil {
		fmt.Println(err)
	}
	Shortlist_View(w, r)
}
