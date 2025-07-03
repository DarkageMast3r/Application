package middleware

import (
	"authentication/pkg/models"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

var jwtSecret []byte // Should be loaded from environment variables or a secure secret manager

func InitJWT(secret string) {
	jwtSecret = []byte(secret)
}

// GenerateJWT generates a new JWT token
func GenerateJWT(userID uuid.UUID, username string, roles, permissions []string) (string, error) {
	if jwtSecret == nil {
		return "", errors.New("JWT secret not initialized")
	}

	claims := models.JWTClaims{
		UserID:      userID,
		Username:    username,
		Roles:       roles,
		Permissions: permissions,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * 15).Unix(), // Access token valid for 15 minutes
			IssuedAt:  time.Now().Unix(),
			Issuer:    "nietgrappig",
			Subject:   userID.String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", errors.New("failed to sign token")
	}
	return tokenString, nil
}

// ValidateJWT validates a JWT token and returns its claims
func ValidateJWT(tokenString string) (*models.JWTClaims, error) {
	if jwtSecret == nil {
		return nil, errors.New("JWT secret not initialized")
	}

	token, err := jwt.ParseWithClaims(tokenString, &models.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*models.JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// AuthMiddleware is a Gin middleware to validate JWT tokens
func AuthMiddleware(jwtSecret string) gin.HandlerFunc {
	InitJWT(jwtSecret) // Initialize secret for this middleware instance
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

		claims, err := ValidateJWT(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token: " + err.Error()})
			c.Abort()
			return
		}

		// Zet de claims in de context voor verdere handlers
		c.Set("userID", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("userRoles", claims.Roles)
		c.Set("userPermissions", claims.Permissions)
		c.Set("jwtClaims", claims) // Optioneel: de volledige claims struct

		c.Next()
	}
}
