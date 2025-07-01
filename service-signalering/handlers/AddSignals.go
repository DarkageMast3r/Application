package handlers

import (
	"log"
	"net/http"
	"service-signalering/database"
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

	for _, signaal := range request.Signalen {
		query := `
            INSERT INTO signals (client_id, type, waarde, tijdstip, bron) 
            VALUES ($1, $2, $3, $4, $5)`

		_, err := database.DB.Exec(query, clientID, signaal.Type, signaal.Waarde, signaal.Tijdstip, signaal.Bron)
		if err != nil {
			log.Printf("Error saving signal: %v", err)
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{
				Code:    "DATABASE_ERROR",
				Message: "Failed to save signal",
			})
			return
		}
	}

	log.Printf("Successfully saved %d signals for client %s", len(request.Signalen), clientID)

	c.JSON(http.StatusCreated, gin.H{
		"message":     "Signalen succesvol geregistreerd",
		"client_id":   clientID,
		"signalen":    request.Signalen,
		"saved_count": len(request.Signalen),
	})
}
