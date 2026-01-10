package models

import "github.com/google/uuid"

// UserQuery for retrieving user information
type UserQuery struct {
	UserID   *uuid.UUID `json:"userId"`
	Username *string    `json:"username"`
	Email    *string    `json:"email"`
}

// RoleQuery for retrieving role information
type RoleQuery struct {
	RoleID *uuid.UUID `json:"roleId"`
	Name   *string    `json:"name"`
}

// PermissionQuery for retrieving permissions
type PermissionQuery struct {
	Resource *string `json:"resource"`
	Action   *string `json:"action"`
}

// EndpointQuery for querying registered endpoints
type EndpointQuery struct {
	ServiceName *string `json:"serviceName"`
	Path        *string `json:"path"`
	Method      *string `json:"method"`
	Version     *string `json:"version"`
}

// AuthorizationQuery for checking permissions
type AuthorizationQuery struct {
	UserID   uuid.UUID `json:"userId" validate:"required"`
	Resource string    `json:"resource" validate:"required"`
	Action   string    `json:"action" validate:"required"`
}
