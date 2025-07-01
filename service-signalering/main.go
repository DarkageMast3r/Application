package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"service-signalering/models"

	"github.com/google/uuid"
)

func main() {
	fmt.Println("Testing ZorgSignalering Models...")

	// Test Signaal
	signaal := models.Signaal{
		Type:     "bloeddruk",
		Waarde:   140.5,
		Tijdstip: time.Now(),
		Bron:     "sensor_001",
	}

	// Test ToestandClassificatie
	classificatie := models.ToestandClassificatie{
		Categorie: "Cardiovasculair risico",
		Ernst:     "matig",
		Motivatie: "Bloeddruk verhoogd",
	}

	// Test Beoordeling
	beoordeling := models.Beoordeling{
		Conclusie:       "Monitoring verhogen",
		Urgentie:        "laag",
		GevalideerdDoor: "verpleegkundige.jansen",
		Tijdstip:        time.Now(),
	}

	// Test ClientToestand
	clientID := uuid.New()
	toestand := models.ClientToestand{
		ToestandID:          uuid.New(),
		ClientID:            clientID,
		Signalen:            []models.Signaal{signaal},
		Classificatie:       &classificatie,
		Beoordeling:         &beoordeling,
		Status:              "beoordeeld",
		TijdstipRegistratie: time.Now(),
	}

	// Test JSON serialization
	jsonData, err := json.MarshalIndent(toestand, "", "  ")
	if err != nil {
		log.Fatal("Error marshaling JSON:", err)
	}

	fmt.Println("\nClientToestand as JSON:")
	fmt.Println(string(jsonData))

	// Test Request structs
	req := models.RegistreerAchteruitgangRequest{
		Signalen: []models.Signaal{signaal},
	}

	reqJSON, _ := json.MarshalIndent(req, "", "  ")
	fmt.Println("\nRegistreerAchteruitgangRequest as JSON:")
	fmt.Println(string(reqJSON))

	// Test Response structs
	response := models.ToestandResponse{
		ToestandID:          toestand.ToestandID,
		ClientID:            toestand.ClientID,
		Status:              toestand.Status,
		TijdstipRegistratie: toestand.TijdstipRegistratie,
		Signalen:            toestand.Signalen,
		Classificatie:       toestand.Classificatie,
		Beoordeling:         toestand.Beoordeling,
	}

	respJSON, _ := json.MarshalIndent(response, "", "  ")
	fmt.Println("\nToestandResponse as JSON:")
	fmt.Println(string(respJSON))

	fmt.Println("Woohoo het werkt")
}
