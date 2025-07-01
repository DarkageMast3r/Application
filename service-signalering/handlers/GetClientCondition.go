package handlers

import (
	"net/http"
	"service-signalering/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetClientCondition(c *gin.Context) {
	clientIDStr := c.Param("id")

	// Parse client ID
	clientID, err := uuid.Parse(clientIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Code:    "INVALID_CLIENT_ID",
			Message: "Invalid client ID format",
		})
		return
	}

	response := models.ToestandResponse{
		ToestandID:          uuid.New(),
		ClientID:            clientID,
		Status:              "actief",
		TijdstipRegistratie: time.Now(),
		Signalen: []models.Signaal{
			{
				Type:     "hartslag",
				Waarde:   72.5,
				Tijdstip: time.Now().Add(-1 * time.Hour),
				Bron:     "sensor",
			},
		},
		Classificatie: &models.ToestandClassificatie{
			Categorie: "cardiovasculair",
			Ernst:     "normaal",
			Motivatie: "Vitalen binnen normale grenzen",
		},
	}

	c.JSON(http.StatusOK, response)
}
