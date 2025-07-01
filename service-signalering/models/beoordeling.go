package models

import "time"

type Beoordeling struct {
	Conclusie       string    `json:"conclusie"`
	Urgentie        string    `json:"urgentie"`
	GevalideerdDoor string    `json:"gevalideerd_door"`
	Tijdstip        time.Time `json:"tijdstip"`
}
