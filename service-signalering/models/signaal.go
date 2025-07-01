package models

import "time"

type Signaal struct {
	Type     string    `json:"type"`
	Waarde   float64   `json:"waarde"`
	Tijdstip time.Time `json:"tijdstip"`
	Bron     string    `json:"bron"`
}
