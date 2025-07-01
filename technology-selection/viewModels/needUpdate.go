package viewModels

import "service/models"

type NeedUpdate struct {
	Need         *models.Need
	Technologies []models.Tech
}
