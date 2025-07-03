package api

import (
	"ZorgTechCatalogus/pkg/cache"
	"ZorgTechCatalogus/pkg/database"
	models "ZorgTechCatalogus/pkg/models"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CatalogusRepository interface {
	Healthcheck(c *gin.Context)
	MaakZorgTechProduct(c *gin.Context)
	WijzigZorgTechProduct(c *gin.Context)
	VerwijderZorgTechProduct(c *gin.Context)
	VoegTechnischDetailToe(c *gin.Context)
	VerwijderTechnischDetail(c *gin.Context)
	GetProductById(c *gin.Context)
	FindByCategorie(c *gin.Context)
	ListAlleProducten(c *gin.Context)
	ZoekOpNaam(c *gin.Context)
}

type catalogusRepository struct {
	DB          database.Database
	RedisClient cache.Cache
	Ctx         *context.Context
}

func NewCatalogusRepository(db database.Database, redisClient cache.Cache, ctx *context.Context) *catalogusRepository {
	return &catalogusRepository{
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
func (r *catalogusRepository) Healthcheck(c *gin.Context) {
	c.JSON(http.StatusOK, "ok")
}

// MaakZorgTechProduct godoc
// @Summary Registreer een nieuw zorgtechnologieproduct
// @Description Maakt een nieuw zorgtechnologieproduct aan in de catalogus
// @Tags catalogus
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param input body models.MaakZorgTechProductCommand true "Product gegevens"
// @Success 201 {object} models.ZorgTechProduct "Successfully created product"
// @Failure 400 {string} string "Bad Request"
// @Router /catalogus/product [post]
func (r *catalogusRepository) MaakZorgTechProduct(c *gin.Context) {
	var input models.MaakZorgTechProductCommand

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Maak nieuw product aan
	product := models.ZorgTechProduct{
		ID:                uuid.New(),
		Naam:              input.Naam,
		Beschrijving:      input.Beschrijving,
		Categorie:         input.Categorie,
		TechnischeDetails: input.TechnischeDetails,
		Prijs:             input.Prijs,
		Leverancier:       input.Leverancier,
		IsActief:          true,
	}

	// Sla op in database
	if err := r.DB.Create(&product).Error(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Kon product niet aanmaken"})
		return
	}

	// Registreer domain event
	event := models.ProductEvent{
		EventID:     uuid.New(),
		ProductID:   product.ID,
		Type:        "ZorgTechProductAangemaakt",
		Payload:     product.Naam,
		TriggeredBy: "Catalogusbeheerder",
	}
	if err := r.DB.Create(&event).Error(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Kon event niet registreren"})
		return
	}

	// Invalideer cache
	r.RedisClient.Del(*r.Ctx, "catalogus_alle_producten")
	r.RedisClient.Del(*r.Ctx, "catalogus_categorie_"+product.Categorie)

	c.JSON(http.StatusCreated, gin.H{"data": product})
}

// WijzigZorgTechProduct godoc
// @Summary Wijzig een bestaand zorgtechnologieproduct
// @Description Past de eigenschappen van een bestaand zorgtechnologieproduct aan
// @Tags catalogus
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param input body models.WijzigZorgTechProductCommand true "Wijzigingsgegevens"
// @Success 200 {object} models.ZorgTechProduct "Successfully updated product"
// @Failure 400 {string} string "Bad Request"
// @Failure 404 {string} string "Product niet gevonden"
// @Router /catalogus/product [put]
func (r *catalogusRepository) WijzigZorgTechProduct(c *gin.Context) {
	var input models.WijzigZorgTechProductCommand
	var product models.ZorgTechProduct

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Zoek product op basis van ID
	if err := r.DB.Where("id = ?", input.ProductID).First(&product).Error(); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product niet gevonden"})
		return
	}

	// Track wijzigingen voor event
	wijzigingen := make(map[string]interface{})

	// Pas alleen velden aan die zijn meegegeven
	if input.Naam != nil {
		wijzigingen["naam"] = product.Naam + " -> " + *input.Naam
		product.Naam = *input.Naam
	}
	if input.Beschrijving != nil {
		wijzigingen["beschrijving"] = "beschrijving gewijzigd"
		product.Beschrijving = *input.Beschrijving
	}
	if input.Categorie != nil {
		wijzigingen["categorie"] = product.Categorie + " -> " + *input.Categorie
		product.Categorie = *input.Categorie
	}
	if input.TechnischeDetails != nil {
		wijzigingen["technischeDetails"] = "technische details gewijzigd"
		product.TechnischeDetails = *input.TechnischeDetails
	}
	if input.Prijs != nil {
		wijzigingen["prijs"] = fmt.Sprintf("%.2f -> %.2f", product.Prijs, *input.Prijs)
		product.Prijs = *input.Prijs
	}
	if input.Leverancier != nil {
		wijzigingen["leverancier"] = product.Leverancier + " -> " + *input.Leverancier
		product.Leverancier = *input.Leverancier
	}

	// Sla update op
	if err := r.DB.Save(&product).Error(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Kon product niet updaten"})
		return
	}

	// Registreer domain event
	payloadBytes, err := json.Marshal(wijzigingen)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Kon wijzigingen niet serialiseren"})
		return
	}

	event := models.ProductEvent{
		EventID:     uuid.New(),
		ProductID:   product.ID,
		Type:        "ZorgTechProductGewijzigd",
		Payload:     string(payloadBytes), // JSON string
		TriggeredBy: "Catalogusbeheerder",
	}
	if err := r.DB.Create(&event).Error(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Kon event niet registreren"})
		return
	}

	// Invalideer cache
	r.RedisClient.Del(*r.Ctx, "catalogus_product_"+product.ID.String())
	r.RedisClient.Del(*r.Ctx, "catalogus_alle_producten")
	r.RedisClient.Del(*r.Ctx, "catalogus_categorie_"+product.Categorie)

	c.JSON(http.StatusOK, gin.H{"data": product})
}

// VerwijderZorgTechProduct godoc
// @Summary Archiveer een zorgtechnologieproduct
// @Description Archiveert een product zodat deze niet meer actief is in de catalogus
// @Tags catalogus
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param input body models.VerwijderZorgTechProductCommand true "Product ID"
// @Success 200 {object} models.ZorgTechProduct "Successfully archived product"
// @Failure 400 {string} string "Bad Request"
// @Failure 404 {string} string "Product niet gevonden"
// @Router /catalogus/product [delete]
func (r *catalogusRepository) VerwijderZorgTechProduct(c *gin.Context) {
	var input models.VerwijderZorgTechProductCommand
	var product models.ZorgTechProduct

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Zoek product op basis van ID
	if err := r.DB.Where("id = ?", input.ProductID).First(&product).Error(); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product niet gevonden"})
		return
	}

	// Archiveer product (soft delete)
	product.IsActief = false

	// Sla update op
	if err := r.DB.Save(&product).Error(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Kon product niet archiveren"})
		return
	}

	// Registreer domain event
	event := models.ProductEvent{
		EventID:     uuid.New(),
		ProductID:   product.ID,
		Type:        "ZorgTechProductVerwijderd",
		Payload:     product.Naam,
		TriggeredBy: "Catalogusbeheerder",
	}
	if err := r.DB.Create(&event).Error(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Kon event niet registreren"})
		return
	}

	// Invalideer cache
	r.RedisClient.Del(*r.Ctx, "catalogus_product_"+product.ID.String())
	r.RedisClient.Del(*r.Ctx, "catalogus_alle_producten")
	r.RedisClient.Del(*r.Ctx, "catalogus_categorie_"+product.Categorie)

	c.JSON(http.StatusOK, gin.H{"data": product})
}

// VoegTechnischDetailToe godoc
// @Summary Voeg een technisch detail toe aan een product
// @Description Voegt een nieuw technisch detail toe aan een bestaand product
// @Tags catalogus
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param input body models.VoegTechnischDetailToeCommand true "Technisch detail"
// @Success 200 {object} models.ZorgTechProduct "Successfully updated product"
// @Failure 400 {string} string "Bad Request"
// @Failure 404 {string} string "Product niet gevonden"
// @Router /catalogus/product/technisch-detail [post]
func (r *catalogusRepository) VoegTechnischDetailToe(c *gin.Context) {
	var input models.VoegTechnischDetailToeCommand
	var product models.ZorgTechProduct

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Zoek product op basis van ID
	if err := r.DB.Where("id = ?", input.ProductID).First(&product).Error(); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product niet gevonden"})
		return
	}

	// Voeg nieuw detail toe
	newDetail := models.TechnischDetail{
		Sleutel: input.Sleutel,
		Waarde:  input.Waarde,
	}
	product.TechnischeDetails = append(product.TechnischeDetails, newDetail)

	// Sla update op
	if err := r.DB.Save(&product).Error(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Kon technisch detail niet toevoegen"})
		return
	}

	// Registreer domain event
	event := models.ProductEvent{
		EventID:     uuid.New(),
		ProductID:   product.ID,
		Type:        "TechnischDetailToegevoegd",
		Payload:     input.Sleutel + ": " + input.Waarde,
		TriggeredBy: "Catalogusbeheerder",
	}
	if err := r.DB.Create(&event).Error(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Kon event niet registreren"})
		return
	}

	// Invalideer cache
	r.RedisClient.Del(*r.Ctx, "catalogus_product_"+product.ID.String())

	c.JSON(http.StatusOK, gin.H{"data": product})
}

// VerwijderTechnischDetail godoc
// @Summary Verwijder een technisch detail van een product
// @Description Verwijdert een bestaand technisch detail van een product
// @Tags catalogus
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param input body models.VerwijderTechnischDetailCommand true "Technisch detail sleutel"
// @Success 200 {object} models.ZorgTechProduct "Successfully updated product"
// @Failure 400 {string} string "Bad Request"
// @Failure 404 {string} string "Product of detail niet gevonden"
// @Router /catalogus/product/technisch-detail [delete]
func (r *catalogusRepository) VerwijderTechnischDetail(c *gin.Context) {
	var input models.VerwijderTechnischDetailCommand
	var product models.ZorgTechProduct

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Zoek product op basis van ID
	if err := r.DB.Where("id = ?", input.ProductID).First(&product).Error(); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product niet gevonden"})
		return
	}

	// Verwijder het detail
	found := false
	for i, detail := range product.TechnischeDetails {
		if detail.Sleutel == input.Sleutel {
			product.TechnischeDetails = append(product.TechnischeDetails[:i], product.TechnischeDetails[i+1:]...)
			found = true
			break
		}
	}

	if !found {
		c.JSON(http.StatusNotFound, gin.H{"error": "Technisch detail niet gevonden"})
		return
	}

	// Sla update op
	if err := r.DB.Save(&product).Error(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Kon technisch detail niet verwijderen"})
		return
	}

	// Registreer domain event
	event := models.ProductEvent{
		EventID:     uuid.New(),
		ProductID:   product.ID,
		Type:        "TechnischDetailVerwijderd",
		Payload:     input.Sleutel,
		TriggeredBy: "Catalogusbeheerder",
	}
	if err := r.DB.Create(&event).Error(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Kon event niet registreren"})
		return
	}

	// Invalideer cache
	r.RedisClient.Del(*r.Ctx, "catalogus_product_"+product.ID.String())

	c.JSON(http.StatusOK, gin.H{"data": product})
}

// GetProductById godoc
// @Summary Haal een product op basis van ID
// @Description Geeft het volledige product terug voor het opgegeven ID
// @Tags catalogus
// @Security ApiKeyAuth
// @Produce json
// @Param zorgtechId path string true "ZorgTech Product ID"
// @Success 200 {object} models.ZorgTechProduct "Successfully retrieved product"
// @Failure 404 {string} string "Product niet gevonden"
// @Router /catalogus/product/{zorgtechId} [get]
func (r *catalogusRepository) GetProductById(c *gin.Context) {
	var product models.ZorgTechProduct
	zorgtechID := c.Param("zorgtechId")

	// Probeer eerst cache
	cacheKey := "catalogus_product_" + zorgtechID
	cachedProduct, err := r.RedisClient.Get(*r.Ctx, cacheKey).Result()
	if err == nil {
		var cached models.ZorgTechProduct
		if err := json.Unmarshal([]byte(cachedProduct), &cached); err == nil {
			c.JSON(http.StatusOK, gin.H{"data": cached})
			return
		}
	}

	// Zoek product in database
	if err := r.DB.Where("id = ? AND is_actief = ?", zorgtechID, true).First(&product).Error(); err != nil {
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

// FindByCategorie godoc
// @Summary Zoek producten op categorie
// @Description Geeft alle actieve producten terug voor de opgegeven categorie
// @Tags catalogus
// @Security ApiKeyAuth
// @Produce json
// @Param categorie path string true "Categorie naam"
// @Success 200 {array} models.ZorgTechProduct "Successfully retrieved products"
// @Failure 404 {string} string "Geen producten gevonden"
// @Router /catalogus/categorie/{categorie} [get]
func (r *catalogusRepository) FindByCategorie(c *gin.Context) {
	var producten []models.ZorgTechProduct
	categorie := c.Param("categorie")

	cacheKey := "catalogus_categorie_" + categorie

	// Probeer cache met context timeout
	ctx, cancel := context.WithTimeout(*r.Ctx, 500*time.Millisecond)
	defer cancel()

	if cached, err := r.RedisClient.Get(ctx, cacheKey).Bytes(); err == nil {
		if err := json.Unmarshal(cached, &producten); err == nil {
			c.JSON(http.StatusOK, gin.H{"data": producten})
			return
		}
	}

	// Database query met error handling
	err := r.DB.Where("categorie = ? AND is_actief = ?", categorie, true).
		Find(&producten).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Interne serverfout"})
		return
	}
	if len(producten) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Geen producten gevonden voor deze categorie"})
		return
	}

	// Cache opslag met context timeout
	if serialized, err := json.Marshal(producten); err == nil {
		ctx, cancel := context.WithTimeout(*r.Ctx, 500*time.Millisecond)
		defer cancel()
		r.RedisClient.Set(ctx, cacheKey, serialized, 30*time.Minute)
	}

	c.JSON(http.StatusOK, gin.H{"data": producten})
}

// ListAlleProducten godoc
// @Summary Lijst van alle actieve producten
// @Description Geeft alle actieve producten terug in de catalogus
// @Tags catalogus
// @Security ApiKeyAuth
// @Produce json
// @Success 200 {array} models.ZorgTechProduct "Successfully retrieved products"
// @Failure 404 {string} string "Geen producten gevonden"
// @Router /catalogus/producten [get]
func (r *catalogusRepository) ListAlleProducten(c *gin.Context) {
	var producten []models.ZorgTechProduct
	cacheKey := "catalogus_alle_producten"

	// Paginatie parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "25"))
	offset := (page - 1) * limit

	// Cache check met context
	ctx, cancel := context.WithTimeout(*r.Ctx, 500*time.Millisecond)
	defer cancel()

	if cached, err := r.RedisClient.Get(ctx, cacheKey).Bytes(); err == nil {
		if err := json.Unmarshal(cached, &producten); err == nil {
			c.JSON(http.StatusOK, gin.H{
				"data":  producten,
				"page":  page,
				"limit": limit,
				"total": len(producten),
			})
			return
		}
	}

	// Database query met paginatie
	query := r.DB.Where("is_actief = ?", true)
	var total int64
	query.Model(&models.ZorgTechProduct{}).Count(&total)

	if err := query.Offset(offset).Limit(limit).Find(&producten).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Kon producten niet ophalen"})
		return
	}

	// Cache opslag
	if serialized, err := json.Marshal(producten); err == nil {
		ctx, cancel := context.WithTimeout(*r.Ctx, 500*time.Millisecond)
		defer cancel()
		r.RedisClient.Set(ctx, cacheKey, serialized, 30*time.Minute)
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  producten,
		"page":  page,
		"limit": limit,
		"total": total,
	})
}

// ZoekOpNaam godoc
// @Summary Zoek producten op naam of beschrijving
// @Description Geeft producten terug waar de naam of beschrijving overeenkomt met de zoekterm
// @Tags catalogus
// @Security ApiKeyAuth
// @Produce json
// @Param query query string true "Zoekterm"
// @Success 200 {array} models.ZorgTechProduct "Successfully retrieved products"
// @Failure 404 {string} string "Geen producten gevonden"
// @Router /catalogus/zoek [get]
func (r *catalogusRepository) ZoekOpNaam(c *gin.Context) {
	var producten []models.ZorgTechProduct
	zoekterm := strings.TrimSpace(c.Query("query"))

	if zoekterm == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Zoekterm is verplicht"})
		return
	}

	cacheKey := "catalogus_zoek_" + zoekterm

	// Cache check met context
	ctx, cancel := context.WithTimeout(*r.Ctx, 500*time.Millisecond)
	defer cancel()

	if cached, err := r.RedisClient.Get(ctx, cacheKey).Bytes(); err == nil {
		if err := json.Unmarshal(cached, &producten); err == nil {
			c.JSON(http.StatusOK, gin.H{"data": producten})
			return
		}
	}

	// Database query met full-text search indien beschikbaar
	zoekterm = "%" + zoekterm + "%"
	query := r.DB.Where(
		"(naam ILIKE ? OR beschrijving ILIKE ?) AND is_actief = ?",
		zoekterm, zoekterm, true,
	).Order("naam ASC")

	if err := query.Find(&producten).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Zoekfout"})
		return
	}

	// Cache opslag met kortere TTL voor zoekresultaten
	if serialized, err := json.Marshal(producten); err == nil {
		ctx, cancel := context.WithTimeout(*r.Ctx, 500*time.Millisecond)
		defer cancel()
		r.RedisClient.Set(ctx, cacheKey, serialized, 5*time.Minute) // Kortere TTL voor zoekresultaten
	}

	c.JSON(http.StatusOK, gin.H{"data": producten})
}
