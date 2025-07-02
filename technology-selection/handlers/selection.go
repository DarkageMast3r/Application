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

var cases []models.Case = []models.Case{
	{
		Id:          1,
		ClientId:    "AAAA-AAAA",
		Name:        "John Smith",
		Description: "Complete loss of mobility after getting involved in a car incident.",
	},
	{
		Id:          2,
		ClientId:    "AAAA-AAAB",
		Name:        "Raphael Stones",
		Description: "Paralised from the neck down after football incident",
	},
	{
		Id:          3,
		ClientId:    "AAAA-AAAC",
		Name:        "Jordan Sanderson",
		Description: "Cannot sit comfortably anymore",
	},
	{
		Id:          4,
		ClientId:    "AAAA-AAAD",
		Name:        "Sarah Flowerfield",
		Description: "Lost his left leg in a mountaineering incident, but wishes to continue",
	},
}

func case_get_by_id(id int) *models.Case {
	for _, clientCase := range cases {
		if clientCase.Id == id {
			return &clientCase
		}
	}
	return nil
}

func Start_View(w http.ResponseWriter, r *http.Request) {
	var view viewModels.Start
	view.Cases = cases

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

		clientCase := case_get_by_id(view.Case.Id)
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
		view.Case = case_get_by_id(view.Case.Id)
		view.Choices = repository.TechChoice_Get_All_By_Case(view.Case.Id)
	}

	err := Template_View(w, view, "selection/shortlist", "templates/selection/shortlist.gohtml")
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}
