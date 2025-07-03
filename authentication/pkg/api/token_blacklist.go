package api

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func (r *authRepository) Logout(c *gin.Context) {
	claims, exists := c.Get("jwt_claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	jwtClaims := claims.(*models.JWTClaims)
	expiresAt := time.Unix(jwtClaims.ExpiresAt, 0)

	// Voeg token toe aan blacklist
	blacklisted := models.TokenBlacklist{
		Token:     c.GetHeader("Authorization")[7:], // Strip "Bearer "
		ExpiresAt: expiresAt,
	}
	if err := r.DB.Create(&blacklisted).Error(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not logout"})
		return
	}

	// Verwijder refresh token
	r.DB.Where("user_id = ?", jwtClaims.UserID).Delete(&models.RefreshToken{})

	c.JSON(http.StatusOK, gin.H{"message": "Successfully logged out"})
}

// Pas token validatie aan
func (r *authRepository) validateToken(token string) (*models.JWTClaims, error) {
	// Check blacklist eerst
	var blacklisted models.TokenBlacklist
	if err := r.DB.Where("token = ?", token).First(&blacklisted).Error(); err == nil {
		return nil, errors.New("token is blacklisted")
	}

	claims, err := models.ValidateJWT(token)
	if err != nil {
		return nil, err
	}

	return claims, nil
}
