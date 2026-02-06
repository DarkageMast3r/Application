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

// AuditLog for tracking security-relevant events
type AuditLog struct {
	ID          uuid.UUID  `json:"id" gorm:"primary_key;type:uuid"`
	UserID      *uuid.UUID `json:"userId" gorm:"type:uuid"` // Null for system events
	Action      string     `json:"action" gorm:"not null"`  // e.g., "LOGIN", "PERMISSION_CHANGE"
	Description string     `json:"description"`
	IPAddress   string     `json:"ipAddress"`
	UserAgent   string     `json:"userAgent"`
	CreatedAt   time.Time  `json:"createdAt" gorm:"autoCreateTime"`
}
