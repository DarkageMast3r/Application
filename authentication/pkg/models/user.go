package models

import (
	"errors"
	"time"
	"unicode"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// User represents a user entity
type User struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey"`
	Username     string    `gorm:"uniqueIndex;size:50"`
	Email        string    `gorm:"uniqueIndex;size:100"`
	PasswordHash string    `gorm:"size:255"`
	IsActive     bool      `gorm:"default:true"`
	LastLogin    *time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Roles        []Role `gorm:"many2many:user_roles;"` // Many-to-many relatie
}

// Role represents a role entity
type Role struct {
	ID          uuid.UUID    `gorm:"type:uuid;primaryKey"`
	Name        string       `gorm:"uniqueIndex;size:50"`
	Description string       `gorm:"size:255"`
	Permissions []Permission `gorm:"many2many:role_permissions;"` // Many-to-many relatie
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// Permission represents a permission entity
type Permission struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	Resource  string    `gorm:"size:100"` // e.g., "users", "roles", "endpoints"
	Action    string    `gorm:"size:50"`  // e.g., "read", "create", "update", "delete"
	CreatedAt time.Time
	UpdatedAt time.Time
}

// APIEndpoint represents an API endpoint and its access roles
type APIEndpoint struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey"`
	ServiceName string    `gorm:"size:100"`
	Path        string    `gorm:"size:255"`
	Method      string    `gorm:"size:10"` // GET, POST, PUT, DELETE, PATCH
	Description string    `gorm:"size:255"`
	Version     string    `gorm:"size:20"`
	Roles       []Role    `gorm:"many2many:endpoint_roles;"` // Many-to-many relatie
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// RefreshToken stores valid refresh tokens for users
type RefreshToken struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserID    uuid.UUID `gorm:"type:uuid;index"`
	Token     string    `gorm:"uniqueIndex;size:255"`
	ExpiresAt time.Time
	CreatedAt time.Time
}

// PasswordResetToken stores tokens for password reset requests
type PasswordResetToken struct {
	Token     string    `gorm:"primaryKey;size:255"`
	UserID    uuid.UUID `gorm:"type:uuid;index"`
	ExpiresAt time.Time
	Used      bool `gorm:"default:false"`
	CreatedAt time.Time
}

// TokenBlacklist stores invalidated JWT tokens
type TokenBlacklist struct {
	Token     string    `gorm:"primaryKey;size:255"`
	ExpiresAt time.Time `gorm:"index"`
}

// Password is a value object for password operations
type Password struct {
	value string
}

// NewPassword creates a new Password value object after validation
func NewPassword(plaintext string) (Password, error) {
	if len(plaintext) < 8 {
		return Password{}, errors.New("password must be at least 8 characters long")
	}

	var (
		hasUpper, hasLower, hasNumber, hasSpecial bool
	)

	for _, c := range plaintext {
		switch {
		case unicode.IsUpper(c):
			hasUpper = true
		case unicode.IsLower(c):
			hasLower = true
		case unicode.IsNumber(c):
			hasNumber = true
		case unicode.IsPunct(c) || unicode.IsSymbol(c):
			hasSpecial = true
		}
	}

	if !hasUpper || !hasLower || !hasNumber || !hasSpecial {
		return Password{}, errors.New("password must contain uppercase, lowercase, number, and special character")
	}

	return Password{value: plaintext}, nil
}

// Hash generates a bcrypt hash of the password
func (p Password) Hash() (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(p.value), bcrypt.DefaultCost)
	if err != nil {
		return "", errors.New("failed to hash password")
	}
	return string(hashedBytes), nil
}

// MatchesHash compares the plaintext password with a hashed password
func (p Password) MatchesHash(hash string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(p.value))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return false, nil // Passwords do not match
		}
		return false, errors.New("error comparing passwords") // Other error
	}
	return true, nil // Passwords match
}

// JWTClaims custom claims structure
type JWTClaims struct {
	UserID      uuid.UUID `json:"userId"`
	Username    string    `json:"username"`
	Roles       []string  `json:"roles"`
	Permissions []string  `json:"permissions"` // Flattened permissions from all roles
	jwt.StandardClaims
}

// UserRole links users to roles (join table with additional fields)
type UserRole struct {
	UserID     uuid.UUID `json:"userId" gorm:"primaryKey;type:uuid"`
	RoleID     uuid.UUID `json:"roleId" gorm:"primaryKey;type:uuid"`
	AssignedBy uuid.UUID `json:"assignedBy" gorm:"type:uuid"`
	AssignedAt time.Time `json:"assignedAt" gorm:"autoCreateTime"`
}
