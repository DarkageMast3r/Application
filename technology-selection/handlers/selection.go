package handlers

import (
	"fmt"
	"net/http"
	"service/models"
	"service/repository"
	"service/viewModels"
	"strconv"
)

func tech_has_need(tech *models.Tech, need *models.CheckableOption) bool {
	needId, err := strconv.Atoi(need.Name)
	if err != nil {
		return false
	}
	for _, techNeed := range tech.Needs {
		if techNeed.Id == needId {
			return true
		}
	}
	return false
}

func tech_has_all_needs(tech *models.Tech, needs []models.CheckableOption) bool {
	for _, need := range needs {
		if !need.Selected {
			continue
		}
		if !tech_has_need(tech, &need) {
			return false
		}
	}
	return true
}

func Start_View(w http.ResponseWriter, r *http.Request) {
	var view viewModels.Start
	view.Cases = repository.Case_Get_All()

	err := Template_View(w, view, "selection/start", "templates/selection/start.gohtml")
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func Selection_View(w http.ResponseWriter, r *http.Request) {
	var view viewModels.SelectTechnology
	view.Categories = repository.Category_Get_All()
	view.Needs.Prefix = "needs"

	needs := repository.Need_Get_All()
	view.Needs.Options = make([]models.CheckableOption, len(needs))
	for i, need := range needs {
		view.Needs.Options[i] = models.CheckableOption{
			Selected:    false,
			Description: need.Description,
			Name:        strconv.Itoa(need.Id),
		}
	}
	if r.Method == "POST" {
		r.ParseForm()
		err := decoder.Decode(&view, r.Form)
		if err != nil {
			fmt.Println(err)
		}

		clientCase := repository.Case_Get_By_Id(view.Case.Id)
		if clientCase != nil {
			view.Case = *clientCase
		}
	}
	for i, category := range view.Categories {
		technologies := repository.Tech_Get_All_By_Category(category.Id)
		results := make([]models.Tech, 0)
		for _, tech := range technologies {
			if tech_has_all_needs(&tech, view.Needs.Options) {
				results = append(results, tech)
			}
		}
		view.Categories[i].Technologies = results
	}

	categoryResults := make([]models.Category, 0)
	for _, category := range view.Categories {
		if len(category.Technologies) > 0 {
			categoryResults = append(categoryResults, category)
		}
	}
	view.Categories = categoryResults

	err := Template_View(w, view, "selection/select", "templates/selection/selectTechnology.gohtml")
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func Shortlist_View(w http.ResponseWriter, r *http.Request) {
	var view viewModels.Shortlist
	if r.Method == "POST" {
		r.ParseForm()
		err := decoder.Decode(&view, r.Form)
		if err != nil {
			fmt.Println(err)
		}
		view.Case = repository.Case_Get_By_Id(view.Case.Id)
		view.Choices = repository.TechChoice_Get_All_By_Case(view.Case.Id)
	}

	err := Template_View(w, view, "selection/shortlist", "templates/selection/shortlist.gohtml")
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}
