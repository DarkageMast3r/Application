package models

import (
	"time"

	"github.com/google/uuid"
)

type TechChoice struct {
	Id         int `json:"id" excludeFromCreate:"false"`
	ClientId   uuid.UUID
	Category   Category
	CategoryId int
	Options    []TechOption
	Status     int
	LastUpdate time.Time
}
