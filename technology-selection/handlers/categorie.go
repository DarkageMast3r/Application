package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"service/models"
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

func Category_Create_View(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("<html><head><title>Title</title></head><body>"))
	w.Write([]byte("<form action=\"/Category/Create\" method=\"post\">"))
	w.Write([]byte("<label>Name<input type=\"text\" name=\"name\"/></label>"))
	w.Write([]byte("<label>Description<input type=\"text\" name=\"description\"/></label>"))
	w.Write([]byte("<input type=\"submit\"/>"))
	w.Write([]byte("</form>"))
	w.Write([]byte("</body>"))
	w.Header().Add("content-type", "text/html")
	w.WriteHeader(http.StatusOK)
}

func Category_Create(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		Category_Create_View(w, r)
		return
	}
	r.ParseForm()
	fmt.Print(r.Form, r.Form.Has("name"), r.Form.Has("description"))
	if !r.Form.Has("name") || !r.Form.Has("description") {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var category models.Category
	category.Naam = r.Form.Get("name")
	category.Beschrijving = r.Form.Get("description")
	err := repository.Category_Save(&category)
	if err != nil {
		fmt.Println(err)
	}
	w.WriteHeader(http.StatusOK)
}
