package models

type Need struct {
	Id          int    `json:"id" excludeFromCreate:"true"`
	TechId      int    `json:"tech_id"`
	Description string `json:"string"`
}
