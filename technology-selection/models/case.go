package models

type Case struct {
	Id          int
	Name        string `json:"name"`
	ClientId    string `json:"client_id"`
	Description string `json:"description"`
}
