package viewModels

import "service/models"

type TechUpdate struct {
	Tech       *models.Tech
	Categories []models.Category
	Needs      []models.Need
}
