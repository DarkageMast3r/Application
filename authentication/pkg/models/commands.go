package models

import "github.com/google/uuid"

// RegisterUserCommand for user registration
type RegisterUserCommand struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

// LoginCommand for user authentication
type LoginCommand struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// CreateRoleCommand for role management
type CreateRoleCommand struct {
	Name        string   `json:"name" validate:"required"`
	Description string   `json:"description"`
	Permissions []string `json:"permissions"` // Permission IDs
}

// AssignRoleCommand for assigning roles to users
type AssignRoleCommand struct {
	UserID  uuid.UUID `json:"userId" validate:"required"`
	RoleID  uuid.UUID `json:"roleId" validate:"required"`
	AdminID uuid.UUID `json:"adminId" validate:"required"` // Who is assigning this
}

// RegisterAPIEndpointCommand for endpoint registration
type RegisterAPIEndpointCommand struct {
	ServiceName string   `json:"serviceName" validate:"required"`
	Path        string   `json:"path" validate:"required"`
	Method      string   `json:"method" validate:"required,oneof=GET POST PUT PATCH DELETE"`
	Description string   `json:"description"`
	Version     string   `json:"version" validate:"required"`
	RoleNames   []string `json:"roleNames"` // Which roles can access this
}

type LoginUser struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}
