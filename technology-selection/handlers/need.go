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

func Need_Get_All(w http.ResponseWriter, r *http.Request) {
	result, err := json.Marshal(repository.Need_Get_All())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.Write(result)
}

func Need_Get_By_Id(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	need := repository.Need_Get_By_Id(id)
	if need == nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("{}"))
	}

	result, err := json.Marshal(need)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.Write(result)
}

func Need_Create_View(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("<html><head><title>Create Need</title></head><body>"))
	w.Write([]byte("<form action=\"/Need/Create\" method=\"post\">"))
	w.Write([]byte("<label>Name<input type=\"text\" name=\"name\"/></label>"))
	w.Write([]byte("<label>Source<input type=\"text\" name=\"source\"/></label>"))
	w.Write([]byte("<input type=\"submit\"/>"))
	w.Write([]byte("</form>"))
	w.Write([]byte("</body>"))
	w.Header().Add("content-type", "text/html")
	w.WriteHeader(http.StatusOK)
}

func Need_Create(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		var need models.Need
		util.Crud_View_Create(w, reflect.TypeOf(need), "/Need/Create")
		return
	}
	r.ParseForm()

	var need models.Need
	err := util.Fill_Fields_From_Form(&need, r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = repository.Need_Save(&need)
	if err != nil {
		fmt.Println(err)
	}
	w.WriteHeader(http.StatusOK)
}
