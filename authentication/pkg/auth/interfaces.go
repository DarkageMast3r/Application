package auth

import (
	"authentication/pkg/models"

	"github.com/google/uuid"
)

type RoleService interface {
	CreateRole(cmd CreateRoleCommand) (models.Role, error)
	AssignRole(cmd AssignRoleCommand) error
	RemoveRole(userID uuid.UUID, roleID uuid.UUID) error
	GetUserRoles(userID uuid.UUID) ([]models.Role, error)
}

type APIManagementService interface {
	RegisterEndpoint(cmd RegisterAPIEndpointCommand) (models.APIEndpoint, error)
	UpdateEndpointPermissions(endpointID uuid.UUID, roleNames []string) error
	GetAllEndpoints() ([]models.APIEndpoint, error)
	GetEndpointByID(id uuid.UUID) (models.APIEndpoint, error)
}
