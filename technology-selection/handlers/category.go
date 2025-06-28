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

func Category_Create(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		var category models.Category
		util.Crud_View_Create(w, reflect.TypeOf(category), "/Category/Create")
		return
	}
	r.ParseForm()
	var category models.Category
	err := util.Fill_Fields_From_Form(&category, r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = repository.Category_Save(&category)
	if err != nil {
		fmt.Println(err)
	}
	w.WriteHeader(http.StatusOK)
}
