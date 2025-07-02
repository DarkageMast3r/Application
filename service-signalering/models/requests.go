package models

import "github.com/google/uuid"

type RegistreerAchteruitgangRequest struct {
	Signalen []Signaal `json:"signalen"`
}

type ClassificeerSituatieRequest struct {
	Classificatie ToestandClassificatie `json:"classificatie"`
}

type BeoordeelSituatieRequest struct {
	Conclusie       string `json:"conclusie"`
	Urgentie        string `json:"urgentie"`
	GevalideerdDoor string `json:"gevalideerd_door"`
}

type RegelDoorverwijzingRequest struct {
	ProbleemID uuid.UUID `json:"probleem_id"`
	Type       string    `json:"type"`
	Urgentie   string    `json:"urgentie"`
	Notities   string    `json:"notities"`
}
