package models

type Category struct {
	Id          int    `schema:"id" json:"id"`
	Name        string `schema:"name" json:"name"`
	Description string `schema:"description" json:"description"`
}
