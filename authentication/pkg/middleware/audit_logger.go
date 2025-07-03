package middleware

import (
	"authentication/pkg/database"
	"authentication/pkg/models"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// middleware/audit_logger.go
func AuditLogger(db database.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip voor health checks etc.
		if c.Request.URL.Path == "/health" {
			c.Next()
			return
		}

		c.Next()

		action := getActionFromContext(c)
		if action == "" {
			return
		}

		var userID *uuid.UUID
		if claims, exists := c.Get("jwt_claims"); exists {
			jwtClaims := claims.(*models.JWTClaims)
			userID = &jwtClaims.UserID
		}

		log := models.AuditLog{
			ID:        uuid.New(),
			Action:    action,
			UserID:    userID,
			IPAddress: c.ClientIP(),
			UserAgent: c.Request.UserAgent(),
			Metadata:  getMetadataFromContext(c),
			CreatedAt: time.Now(),
		}

		go func() {
			if err := db.Create(&log).Error(); err != nil {
				fmt.Printf("Failed to log audit: %v\n", err)
			}
		}()
	}
}

// Helper functies
func getActionFromContext(c *gin.Context) models.AuditAction {
	switch c.Request.URL.Path {
	case "/auth/login":
		return models.ActionLogin
	case "/auth/logout":
		return models.ActionLogout
	// Voeg meer routes toe
	default:
		return ""
	}
}

func getMetadataFromContext(c *gin.Context) map[string]interface{} {
	metadata := make(map[string]interface{})
	// Voor password reset:
	if c.Request.URL.Path == "/auth/reset-password" {
		metadata["reset_for"] = c.Query("token")
	}
	return metadata
}
