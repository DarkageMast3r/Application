package models

import (
	"time"

	"github.com/google/uuid"
)

type ClientToestand struct {
	ToestandID          uuid.UUID              `json:"toestand_id"`
	ClientID            uuid.UUID              `json:"client_id"`
	Signalen            []Signaal              `json:"signalen"`
	Classificatie       *ToestandClassificatie `json:"classificatie"`
	Beoordeling         *Beoordeling           `json:"beoordeling"`
	Status              string                 `json:"status"`
	TijdstipRegistratie time.Time              `json:"tijdstip_registratie"`
}
