package models

import (
	"time"

	"github.com/google/uuid"
)

type ToestandResponse struct {
	ToestandID          uuid.UUID              `json:"toestand_id"`
	ClientID            uuid.UUID              `json:"client_id"`
	Status              string                 `json:"status"`
	TijdstipRegistratie time.Time              `json:"tijdstip_registratie"`
	Signalen            []Signaal              `json:"signalen,omitempty"`
	Classificatie       *ToestandClassificatie `json:"classificatie,omitempty"`
	Beoordeling         *Beoordeling           `json:"beoordeling,omitempty"`
}

type BeoordelingResponse struct {
	ClientID         uuid.UUID    `json:"client_id"`
	LaatsBeoordeling *Beoordeling `json:"laatste_beoordeling"`
	Aanbevelingen    []string     `json:"aanbevelingen"`
}

type SimpleErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// Gedetailleerde errors voor validatie
type ErrorResponse struct {
	Code    string   `json:"code"`
	Message string   `json:"message"`
	Details []string `json:"details,omitempty"`
}
