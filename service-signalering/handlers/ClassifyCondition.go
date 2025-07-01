package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log"
	"net/http"
	"service-signalering/database"
	"service-signalering/models"
)

func ClassifyCondition(c *gin.Context) {
	clientIDStr := c.Param("id")
	clientID, err := uuid.Parse(clientIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Code:    "INVALID_CLIENT_ID",
			Message: "Invalid client ID format",
		})
		return
	}

	var request models.ClassificeerSituatieRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Code:    "INVALID_REQUEST",
			Message: "Invalid request format",
		})
		return
	}

	// Save classification to database
	query := `
        INSERT INTO classifications (client_id, categorie, ernst, motivatie) 
        VALUES ($1, $2, $3, $4)`

	_, err = database.DB.Exec(query,
		clientID,
		request.Classificatie.Categorie,
		request.Classificatie.Ernst,
		request.Classificatie.Motivatie)

	if err != nil {
		log.Printf("Error saving classification: %v", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Code:    "DATABASE_ERROR",
			Message: "Failed to save classification",
		})
		return
	}

	log.Printf("Successfully saved classification for client %s", clientID)

	c.JSON(http.StatusOK, gin.H{
		"message":       "Situatie succesvol geclassificeerd",
		"client_id":     clientID,
		"classificatie": request.Classificatie,
	})
}
