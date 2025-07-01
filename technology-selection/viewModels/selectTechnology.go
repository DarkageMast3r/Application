package viewModels

import "service/models"

type SelectTechnology struct {
	Categories []models.Category
	Needs      models.Checklist `schema:"needs"`
	Case       models.Case      `schema:"case"`
}
