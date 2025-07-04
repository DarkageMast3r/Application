package models

import (
	"time"
	"github.com/google/uuid"
)

// ZorgTechProduct represents a healthcare technology product
type ZorgTechProduct struct {
	ZorgtechID       uuid.UUID `json:"zorgtechId" gorm:"primary_key;type:uuid"`
	Name             string    `json:"naam" gorm:"not null"`
	Description      string    `json:"beschrijving"`
	Category         string    `json:"categorie"`
	TechnicalDetails string    `json:"technischeDetails"`
	Price           float64   `json:"prijs"`
	Supplier        string    `json:"leverancier"`
	CreatedAt       time.Time `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt       time.Time `json:"updatedAt" gorm:"autoUpdateTime"`
}

// ImplementatieStatus represents the possible statuses of an implementation
type ImplementatieStatus string

const (
	StatusBesteld       ImplementatieStatus = "Besteld"
	StatusGeleverd      ImplementatieStatus = "Geleverd"
	StatusGeinstalleerd ImplementatieStatus = "Geïnstalleerd"
	StatusGepersonaliseerd ImplementatieStatus = "Gepersonaliseerd"
	StatusVoltooid      ImplementatieStatus = "Voltooid"
)

// PersonalisatieInstellingen is a value object for personalization settings
type PersonalisatieInstellingen struct {
	VolumeNiveau  int    `json:"volumeNiveau"`  // 1-10
	MeldingsType  string `json:"meldingsType"`  // Auditief, visueel, trilling
	Schema        string `json:"schema"`        // dag/nachtpatroon
	Taal          string `json:"taal"`          // Instellingstaal
}

// TrackingInfo is a value object for tracking information
type TrackingInfo struct {
	LeverancierNaam       string    `json:"leverancierNaam"`
	TrackingCode          string    `json:"trackingCode"`
	Status                string    `json:"status"`                // Onderweg, Afgeleverd, etc.
	VerwachteLeverdatum   time.Time `json:"verwachteLeverdatum"`
}


type ImplementatieLog struct {
	ID                    uuid.UUID `gorm:"primary_key;type:uuid"`
	ImplementatieDossierID uuid.UUID `gorm:"type:uuid;not null"` // foreign key
	Timestamp    time.Time `json:"timestamp" gorm:"autoCreateTime"`
	Actie        string    `json:"actie"`
	UitgevoerdDoor string    `json:"uitgevoerdDoor"`
}

// ImplementatieDossier is the root entity for the implementation process
type ImplementatieDossier struct {
	ImplementatieID  uuid.UUID               `json:"implementatieId" gorm:"primary_key;type:uuid"`
	ClientID        uuid.UUID               `json:"clientId" gorm:"type:uuid;not null"`
	ZorgtechID      uuid.UUID               `json:"zorgtechId" gorm:"type:uuid;not null"`
	Status          ImplementatieStatus     `json:"status"`
	Serienummer     string                  `json:"serienummer"`
	InstallatieDatum *time.Time              `json:"installatieDatum"`
	Personalisatie  *PersonalisatieInstellingen `json:"personalisatie" gorm:"embedded;embeddedPrefix:personalisatie_"`
	TrackingInfo    *TrackingInfo            `json:"trackingInfo" gorm:"embedded;embeddedPrefix:tracking_"`
	Logs            []ImplementatieLog       `json:"logs" gorm:"foreignKey:ImplementatieDossierID"`
	CreatedAt       time.Time                `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt       time.Time                `json:"updatedAt" gorm:"autoUpdateTime"`
}

// Command and Request models

type AanvraagProductCommand struct {
	ClientID   uuid.UUID `json:"clientId" binding:"required"`
	ZorgtechID uuid.UUID `json:"zorgtechId" binding:"required"`
}

type OntvangProductCommand struct {
	ClientID   uuid.UUID `json:"clientId" binding:"required"`
	ZorgtechID uuid.UUID `json:"zorgtechId" binding:"required"`
	Serienummer string    `json:"serienummer" binding:"required"`
}

type InstalleerProductCommand struct {
	ClientID    uuid.UUID `json:"clientId" binding:"required"`
	Serienummer string    `json:"serienummer" binding:"required"`
}

type PersonaliseerProductCommand struct {
	ClientID    uuid.UUID                  `json:"clientId" binding:"required"`
	Instellingen PersonalisatieInstellingen `json:"instellingen" binding:"required"`
}

type LeverProductCommand struct {
	ClientID   uuid.UUID `json:"clientId" binding:"required"`
	ZorgtechID uuid.UUID `json:"zorgtechId" binding:"required"`
}

type MarkeerAlsGeimplementeerdCommand struct {
	ClientID uuid.UUID `json:"clientId" binding:"required"`
}

// Query models

type ProductInformatieQuery struct {
	ZorgtechID uuid.UUID `json:"zorgtechId" binding:"required"`
}

type ImplementatieStatusQuery struct {
	ClientID uuid.UUID `json:"clientId" binding:"required"`
}

type InstallatieStatusQuery struct {
	ClientID uuid.UUID `json:"clientId" binding:"required"`
}

type PersoonlijkeInstellingenQuery struct {
	ClientID uuid.UUID `json:"clientId" binding:"required"`
}