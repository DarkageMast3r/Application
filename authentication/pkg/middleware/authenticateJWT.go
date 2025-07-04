package middleware

import (
	"authentication/pkg/models"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

/*
var jwtSecret []byte

func InitJWT(secret string) {
	jwtSecret = []byte(secret)
}
*/

// JWTWrapper handles JWT operations with  encapsulation
type JWTWrapper struct {
	SecretKey []byte
}

// NewJWTWrapper creates a new JWT wrapper instance
func NewJWTWrapper(secret string) *JWTWrapper {
	return &JWTWrapper{
		SecretKey: []byte(secret),
	}
}

// GenerateJWT generates a new JWT token
func (j *JWTWrapper) GenerateJWT(userID uuid.UUID, username string, roles, permissions []string) (string, error) {
	if len(j.SecretKey) == 0 {
		return "", errors.New("JWT secret not initialized")
	}

	claims := models.JWTClaims{
		UserID:      userID,
		Username:    username,
		Roles:       roles,
		Permissions: permissions,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "auth-service",
			Subject:   userID.String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.SecretKey)
}

// ValidateJWT validates a JWT token and returns its claims
func (j *JWTWrapper) ValidateJWT(tokenString string) (*models.JWTClaims, error) {
	if len(j.SecretKey) == 0 {
		return nil, errors.New("JWT secret not initialized")
	}

	token, err := jwt.ParseWithClaims(tokenString, &models.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return j.SecretKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if claims, ok := token.Claims.(*models.JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token claims")
}

// AuthMiddleware is a Gin middleware to validate JWT tokens
func (j *JWTWrapper) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		tokenString := authHeader
		if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
			tokenString = tokenString[7:]
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
			c.Abort()
			return
		}

		claims, err := j.ValidateJWT(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token: " + err.Error()})
			c.Abort()
			return
		}

		// Set claims in context for downstream handlers
		c.Set("userID", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("userRoles", claims.Roles)
		c.Set("userPermissions", claims.Permissions)
		c.Set("jwtClaims", claims)

		c.Next()
	}
}
