package viewModels

import "service/models"

type Start struct {
	Cases  []models.Case
	CaseId string `schema:"case_id"`
}
