package handlers

import (
	"net/http"
	"service-signalering/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func AssessCondition(c *gin.Context) {
	clientIDStr := c.Param("id")

	clientID, err := uuid.Parse(clientIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Code:    "INVALID_CLIENT_ID",
			Message: "Invalid client ID format",
		})
		return
	}

	var request models.BeoordeelSituatieRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Code:    "INVALID_REQUEST",
			Message: "Invalid request format",
		})
		return
	}

	assessment := models.Beoordeling{
		Conclusie:       request.Conclusie,
		Urgentie:        request.Urgentie,
		GevalideerdDoor: request.GevalideerdDoor,
		Tijdstip:        time.Now(),
	}

	response := models.BeoordelingResponse{
		ClientID:         clientID,
		LaatsBeoordeling: &assessment,
		Aanbevelingen:    []string{"Monitor vitalen", "Volg medicatieregime"},
	}

	c.JSON(http.StatusOK, response)
}
