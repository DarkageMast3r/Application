package handlers

import (
	"net/http"
	"service-signalering/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

	c.JSON(http.StatusOK, gin.H{
		"message":       "Situatie succesvol geclassificeerd",
		"client_id":     clientID,
		"classificatie": request.Classificatie,
	})
}
