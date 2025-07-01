package handlers

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log"
	"net/http"
	"service-signalering/database"
	"service-signalering/models"
	"time"
)

func GetClientCondition(c *gin.Context) {
	clientIDStr := c.Param("id")
	clientID, err := uuid.Parse(clientIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Code:    "INVALID_CLIENT_ID",
			Message: "Invalid client ID format",
		})
		return
	}

	// Get latest signals for this client
	signalsQuery := `
        SELECT type, waarde, tijdstip, bron 
        FROM signals 
        WHERE client_id = $1 
        ORDER BY tijdstip DESC 
        LIMIT 10`

	rows, err := database.DB.Query(signalsQuery, clientID)
	if err != nil {
		log.Printf("Error querying signals: %v", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Code:    "DATABASE_ERROR",
			Message: "Failed to retrieve signals",
		})
		return
	}
	defer rows.Close()

	var signalen []models.Signaal
	for rows.Next() {
		var signaal models.Signaal
		err := rows.Scan(&signaal.Type, &signaal.Waarde, &signaal.Tijdstip, &signaal.Bron)
		if err != nil {
			log.Printf("Error scanning signal: %v", err)
			continue
		}
		signalen = append(signalen, signaal)
	}

	// Get latest classification for this client
	var classificatie *models.ToestandClassificatie
	classificationQuery := `
        SELECT categorie, ernst, motivatie 
        FROM classifications 
        WHERE client_id = $1 
        ORDER BY created_at DESC 
        LIMIT 1`

	var cat, ernst, motivatie string
	err = database.DB.QueryRow(classificationQuery, clientID).Scan(&cat, &ernst, &motivatie)
	if err != nil && err != sql.ErrNoRows {
		log.Printf("Error querying classification: %v", err)
	} else if err == nil {
		classificatie = &models.ToestandClassificatie{
			Categorie: cat,
			Ernst:     ernst,
			Motivatie: motivatie,
		}
	}

	// Get latest assessment for this client
	var beoordeling *models.Beoordeling
	assessmentQuery := `
        SELECT conclusie, urgentie, gevalideerd_door, tijdstip 
        FROM assessments 
        WHERE client_id = $1 
        ORDER BY created_at DESC 
        LIMIT 1`

	var conclusie, urgentie, gevalideerdDoor string
	var tijdstip time.Time
	err = database.DB.QueryRow(assessmentQuery, clientID).Scan(&conclusie, &urgentie, &gevalideerdDoor, &tijdstip)
	if err != nil && err != sql.ErrNoRows {
		log.Printf("Error querying assessment: %v", err)
	} else if err == nil {
		beoordeling = &models.Beoordeling{
			Conclusie:       conclusie,
			Urgentie:        urgentie,
			GevalideerdDoor: gevalideerdDoor,
			Tijdstip:        tijdstip,
		}
	}

	// Create response with real data from database
	response := models.ToestandResponse{
		ToestandID:          uuid.New(), // Generate new UUID for this request
		ClientID:            clientID,
		Status:              "actief",
		TijdstipRegistratie: time.Now(),
		Signalen:            signalen,
		Classificatie:       classificatie,
		Beoordeling:         beoordeling,
	}

	c.JSON(http.StatusOK, response)
}
