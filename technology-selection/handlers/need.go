package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"service/models"
	"service/repository"
	"service/viewModels"
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

func Need_Update(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	var need models.Need
	err := decoder.Decode(&need, r.Form)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = repository.Need_Save(&need)
	if err != nil {
		fmt.Println(err)
	}
	http.Redirect(w, r, "/View/Need", http.StatusSeeOther)
}

func Need_Create(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	var need models.Need
	err := decoder.Decode(&need, r.Form)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = repository.Need_Save(&need)
	if err != nil {
		fmt.Println(err)
	}
	http.Redirect(w, r, "/View/Need", http.StatusSeeOther)
}
func Need_Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	repository.Need_Delete(id)
	http.Redirect(w, r, "/View/Need", http.StatusSeeOther)
}

func Need_View(w http.ResponseWriter, r *http.Request) {
	err := Template_View(w, repository.Need_Get_All(), "need/view", "templates/need/view.gohtml")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func Need_View_Create(w http.ResponseWriter, r *http.Request) {
	var view viewModels.NeedCreate
	view.Technologies = repository.Tech_Get_All()

	err := Template_View(w, view, "need/create", "templates/need/create.gohtml")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func Need_View_Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var view viewModels.NeedUpdate
	view.Need = repository.Need_Get_By_Id(id)
	view.Technologies = repository.Tech_Get_All()

	err = Template_View(w, view, "need/update", "templates/need/update.gohtml")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
