package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"service/models"
	"service/repository"
	"service/util"
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
func TechChoice_Create(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		var need models.TechChoice
		util.Crud_View_Create(w, reflect.TypeOf(need), "/TechChoice/Create")
		return
	}
	r.ParseForm()

	var techChoice models.TechChoice
	err := util.Fill_Fields_From_Form(&techChoice, r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = repository.TechChoice_Save(&techChoice)
	if err != nil {
		fmt.Println(err)
	}
	w.WriteHeader(http.StatusOK)
}
