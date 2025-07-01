package models

type Tech struct {
	Id         int `schema:"id" json:"id"`
	Category   *Category
	CategoryId int     `schema:"category_id" json:"category_id"`
	Needs      []Need  `schema:"needs"`
	Name       string  `schema:"name" json:"name"`
	Cost       float64 `schema:"cost" json:"cost"`
}
