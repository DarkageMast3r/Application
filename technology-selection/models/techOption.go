package models

import "github.com/google/uuid"

type TechOption struct {
	Id         int `json:"id" excludeFromCreate:"false"`
	ClientId   uuid.UUID
	Category   Category
	CategoryId int
}
