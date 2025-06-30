package models

type Need struct {
	Id          int    `schema:"id" json:"id"`
	Description string `schema:"description" json:"description"`
}
