package models

import "time"

type Category struct {
	Id            int       `json:"id"`
	Naam          string    `json:"name"`
	Beschrijving  string    `json:"description"`
	GegenereerdOp time.Time `json:"generated_on"`
}
