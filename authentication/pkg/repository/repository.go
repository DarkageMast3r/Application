package repository

import (
	"authentication/pkg/models"
	"context"
	"time"

	"github.com/google/uuid"
)

// UserRepository defines the interface for user data operations
type UserRepository interface {
	FindByUsernameOrEmail(ctx context.Context, username, email string) (*models.User, error)
	FindByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	Save(ctx context.Context, user *models.User) error
	UpdateLastLogin(ctx context.Context, user *models.User) error
	LoadUserRoles(ctx context.Context, user *models.User) error       // Laad de rollen voor een gebruiker
	LoadRolePermissions(ctx context.Context, role *models.Role) error // Laad de permissies voor een rol
}

// RoleRepository defines the interface for role data operations
type RoleRepository interface {
	Create(ctx context.Context, role *models.Role) error
	FindByID(ctx context.Context, id uuid.UUID) (*models.Role, error)
	FindByName(ctx context.Context, name string) (*models.Role, error)
	FindAll(ctx context.Context) ([]models.Role, error)
	Save(ctx context.Context, role *models.Role) error // Voor updates
	Delete(ctx context.Context, id uuid.UUID) error
	LoadRolePermissions(ctx context.Context, role *models.Role) error
	// Optioneel: FindByNames(ctx context.Context, names []string) ([]models.Role, error)
}

// EndpointRepository defines the interface for API endpoint data operations
type EndpointRepository interface {
	Create(ctx context.Context, endpoint *models.APIEndpoint) error
	FindByID(ctx context.Context, id uuid.UUID) (*models.APIEndpoint, error)
	FindByPathAndMethod(ctx context.Context, path, method string) (*models.APIEndpoint, error)
	FindAll(ctx context.Context) ([]models.APIEndpoint, error)
	Save(ctx context.Context, endpoint *models.APIEndpoint) error // Voor updates
	Delete(ctx context.Context, id uuid.UUID) error
	LoadEndpointRoles(ctx context.Context, endpoint *models.APIEndpoint) error
}

// AuthTokenRepository defines the interface for authentication token operations
type AuthTokenRepository interface {
	Create(ctx context.Context, token *models.RefreshToken) error
	DeleteByID(ctx context.Context, id uuid.UUID) error
	DeleteExpiredTokens(ctx context.Context) error
	DeleteUserTokens(ctx context.Context, userID uuid.UUID) error // Voor logout all devices
	Save(ctx context.Context, token *models.RefreshToken) error   // Voor updates (e.g., used_at)

	// Blacklist Access Tokens
	AddBlacklistedAccessToken(ctx context.Context, token string, expiresAt time.Time) error
	IsAccessTokenBlacklisted(ctx context.Context, token string) (bool, error)
	CleanExpiredBlacklistedTokens(ctx context.Context) error
	SaveRefreshToken(ctx context.Context, token *models.RefreshToken) error
	FindByRefreshToken(ctx context.Context, token string) (*models.RefreshToken, error)
	DeleteRefreshToken(ctx context.Context, id uuid.UUID) error
	DeleteRefreshTokensByUserID(ctx context.Context, userID uuid.UUID) error
	DeleteExpiredRefreshTokens(ctx context.Context) error

	// Password Reset Token operations
	SavePasswordResetToken(ctx context.Context, token *models.PasswordResetToken) error
	FindValidPasswordResetToken(ctx context.Context, token string) (*models.PasswordResetToken, error)
	MarkPasswordResetTokenUsed(ctx context.Context, token *models.PasswordResetToken) error
	DeleteExpiredPasswordResetTokens(ctx context.Context) error

	// Blacklisted Token operations
	AddBlacklistedToken(ctx context.Context, token *BlacklistedToken) error
	IsTokenBlacklisted(ctx context.Context, token string) (bool, error)
	DeleteExpiredBlacklistedTokens(ctx context.Context) error
}

// CacheRepository defines the interface for caching operations
type CacheRepository interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Del(ctx context.Context, key string) error
	Exists(ctx context.Context, key string) (bool, error)                   // Handig voor blacklist checks
	Increment(ctx context.Context, key string) (int64, error)               // Voor rate limiting
	Expire(ctx context.Context, key string, expiration time.Duration) error // Voor het instellen van een vervaldatum na increment
	// Specifiek voor blacklisted access tokens
	AddBlacklistedAccessToken(ctx context.Context, token string, expiresAt time.Time) error
	IsAccessTokenBlacklisted(ctx context.Context, token string) (bool, error)
	// CleanExpiredBlacklistedTokens hoeft niet apart, want Redis doet dit met TTL
}

// EventPublisher defines the interface for publishing domain events
type EventPublisher interface {
	Publish(ctx context.Context, event interface{}) error
}
