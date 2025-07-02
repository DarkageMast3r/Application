package models

type ToestandClassificatie struct {
	Categorie string `json:"categorie"`
	Ernst     string `json:"ernst"`
	Motivatie string `json:"motivatie"`
}
