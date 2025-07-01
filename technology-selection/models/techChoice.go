package models

import "database/sql"

type TechChoice struct {
	Id        int `json:"id" excludeFromCreate:"true"`
	TechId    int `json:"tech_id"`
	Tech      *Tech
	Reasoning sql.NullString `schema:"reasoning"`
	CaseId    int            `json:"case_id"`
	Status    int            `json:"status"`
}
