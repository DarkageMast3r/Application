package viewModels

import "service/models"

type TechCreate struct {
	Categories []models.Category
	Needs      []models.Need
}
