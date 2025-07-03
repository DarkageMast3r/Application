package models

import (
	"time"

	"github.com/google/uuid"
)

// TechnischDetail representeert één technisch detail item
type TechnischDetail struct {
	Sleutel string `json:"sleutel"`
	Waarde  string `json:"waarde"`
}

// ZorgTechProduct is de root entity voor zorgtechnologieproducten
type ZorgTechProduct struct {
	ID                uuid.UUID         `json:"zorgtechId" gorm:"primary_key;type:uuid"`
	Naam              string            `json:"naam"`
	Beschrijving      string            `json:"beschrijving"`
	Categorie         string            `json:"categorie"` // Vrije tekst ipv enum
	TechnischeDetails []TechnischDetail `json:"technischeDetails" gorm:"type:jsonb"`
	Prijs             float64           `json:"prijs"`
	Leverancier       string            `json:"leverancier"`
	IsActief          bool              `json:"isActief" gorm:"default:true"`
	CreatedAt         time.Time         `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt         time.Time         `json:"updatedAt" gorm:"autoUpdateTime"`
	Events            []ProductEvent    `json:"-" gorm:"foreignKey:ProductID"`
}

// ProductEvent is een base struct voor domain events
type ProductEvent struct {
	EventID     uuid.UUID `json:"eventId" gorm:"primary_key;type:uuid"`
	ProductID   uuid.UUID `json:"productId" gorm:"type:uuid;not null"`
	Type        string    `json:"type"`
	Payload     string    `json:"payload"`
	OccurredAt  time.Time `json:"occurredAt" gorm:"autoCreateTime"`
	TriggeredBy string    `json:"triggeredBy"`
}

// Domain Events
type ZorgTechProductAangemaakt struct {
	ProductID uuid.UUID `json:"productId"`
	Naam      string    `json:"naam"`
}

type ZorgTechProductGewijzigd struct {
	ProductID   uuid.UUID              `json:"productId"`
	Wijzigingen map[string]interface{} `json:"wijzigingen"`
}

type ZorgTechProductVerwijderd struct {
	ProductID uuid.UUID `json:"productId"`
}

// Command models
type MaakZorgTechProductCommand struct {
	Naam              string            `json:"naam" binding:"required"`
	Beschrijving      string            `json:"beschrijving" binding:"required"`
	Categorie         string            `json:"categorie" binding:"required"`
	TechnischeDetails []TechnischDetail `json:"technischeDetails"`
	Prijs             float64           `json:"prijs" binding:"required,min=0"`
	Leverancier       string            `json:"leverancier" binding:"required"`
}

type WijzigZorgTechProductCommand struct {
	ProductID         uuid.UUID          `json:"productId" binding:"required"`
	Naam              *string            `json:"naam,omitempty"`
	Beschrijving      *string            `json:"beschrijving,omitempty"`
	Categorie         *string            `json:"categorie,omitempty"`
	TechnischeDetails *[]TechnischDetail `json:"technischeDetails,omitempty"`
	Prijs             *float64           `json:"prijs,omitempty"`
	Leverancier       *string            `json:"leverancier,omitempty"`
}

type VerwijderZorgTechProductCommand struct {
	ProductID uuid.UUID `json:"productId" binding:"required"`
}

type VoegTechnischDetailToeCommand struct {
	ProductID uuid.UUID `json:"productId" binding:"required"`
	Sleutel   string    `json:"sleutel" binding:"required"`
	Waarde    string    `json:"waarde" binding:"required"`
}

type VerwijderTechnischDetailCommand struct {
	ProductID uuid.UUID `json:"productId" binding:"required"`
	Sleutel   string    `json:"sleutel" binding:"required"`
}

// Query models
type GetProductByIdQuery struct {
	ZorgtechID uuid.UUID `json:"zorgtechId" binding:"required"`
}

type FindByCategorieQuery struct {
	Categorie string `json:"categorie" binding:"required"`
}

type ZoekOpNaamQuery struct {
	Zoekterm string `json:"zoekterm" binding:"required"`
}
