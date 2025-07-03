package auth

import (
	"authentication/pkg/models"
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// generateTokens is an internal helper for authService
func (s *authService) generateTokens(ctx context.Context, user models.User) (string, string, error) {
	// ... (your existing generateTokens implementation) ...
	// Ensure calls like s.userRepo.LoadUserRoles and s.roleRepo.LoadRolePermissions use `ctx`
	// Example:
	err := s.userRepo.LoadUserRoles(ctx, &user)
	if err != nil {
		return "", "", fmt.Errorf("failed to load user roles for token generation: %w", err)
	}
	for i := range user.Roles {
		err := s.roleRepo.LoadRolePermissions(ctx, &user.Roles[i])
		if err != nil {
			log.Printf("Warning: Could not load permissions for role %s: %v\n", user.Roles[i].Name, err)
		}
	}

	var permissions []string
	for _, role := range user.Roles {
		for _, perm := range role.Permissions {
			permissions = append(permissions, fmt.Sprintf("%s:%s", perm.Resource, perm.Action))
		}
	}

	accessClaims := JWTClaims{
		UserID:      user.ID,
		Username:    user.Username,
		Roles:       s.getRoleNames(user.Roles),
		Permissions: permissions,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(15 * time.Minute).Unix(),
			Issuer:    "auth-service",
		},
	}

	// Assuming jwt.NewWithClaims and accessToken.SignedString are correctly handled by a JWT utility
	// or directly imported from a JWT library.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := token.SignedString(s.jwtSecret)
	if err != nil {
		return "", "", fmt.Errorf("failed to sign access token: %w", err)
	}

	refreshTokenString := uuid.New().String()
	refreshExpiry := time.Now().Add(24 * time.Hour * 7)

	authToken := models.AuthToken{
		ID:        uuid.New(),
		Token:     refreshTokenString,
		UserID:    user.ID,
		ExpiresAt: refreshExpiry,
	}

	if err := s.authTokenRepo.Create(ctx, &authToken); err != nil { // Use ctx here
		return "", "", fmt.Errorf("failed to save refresh token: %w", err)
	}

	return accessTokenString, refreshTokenString, nil
}

// RegisterAPIEndpoint implements AuthService.RegisterAPIEndpoint.
// Note: It's crucial that this method now belongs to the 'authService' struct.
// It also needs to accept a context.Context.
func (s *authService) RegisterAPIEndpoint(ctx context.Context, cmd models.RegisterAPIEndpointCommand) (*models.APIEndpoint, error) { // Added ctx, changed return to pointer
	// Check if endpoint already exists using the injected endpointRepo
	existing, err := s.endpointRepo.FindByPathAndMethod(ctx, cmd.Path, cmd.Method) // Pass ctx
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {                     // Check for actual DB errors, not just "not found"
		return nil, fmt.Errorf("database error checking for existing endpoint: %w", err)
	}
	if existing != nil {
		return nil, errors.New("endpoint already exists") // Changed return to nil, errors.New for consistency
	}

	// Find roles
	var roles []models.Role // Use models.Role
	for _, roleName := range cmd.RoleNames {
		role, err := s.roleRepo.FindByName(ctx, roleName) // Pass ctx
		if err != nil {
			return nil, fmt.Errorf("database error finding role %s: %w", roleName, err)
		}
		if role == nil { // Check if role was actually found
			return nil, fmt.Errorf("role not found: %s", roleName)
		}
		roles = append(roles, *role) // Append the dereferenced role
	}

	// Create endpoint
	endpoint := models.APIEndpoint{ // Use models.APIEndpoint
		ID:          uuid.New(),
		ServiceName: cmd.ServiceName,
		Path:        cmd.Path,
		Method:      cmd.Method,
		Description: cmd.Description,
		Version:     cmd.Version,
		Roles:       roles, // Assign the fetched roles
		CreatedAt:   time.Now(),
	}

	// Create the endpoint using the injected endpointRepo
	// The Create method should accept a pointer if it's modifying the object or needs its ID back.
	if err := s.endpointRepo.Create(ctx, &endpoint); err != nil { // Pass ctx and pointer
		return nil, fmt.Errorf("failed to create API endpoint: %w", err)
	}

	// Publish event
	// Ensure EndpointRegisteredEvent is defined in models or an events package
	// Also ensure s.eventPublisher is properly initialized and its Publish method is correctly defined.
	if s.eventPublisher != nil { // Check if eventPublisher is set
		event := models.EndpointRegisteredEvent{ // Assuming this struct is in models
			EndpointID:   endpoint.ID,
			Path:         endpoint.Path,
			Method:       endpoint.Method,
			Service:      endpoint.ServiceName,
			RegisteredBy: cmd.RegisteredBy,
			Timestamp:    time.Now(),
		}
		if err := s.eventPublisher.Publish(event); err != nil {
			log.Printf("Warning: Failed to publish EndpointRegisteredEvent for %s: %v", endpoint.Path, err)
			// Decide if this is a critical error or just a warning.
		}
	}

	return &endpoint, nil // Return pointer to the created endpoint
}

// ... (Other AuthService methods will go here, using s.userRepo, s.roleRepo, etc.)

// Placeholder for gorm.ErrRecordNotFound if not imported directly
var gormErrRecordNotFound = errors.New("record not found") // You should import "gorm.io/gorm" and use gorm.ErrRecordNotFound
