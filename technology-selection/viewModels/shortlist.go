package viewModels

import "service/models"

type Shortlist struct {
	Choices []models.TechChoice
	Case    *models.Case `json:"Case"`
}
