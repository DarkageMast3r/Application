package models

type Need struct {
	Id     int    `json:"id" excludeFromCreate:"false"`
	Name   string `json:"name"`
	Source string `json:"source"`
}
