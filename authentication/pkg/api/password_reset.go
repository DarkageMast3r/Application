package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// api/auth.go
func (r *authRepository) RequestPasswordReset(c *gin.Context) {
	var input struct {
		Email string `json:"email" binding:"required,email"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := r.DB.Where("email = ?", input.Email).First(&user).Error(); err != nil {
		// Liever geen foutmelding geven of het email bestaat (security)
		c.JSON(http.StatusOK, gin.H{"message": "If the email exists, a reset link has been sent"})
		return
	}

	token := uuid.New().String()
	resetToken := models.PasswordResetToken{
		Token:     token,
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(1 * time.Hour),
	}

	if err := r.DB.Create(&resetToken).Error(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create reset token"})
		return
	}

	// In productie: Stuur email met reset link
	resetLink := fmt.Sprintf("https://arjanslab.nl/reset-password?token=%s", token)
	fmt.Printf("DEV: Password reset link: %s\n", resetLink)

	c.JSON(http.StatusOK, gin.H{"message": "Reset link sent"})
}

func (r *authRepository) ResetPassword(c *gin.Context) {
	var input struct {
		Token    string `json:"token" binding:"required"`
		Password string `json:"password" binding:"required,min=8"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var resetToken models.PasswordResetToken
	if err := r.DB.Where("token = ? AND used = ? AND expires_at > ?",
		input.Token, false, time.Now()).First(&resetToken).Error(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid or expired token"})
		return
	}

	hashedPassword, err := models.NewPassword(input.Password).Hash()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not hash password"})
		return
	}

	// Update user password
	if err := r.DB.Model(&models.User{}).Where("id = ?", resetToken.UserID).
		Update("password_hash", hashedPassword).Error(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update password"})
		return
	}

	// Mark token as used
	r.DB.Model(&resetToken).Update("used", true)

	// Invalidate all sessions
	r.DB.Where("user_id = ?", resetToken.UserID).Delete(&models.RefreshToken{})

	c.JSON(http.StatusOK, gin.H{"message": "Password reset successfully"})
}
