package models

type Behoefte struct {
	Id           int    `json:"id"`
	Beschrijving string `json:"description"`
	Bron         string `json:"source"`
}
