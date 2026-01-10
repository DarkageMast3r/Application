package models

import (
	"time"

	"github.com/google/uuid"
)

type AuditAction string

const (
	ActionLogin         AuditAction = "LOGIN"
	ActionLogout        AuditAction = "LOGOUT"
	ActionPasswordReset AuditAction = "PASSWORD_RESET"
	ActionRoleChanged   AuditAction = "ROLE_CHANGED"
)

type AuditLog struct {
	ID        uuid.UUID   `gorm:"primaryKey;type:uuid"`
	Action    AuditAction `gorm:"index"`
	UserID    *uuid.UUID  `gorm:"type:uuid;index"` // Null voor system events
	IPAddress string
	UserAgent string
	Metadata  map[string]interface{} `gorm:"serializer:json"`
	CreatedAt time.Time
}
