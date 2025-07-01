package main

import (
	"fmt"
	"net/http"
	"service-signalering/database"
	"service-signalering/handlers"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func main() {

	database.Init()

	// Maak nep-UUIDs voor makkelijk testen
	fmt.Println("=== Test UUIDs ===")
	for i := 1; i <= 10; i++ {
		testUUID := uuid.New()
		fmt.Printf("Client %d: %s\n", i, testUUID.String())
	}
	fmt.Println("=========================")
	fmt.Println()

	r := gin.Default()

	// Check of de API werkt
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"message": "Ik ben er! Geen zorgen!",
		})
	})

	// Daadwerkelijke endpoints
	v1 := r.Group("/api/v1")
	{
		// Client monitoring endpoints
		v1.GET("/clients/:id/condition", handlers.GetClientCondition)
		v1.POST("/clients/:id/signals", handlers.AddSignals)
		v1.POST("/clients/:id/classify", handlers.ClassifyCondition)
		v1.POST("/clients/:id/assess", handlers.AssessCondition)
	}

	r.Run(":8080")
}
