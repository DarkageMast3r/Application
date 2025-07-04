package handlers

import (
	"encoding/json"
	"net/http"
	"service/models"
	"service/repository"
	"service/service"
	"service/viewModels"
	"strconv"

	"github.com/gorilla/schema"
)

var decoder = schema.NewDecoder()

func Tech_Get_All(w http.ResponseWriter, r *http.Request) {
	result, err := json.Marshal(repository.Tech_Get_All())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.Write(result)
}

func Tech_Get_By_Id(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	tech := repository.Tech_Get_By_Id(id)
	if tech == nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("{}"))
	}

	result, err := json.Marshal(tech)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.Write(result)
}

func Tech_Update(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	var tech models.Tech
	err := decoder.Decode(&tech, r.Form)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = repository.Tech_Save(&tech)
	if err != nil {
		service.LogError(err)
	}
	http.Redirect(w, r, "/View/Tech", http.StatusSeeOther)
}

func Tech_Create(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	var tech models.Tech
	err := decoder.Decode(&tech, r.Form)
	if err != nil {
		service.LogWarning("Could not decode: ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = repository.Tech_Save(&tech)
	if err != nil {
		service.LogWarning(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/View/Tech", http.StatusSeeOther)
}

func Tech_Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	repository.Tech_Delete(id)
	http.Redirect(w, r, "/View/Tech", http.StatusSeeOther)
}

func Tech_View(w http.ResponseWriter, r *http.Request) {
	err := Template_View(w, repository.Tech_Get_All(), "tech/view", "templates/tech/view.gohtml")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func Tech_View_Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var view viewModels.TechUpdate
	view.Tech = repository.Tech_Get_By_Id(id)
	view.Categories = repository.Category_Get_All()
	view.Needs = repository.Need_Get_All()
	err = Template_View(w, view, "tech/update", "templates/tech/update.gohtml")
	if err != nil {
		service.LogError(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func Tech_View_Create(w http.ResponseWriter, r *http.Request) {
	var view viewModels.TechCreate
	view.Categories = repository.Category_Get_All()
	view.Needs = repository.Need_Get_All()

	err := Template_View(w, view, "tech/create", "templates/tech/create.gohtml")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func Tech_Shortlist(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		service.LogError(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if repository.Tech_Get_By_Id(id) == nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	r.ParseForm()
	var selectForm viewModels.SelectTechnology
	err = decoder.Decode(&selectForm, r.Form)

	var techChoice models.TechChoice
	techChoice.TechId = id
	techChoice.CaseId = selectForm.Case.Id
	techChoice.Status = models.SelectionStatus_Shortlist
	if err != nil {
		service.LogError(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = repository.TechChoice_Save(&techChoice)
	if err != nil {
		service.LogError(err)
	}
	Selection_View(w, r)
}
