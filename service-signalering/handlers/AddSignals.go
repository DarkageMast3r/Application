package handlers

import (
	"net/http"
	"service-signalering/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func AddSignals(c *gin.Context) {
	clientIDStr := c.Param("id")

	clientID, err := uuid.Parse(clientIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Code:    "INVALID_CLIENT_ID",
			Message: "Invalid client ID format",
		})
		return
	}

	var request models.RegistreerAchteruitgangRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Code:    "INVALID_REQUEST",
			Message: "Invalid request format",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":   "Signalen succesvol geregistreerd",
		"client_id": clientID,
		"signalen":  request.Signalen,
	})
}
