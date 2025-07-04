package api

import (
	"ZorgTechImplementatie/pkg/cache"
	"ZorgTechImplementatie/pkg/database"
	"ZorgTechImplementatie/pkg/models"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ImplementatieRepository interface {
	Healthcheck(c *gin.Context)
	AanvraagProduct(c *gin.Context)
	OntvangProduct(c *gin.Context)
	InstalleerProduct(c *gin.Context)
	PersonaliseerProduct(c *gin.Context)
	LeverProduct(c *gin.Context)
	MarkeerAlsGeimplementeerd(c *gin.Context)
	GetProductInformatie(c *gin.Context)
	GetImplementatieStatus(c *gin.Context)
	GetInstallatieStatus(c *gin.Context)
	GetPersoonlijkeInstellingen(c *gin.Context)
}

type implementatieRepository struct {
	DB          database.Database
	RedisClient cache.Cache
	Ctx         *context.Context
}

func NewImplementatieRepository(db database.Database, redisClient cache.Cache, ctx *context.Context) *implementatieRepository {
	return &implementatieRepository{
		DB:          db,
		RedisClient: redisClient,
		Ctx:         ctx,
	}
}

// @BasePath /api/v1

// Healthcheck godoc
// @Summary ping example
// @Schemes
// @Description do ping
// @Tags example
// @Accept json
// @Produce json
// @Success 200 {string} ok
// @Router / [get]
func (r *implementatieRepository) Healthcheck(c *gin.Context) {
	c.JSON(http.StatusOK, "ok")
}

// AanvraagProduct godoc
// @Summary Registreer een productaanvraag
// @Description Registreert dat een product is aangevraagd voor een cliënt
// @Tags implementatie
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param input body models.AanvraagProductCommand true "Aanvraag gegevens"
// @Success 201 {object} models.ImplementatieDossier "Successfully created aanvraag"
// @Failure 400 {string} string "Bad Request"
// @Router /implementatie/aanvraag [post]
func (r *implementatieRepository) AanvraagProduct(c *gin.Context) {
	var input models.AanvraagProductCommand
	fmt.Println("Received request")
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Maak nieuw implementatiedossier aan
	dossier := models.ImplementatieDossier{
		ImplementatieID: uuid.New(),
		ClientID:        input.ClientID,
		ZorgtechID:      input.ZorgtechID,
		Status:          models.StatusBesteld,
		Logs: []models.ImplementatieLog{
			{
				Actie:          "Aanvraag geregistreerd",
				UitgevoerdDoor: "Systeem",
			},
		},
	}

	fmt.Println("Received request 2")
	// Sla op in database
	if err := r.DB.Create(&dossier).Error(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Kon dossier niet aanmaken"})
		return
	}
	fmt.Println("Received request 3")

	// TODO: Roep ReserveerBudget aan bij BC Financiëring

	c.JSON(http.StatusCreated, gin.H{"data": dossier})
}

// OntvangProduct godoc
// @Summary Registreer ontvangst van product
// @Description Registreert dat het product is ontvangen bij de instelling
// @Tags implementatie
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param input body models.OntvangProductCommand true "Ontvangst gegevens"
// @Success 200 {object} models.ImplementatieDossier "Successfully updated dossier"
// @Failure 400 {string} string "Bad Request"
// @Failure 404 {string} string "Dossier niet gevonden"
// @Router /implementatie/ontvang [post]
func (r *implementatieRepository) OntvangProduct(c *gin.Context) {
	var input models.OntvangProductCommand
	var dossier models.ImplementatieDossier

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Zoek dossier op basis van client en product
	if err := r.DB.Where("client_id = ? AND zorgtech_id = ?", input.ClientID, input.ZorgtechID).
		First(&dossier).Error(); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Dossier niet gevonden"})
		return
	}

	// Update dossier status en serienummer
	dossier.Status = models.StatusGeleverd
	dossier.Serienummer = input.Serienummer
	dossier.Logs = append(dossier.Logs, models.ImplementatieLog{
		Actie:          "Product ontvangen",
		UitgevoerdDoor: "Systeem",
	})

	// Sla update op
	if err := r.DB.Save(&dossier).Error(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Kon dossier niet updaten"})
		return
	}

	// Invalideer cache voor dit dossier
	cacheKey := "implementatie_status_" + dossier.ClientID.String()
	r.RedisClient.Del(*r.Ctx, cacheKey)

	c.JSON(http.StatusOK, gin.H{"data": dossier})
}

// InstalleerProduct godoc
// @Summary Voer installatie uit
// @Description Registreert dat het product is geïnstalleerd bij de cliënt
// @Tags implementatie
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param input body models.InstalleerProductCommand true "Installatie gegevens"
// @Success 200 {object} models.ImplementatieDossier "Successfully updated dossier"
// @Failure 400 {string} string "Bad Request"
// @Failure 404 {string} string "Dossier niet gevonden"
// @Failure 409 {string} string "Conflict - Product nog niet geleverd"
// @Router /implementatie/installeer [post]
func (r *implementatieRepository) InstalleerProduct(c *gin.Context) {
	var input models.InstalleerProductCommand
	var dossier models.ImplementatieDossier

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Zoek dossier op basis van client en serienummer
	if err := r.DB.Where("client_id = ? AND serienummer = ?", input.ClientID, input.Serienummer).
		First(&dossier).Error(); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Dossier niet gevonden"})
		return
	}

	// Valideer business rules
	if dossier.Status != models.StatusGeleverd {
		c.JSON(http.StatusConflict, gin.H{"error": "Product moet eerst geleverd zijn voor installatie"})
		return
	}

	// Update dossier status en installatiedatum
	now := time.Now()
	dossier.Status = models.StatusGeinstalleerd
	dossier.InstallatieDatum = &now
	dossier.Logs = append(dossier.Logs, models.ImplementatieLog{
		Actie:          "Product geïnstalleerd",
		UitgevoerdDoor: "Systeem",
	})

	// Sla update op
	if err := r.DB.Save(&dossier).Error(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Kon dossier niet updaten"})
		return
	}

	// Invalideer cache voor dit dossier
	cacheKey := "implementatie_status_" + dossier.ClientID.String()
	r.RedisClient.Del(*r.Ctx, cacheKey)

	c.JSON(http.StatusOK, gin.H{"data": dossier})
}

// PersonaliseerProduct godoc
// @Summary Personaliseer product voor cliënt
// @Description Past het product aan op zorgvraag of voorkeuren van de cliënt
// @Tags implementatie
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param input body models.PersonaliseerProductCommand true "Personalisatie gegevens"
// @Success 200 {object} models.ImplementatieDossier "Successfully updated dossier"
// @Failure 400 {string} string "Bad Request"
// @Failure 404 {string} string "Dossier niet gevonden"
// @Failure 409 {string} string "Conflict - Product nog niet geïnstalleerd"
// @Router /implementatie/personaliseer [post]
func (r *implementatieRepository) PersonaliseerProduct(c *gin.Context) {
	var input models.PersonaliseerProductCommand
	var dossier models.ImplementatieDossier

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Zoek dossier op basis van client ID
	if err := r.DB.Where("client_id = ?", input.ClientID).
		First(&dossier).Error(); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Dossier niet gevonden"})
		return
	}

	// Valideer business rules
	if dossier.Status != models.StatusGeinstalleerd {
		c.JSON(http.StatusConflict, gin.H{"error": "Product moet eerst geïnstalleerd zijn voor personalisatie"})
		return
	}

	// Update dossier status en personalisatie
	dossier.Status = models.StatusGepersonaliseerd
	dossier.Personalisatie = &input.Instellingen
	dossier.Logs = append(dossier.Logs, models.ImplementatieLog{
		Actie:          "Product gepersonaliseerd",
		UitgevoerdDoor: "Systeem",
	})

	// Sla update op
	if err := r.DB.Save(&dossier).Error(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Kon dossier niet updaten"})
		return
	}

	// Invalideer cache voor dit dossier
	cacheKey := "implementatie_status_" + dossier.ClientID.String()
	r.RedisClient.Del(*r.Ctx, cacheKey)
	cacheKey = "persoonlijke_instellingen_" + dossier.ClientID.String()
	r.RedisClient.Del(*r.Ctx, cacheKey)

	c.JSON(http.StatusOK, gin.H{"data": dossier})
}

// LeverProduct godoc
// @Summary Registreer levering aan cliënt
// @Description Registreert dat het product is geleverd bij de cliënt
// @Tags implementatie
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param input body models.LeverProductCommand true "Leveringsgegevens"
// @Success 200 {object} models.ImplementatieDossier "Successfully updated dossier"
// @Failure 400 {string} string "Bad Request"
// @Failure 404 {string} string "Dossier niet gevonden"
// @Router /implementatie/lever [post]
func (r *implementatieRepository) LeverProduct(c *gin.Context) {
	var input models.LeverProductCommand
	var dossier models.ImplementatieDossier

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Zoek dossier op basis van client en product
	if err := r.DB.Where("client_id = ? AND zorgtech_id = ?", input.ClientID, input.ZorgtechID).
		First(&dossier).Error(); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Dossier niet gevonden"})
		return
	}

	// Update dossier logs
	dossier.Logs = append(dossier.Logs, models.ImplementatieLog{
		Actie:          "Product geleverd aan cliënt",
		UitgevoerdDoor: "Systeem",
	})

	// Sla update op
	if err := r.DB.Save(&dossier).Error(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Kon dossier niet updaten"})
		return
	}

	// Invalideer cache voor dit dossier
	cacheKey := "implementatie_status_" + dossier.ClientID.String()
	r.RedisClient.Del(*r.Ctx, cacheKey)

	c.JSON(http.StatusOK, gin.H{"data": dossier})
}

// MarkeerAlsGeimplementeerd godoc
// @Summary Markeer implementatie als voltooid
// @Description Sluit het implementatieproces formeel af
// @Tags implementatie
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param input body models.MarkeerAlsGeimplementeerdCommand true "Client ID"
// @Success 200 {object} models.ImplementatieDossier "Successfully updated dossier"
// @Failure 400 {string} string "Bad Request"
// @Failure 404 {string} string "Dossier niet gevonden"
// @Failure 409 {string} string "Conflict - Implementatie nog niet klaar voor afronding"
// @Router /implementatie/voltooi [post]
func (r *implementatieRepository) MarkeerAlsGeimplementeerd(c *gin.Context) {
	var input models.MarkeerAlsGeimplementeerdCommand
	var dossier models.ImplementatieDossier

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Zoek dossier op basis van client ID
	if err := r.DB.Where("client_id = ?", input.ClientID).
		First(&dossier).Error(); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Dossier niet gevonden"})
		return
	}

	// Valideer business rules
	if dossier.Status != models.StatusGepersonaliseerd {
		c.JSON(http.StatusConflict, gin.H{"error": "Implementatie kan alleen worden voltooid na personalisatie"})
		return
	}

	// Update dossier status
	dossier.Status = models.StatusVoltooid
	dossier.Logs = append(dossier.Logs, models.ImplementatieLog{
		Actie:          "Implementatie voltooid",
		UitgevoerdDoor: "Systeem",
	})

	// Sla update op
	if err := r.DB.Save(&dossier).Error(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Kon dossier niet updaten"})
		return
	}

	// Invalideer cache voor dit dossier
	cacheKey := "implementatie_status_" + dossier.ClientID.String()
	r.RedisClient.Del(*r.Ctx, cacheKey)

	c.JSON(http.StatusOK, gin.H{"data": dossier})
}

// GetProductInformatie godoc
// @Summary Haal productinformatie op
// @Description Toont de productinformatie van een gespecificeerd product
// @Tags implementatie
// @Security ApiKeyAuth
// @Produce json
// @Param zorgtechId path string true "ZorgTech Product ID"
// @Success 200 {object} models.ZorgTechProduct "Successfully retrieved product"
// @Failure 404 {string} string "Product niet gevonden"
// @Router /implementatie/product/{zorgtechId} [get]
func (r *implementatieRepository) GetProductInformatie(c *gin.Context) {
	var product models.ZorgTechProduct
	zorgtechID := c.Param("zorgtechId")

	// Probeer eerst cache
	cacheKey := "product_info_" + zorgtechID
	cachedProduct, err := r.RedisClient.Get(*r.Ctx, cacheKey).Result()
	if err == nil {
		var cached models.ZorgTechProduct
		if err := json.Unmarshal([]byte(cachedProduct), &cached); err == nil {
			c.JSON(http.StatusOK, gin.H{"data": cached})
			return
		}
	}

	// Zoek product in database
	if err := r.DB.Where("zorgtech_id = ?", zorgtechID).First(&product).Error(); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product niet gevonden"})
		return
	}

	// Sla in cache op
	serialized, err := json.Marshal(product)
	if err == nil {
		r.RedisClient.Set(*r.Ctx, cacheKey, serialized, time.Hour)
	}

	c.JSON(http.StatusOK, gin.H{"data": product})
}

// GetImplementatieStatus godoc
// @Summary Haal implementatiestatus op
// @Description Toont de voortgang en status van de implementatie voor een cliënt
// @Tags implementatie
// @Security ApiKeyAuth
// @Produce json
// @Param clientId path string true "Client ID"
// @Success 200 {object} models.ImplementatieDossier "Successfully retrieved status"
// @Failure 404 {string} string "Dossier niet gevonden"
// @Router /implementatie/status/{clientId} [get]
func (r *implementatieRepository) GetImplementatieStatus(c *gin.Context) {
	var dossier models.ImplementatieDossier
	clientID := c.Param("clientId")

	// Probeer eerst cache
	cacheKey := "implementatie_status_" + clientID
	cachedStatus, err := r.RedisClient.Get(*r.Ctx, cacheKey).Result()
	if err == nil {
		var cached models.ImplementatieDossier
		if err := json.Unmarshal([]byte(cachedStatus), &cached); err == nil {
			c.JSON(http.StatusOK, gin.H{"data": cached})
			return
		}
	}

	// Zoek dossier in database
	if err := r.DB.Where("client_id = ?", clientID).First(&dossier).Error(); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Dossier niet gevonden"})
		return
	}

	// Sla in cache op
	serialized, err := json.Marshal(dossier)
	if err == nil {
		r.RedisClient.Set(*r.Ctx, cacheKey, serialized, 30*time.Minute)
	}

	c.JSON(http.StatusOK, gin.H{"data": dossier})
}

// GetInstallatieStatus godoc
// @Summary Haal installatiestatus op
// @Description Toont de voortgang en status van de installatie en personalisatie voor een cliënt
// @Tags implementatie
// @Security ApiKeyAuth
// @Produce json
// @Param clientId path string true "Client ID"
// @Success 200 {object} object "Successfully retrieved status"
// @Failure 404 {string} string "Dossier niet gevonden"
// @Router /implementatie/installatie-status/{clientId} [get]
func (r *implementatieRepository) GetInstallatieStatus(c *gin.Context) {
	var dossier models.ImplementatieDossier
	clientID := c.Param("clientId")

	// Zoek dossier in database
	if err := r.DB.Where("client_id = ?", clientID).First(&dossier).Error(); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Dossier niet gevonden"})
		return
	}

	response := gin.H{
		"status":             dossier.Status,
		"installatieDatum":   dossier.InstallatieDatum,
		"isGepersonaliseerd": dossier.Status == models.StatusGepersonaliseerd || dossier.Status == models.StatusVoltooid,
	}

	c.JSON(http.StatusOK, gin.H{"data": response})
}

// GetPersoonlijkeInstellingen godoc
// @Summary Haal persoonlijke instellingen op
// @Description Haalt de opgeslagen gebruikersinstellingen op
// @Tags implementatie
// @Security ApiKeyAuth
// @Produce json
// @Param clientId path string true "Client ID"
// @Success 200 {object} models.PersonalisatieInstellingen "Successfully retrieved instellingen"
// @Failure 404 {string} string "Dossier niet gevonden of instellingen niet beschikbaar"
// @Router /implementatie/instellingen/{clientId} [get]
func (r *implementatieRepository) GetPersoonlijkeInstellingen(c *gin.Context) {
	var dossier models.ImplementatieDossier
	clientID := c.Param("clientId")

	// Probeer eerst cache
	cacheKey := "persoonlijke_instellingen_" + clientID
	cachedSettings, err := r.RedisClient.Get(*r.Ctx, cacheKey).Result()
	if err == nil {
		var cached models.PersonalisatieInstellingen
		if err := json.Unmarshal([]byte(cachedSettings), &cached); err == nil {
			c.JSON(http.StatusOK, gin.H{"data": cached})
			return
		}
	}

	// Zoek dossier in database
	if err := r.DB.Where("client_id = ?", clientID).First(&dossier).Error(); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Dossier niet gevonden"})
		return
	}

	// Controleer of personalisatie is uitgevoerd
	if dossier.Personalisatie == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Geen persoonlijke instellingen beschikbaar"})
		return
	}

	// Sla in cache op
	serialized, err := json.Marshal(dossier.Personalisatie)
	if err == nil {
		r.RedisClient.Set(*r.Ctx, cacheKey, serialized, time.Hour)
	}

	c.JSON(http.StatusOK, gin.H{"data": dossier.Personalisatie})
}
