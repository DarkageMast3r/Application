package models

import (
	"time"

	"github.com/google/uuid"
)

// UserRegisteredEvent triggered after successful registration
type UserRegisteredEvent struct {
	UserID    uuid.UUID `json:"userId"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Timestamp time.Time `json:"timestamp"`
}

// UserLoggedInEvent for tracking logins
type UserLoggedInEvent struct {
	UserID    uuid.UUID `json:"userId"`
	IPAddress string    `json:"ipAddress"`
	UserAgent string    `json:"userAgent"`
	Timestamp time.Time `json:"timestamp"`
}

// RoleAssignedEvent when a user gets a new role
type RoleAssignedEvent struct {
	UserID     uuid.UUID `json:"userId"`
	RoleID     uuid.UUID `json:"roleId"`
	AssignedBy uuid.UUID `json:"assignedBy"`
	Timestamp  time.Time `json:"timestamp"`
}

// EndpointRegisteredEvent when a new API endpoint is added
type EndpointRegisteredEvent struct {
	EndpointID   uuid.UUID `json:"endpointId"`
	Path         string    `json:"path"`
	Method       string    `json:"method"`
	Service      string    `json:"service"`
	RegisteredBy uuid.UUID `json:"registeredBy"`
	Timestamp    time.Time `json:"timestamp"`
}
