package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log"
	"net/http"
	"service-signalering/models"
	"service-signalering/pkg/database"
	"service-signalering/pkg/validation"
	"time"
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

	validationErrors := validation.ValidateBeoordelingRequest(request)
	if len(validationErrors) > 0 {
		log.Printf("Assessment validation failed for client %s: %v", clientID, validationErrors)

		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Code:    "VALIDATION_FAILED",
			Message: "Beoordeling validatie gefaald",
			Details: validationErrors,
		})
		return
	}

	assessmentTime := time.Now()
	assessment := models.Beoordeling{
		Conclusie:       request.Conclusie,
		Urgentie:        request.Urgentie,
		GevalideerdDoor: request.GevalideerdDoor,
		Tijdstip:        assessmentTime,
	}

	query := `
        INSERT INTO assessments (client_id, conclusie, urgentie, gevalideerd_door, tijdstip) 
        VALUES ($1, $2, $3, $4, $5)`

	_, err = database.DB.Exec(query,
		clientID,
		assessment.Conclusie,
		assessment.Urgentie,
		assessment.GevalideerdDoor,
		assessment.Tijdstip)

	if err != nil {
		log.Printf("Error saving assessment: %v", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Code:    "DATABASE_ERROR",
			Message: "Failed to save assessment",
		})
		return
	}

	log.Printf("Successfully saved assessment for client %s by %s", clientID, assessment.GevalideerdDoor)

	response := models.BeoordelingResponse{
		ClientID:         clientID,
		LaatsBeoordeling: &assessment,
		Aanbevelingen:    []string{"Monitor vitalen", "Volg medicatieregime"},
	}

	c.JSON(http.StatusOK, response)
}
