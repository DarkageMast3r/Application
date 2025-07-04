package handlers

import (
	"encoding/json"
	"net/http"
	"service/models"
	"service/repository"
	"service/service"
	"service/viewModels"
	"strconv"
)

func Category_Get_All(w http.ResponseWriter, r *http.Request) {
	result, err := json.Marshal(repository.Category_Get_All())
	if err != nil {
		service.LogError(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.Write(result)
}

func Category_Get_By_Id(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		service.LogError(err)
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
		service.LogError(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.Write(result)
}

func Category_Update(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	var category models.Category
	err := decoder.Decode(&category, r.Form)
	if err != nil {
		service.LogError(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = repository.Category_Save(&category)
	if err != nil {
		service.LogError(err)
	}
	http.Redirect(w, r, "/View/Category", http.StatusSeeOther)
}

func Category_Create(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	var category models.Category
	err := decoder.Decode(&category, r.Form)
	if err != nil {
		service.LogError(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = repository.Category_Save(&category)
	if err != nil {
		service.LogError(err)
	}
	http.Redirect(w, r, "/View/Category", http.StatusSeeOther)
}
func Category_Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		service.LogError(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	repository.Category_Delete(id)
	http.Redirect(w, r, "/View/Category", http.StatusSeeOther)
}

func Category_View(w http.ResponseWriter, r *http.Request) {
	err := Template_View(w, repository.Category_Get_All(), "category/view", "templates/category/view.gohtml")
	if err != nil {
		service.LogError(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func Category_View_Create(w http.ResponseWriter, r *http.Request) {
	var view viewModels.CategoryCreate

	err := Template_View(w, view, "category/create", "templates/category/create.gohtml")
	if err != nil {
		service.LogError(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func Category_View_Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		service.LogError(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var view viewModels.CategoryUpdate
	view.Category = repository.Category_Get_By_Id(id)

	err = Template_View(w, view, "category/update", "templates/category/update.gohtml")
	if err != nil {
		service.LogError(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}
