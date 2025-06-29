package models

import (
	"github.com/google/uuid"
)

type TechChoice struct {
	Id       int `json:"id" excludeFromCreate:"true"`
	TechId   int `json:"tech_id"`
	Tech     *Tech
	ClientId uuid.UUID `json:"client_id"`
	Status   int       `json:"status"`
}
