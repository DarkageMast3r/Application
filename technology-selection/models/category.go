package models

import (
	"time"
)

type Category struct {
	Id          int       `json:"id" excludeFromCreate:"true"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	GeneratedOn time.Time `json:"generated_on"  excludeFromCreate:"true"`
}
