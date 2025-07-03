package api

import (
	"ZorgTechCatalogus/pkg/cache"
	"ZorgTechCatalogus/pkg/database"
	zorgtechproduct "ZorgTechCatalogus/pkg/models"
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type ZorgTechProductRepository interface {
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

type zorgTechProductRepository struct {
	DB          database.Database
	RedisClient cache.Cache
	Ctx         *context.Context
}

func NewZorgTechProductRepository(db database.Database, redisClient cache.Cache, ctx *context.Context) *zorgTechProductRepository {
	return &zorgTechProductRepository{
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
func (r *zorgTechProductRepository) Healthcheck(c *gin.Context) {
	c.JSON(http.StatusOK, "ok")
}

// MaakZorgTechProduct godoc
// @Summary Maak een nieuw zorgtech product aan
// @Description Registreert een nieuw zorgtechnologie product in de catalogus
// @Tags catalogus
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param input body MaakZorgTechProductCommand true "Product gegevens"
// @Success 201 {object} ZorgTechProduct "Successfully created product"
// @Failure 400 {string} string "Bad Request"
// @Router /catalogus/product [post]
func (r *zorgTechProductRepository) MaakZorgTechProduct(c *gin.Context) {
	var input zorgtechproduct.MaakZorgTechProductCommand

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product, err := zorgtechproduct.NewZorgTechProduct(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Kon product niet aanmaken"})
		return
	}

	if err := r.DB.Create(product).Error(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Kon product niet opslaan"})
		return
	}

	// Invalideer cache voor productlijsten
	r.RedisClient.Del(*r.Ctx, "alle_producten")
	r.RedisClient.Del(*r.Ctx, "categorie_"+product.Categorie)

	c.JSON(http.StatusCreated, gin.H{"data": product})
}

// WijzigZorgTechProduct godoc
// @Summary Wijzig een bestaand zorgtech product
// @Description Past de eigenschappen van een bestaand product aan
// @Tags catalogus
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param input body WijzigZorgTechProductCommand true "Wijzigingsgegevens"
// @Success 200 {object} ZorgTechProduct "Successfully updated product"
// @Failure 400 {string} string "Bad Request"
// @Failure 404 {string} string "Product niet gevonden"
// @Router /catalogus/product [put]
func (r *zorgTechProductRepository) WijzigZorgTechProduct(c *gin.Context) {
	var input zorgtechproduct.WijzigZorgTechProductCommand
	var product zorgtechproduct.ZorgTechProduct

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := r.DB.Where("id = ?", input.ProductID).First(&product).Error(); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product niet gevonden"})
		return
	}

	if err := product.Wijzig(input); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Kon product niet wijzigen"})
		return
	}

	if err := r.DB.Save(&product).Error(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Kon wijzigingen niet opslaan"})
		return
	}

	// Invalideer relevante cache
	r.RedisClient.Del(*r.Ctx, "product_"+product.ID.String())
	r.RedisClient.Del(*r.Ctx, "alle_producten")
	r.RedisClient.Del(*r.Ctx, "categorie_"+product.Categorie)

	c.JSON(http.StatusOK, gin.H{"data": product})
}

// VerwijderZorgTechProduct godoc
// @Summary Verwijder een zorgtech product
// @Description Archiveert een product zodat het niet meer actief is in de catalogus
// @Tags catalogus
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param input body VerwijderZorgTechProductCommand true "Product ID"
// @Success 200 {object} ZorgTechProduct "Successfully deleted product"
// @Failure 400 {string} string "Bad Request"
// @Failure 404 {string} string "Product niet gevonden"
// @Router /catalogus/product [delete]
func (r *zorgTechProductRepository) VerwijderZorgTechProduct(c *gin.Context) {
	var input zorgtechproduct.VerwijderZorgTechProductCommand
	var product zorgtechproduct.ZorgTechProduct

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := r.DB.Where("id = ?", input.ProductID).First(&product).Error(); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product niet gevonden"})
		return
	}

	product.Verwijder()

	if err := r.DB.Save(&product).Error(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Kon product niet archiveren"})
		return
	}

	// Invalideer relevante cache
	r.RedisClient.Del(*r.Ctx, "product_"+product.ID.String())
	r.RedisClient.Del(*r.Ctx, "alle_producten")
	r.RedisClient.Del(*r.Ctx, "categorie_"+product.Categorie)

	c.JSON(http.StatusOK, gin.H{"data": product})
}

// VoegTechnischDetailToe godoc
// @Summary Voeg een technisch detail toe aan een product
// @Description Voegt een nieuw technisch detail toe of werkt een bestaand detail bij
// @Tags catalogus
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param input body VoegTechnischDetailToeCommand true "Detail gegevens"
// @Success 200 {object} ZorgTechProduct "Successfully updated product"
// @Failure 400 {string} string "Bad Request"
// @Failure 404 {string} string "Product niet gevonden"
// @Router /catalogus/product/technisch-detail [post]
func (r *zorgTechProductRepository) VoegTechnischDetailToe(c *gin.Context) {
	var input zorgtechproduct.VoegTechnischDetailToeCommand
	var product zorgtechproduct.ZorgTechProduct

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := r.DB.Where("id = ?", input.ProductID).First(&product).Error(); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product niet gevonden"})
		return
	}

	product.VoegTechnischDetailToe(input)

	if err := r.DB.Save(&product).Error(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Kon technisch detail niet toevoegen"})
		return
	}

	// Invalideer cache voor dit product
	r.RedisClient.Del(*r.Ctx, "product_"+product.ID.String())

	c.JSON(http.StatusOK, gin.H{"data": product})
}

// VerwijderTechnischDetail godoc
// @Summary Verwijder een technisch detail van een product
// @Description Verwijdert een specifiek technisch detail van een product
// @Tags catalogus
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param input body VerwijderTechnischDetailCommand true "Detail gegevens"
// @Success 200 {object} ZorgTechProduct "Successfully updated product"
// @Failure 400 {string} string "Bad Request"
// @Failure 404 {string} string "Product of detail niet gevonden"
// @Router /catalogus/product/technisch-detail [delete]
func (r *zorgTechProductRepository) VerwijderTechnischDetail(c *gin.Context) {
	var input zorgtechproduct.VerwijderTechnischDetailCommand
	var product zorgtechproduct.ZorgTechProduct

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := r.DB.Where("id = ?", input.ProductID).First(&product).Error(); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product niet gevonden"})
		return
	}

	if !product.VerwijderTechnischDetail(input) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Technisch detail niet gevonden"})
		return
	}

	if err := r.DB.Save(&product).Error(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Kon technisch detail niet verwijderen"})
		return
	}

	// Invalideer cache voor dit product
	r.RedisClient.Del(*r.Ctx, "product_"+product.ID.String())

	c.JSON(http.StatusOK, gin.H{"data": product})
}

// GetProductById godoc
// @Summary Haal een product op basis van ID op
// @Description Geeft alle details van een specifiek zorgtech product
// @Tags catalogus
// @Security ApiKeyAuth
// @Produce json
// @Param zorgtechId path string true "ZorgTech Product ID"
// @Success 200 {object} ZorgTechProduct "Successfully retrieved product"
// @Failure 404 {string} string "Product niet gevonden"
// @Router /catalogus/product/{zorgtechId} [get]
func (r *zorgTechProductRepository) GetProductById(c *gin.Context) {
	var product zorgtechproduct.ZorgTechProduct
	zorgtechID := c.Param("zorgtechId")

	// Probeer eerst cache
	cacheKey := "product_" + zorgtechID
	cachedProduct, err := r.RedisClient.Get(*r.Ctx, cacheKey).Result()
	if err == nil {
		var cached zorgtechproduct.ZorgTechProduct
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
// @Description Geeft een lijst van alle actieve producten in een specifieke categorie
// @Tags catalogus
// @Security ApiKeyAuth
// @Produce json
// @Param categorie path string true "Categorie naam"
// @Success 200 {array} ZorgTechProduct "Successfully retrieved products"
// @Router /catalogus/categorie/{categorie} [get]
func (r *zorgTechProductRepository) FindByCategorie(c *gin.Context) {
	var producten []zorgtechproduct.ZorgTechProduct
	categorie := c.Param("categorie")

	// Probeer eerst cache
	cacheKey := "categorie_" + categorie
	cachedProducts, err := r.RedisClient.Get(*r.Ctx, cacheKey).Result()
	if err == nil {
		var cached []zorgtechproduct.ZorgTechProduct
		if err := json.Unmarshal([]byte(cachedProducts), &cached); err == nil {
			c.JSON(http.StatusOK, gin.H{"data": cached})
			return
		}
	}

	// Zoek producten in database
	if err := r.DB.Where("categorie = ? AND is_actief = ?", categorie, true).Find(&producten).Error(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Kon producten niet ophalen"})
		return
	}

	// Sla in cache op
	serialized, err := json.Marshal(producten)
	if err == nil {
		r.RedisClient.Set(*r.Ctx, cacheKey, serialized, 30*time.Minute)
	}

	c.JSON(http.StatusOK, gin.H{"data": producten})
}

// ListAlleProducten godoc
// @Summary Lijst van alle actieve producten
// @Description Geeft een lijst van alle actieve zorgtech producten in de catalogus
// @Tags catalogus
// @Security ApiKeyAuth
// @Produce json
// @Success 200 {array} ZorgTechProduct "Successfully retrieved products"
// @Router /catalogus/producten [get]
func (r *zorgTechProductRepository) ListAlleProducten(c *gin.Context) {
	var producten []zorgtechproduct.ZorgTechProduct

	// Probeer eerst cache
	cacheKey := "alle_producten"
	cachedProducts, err := r.RedisClient.Get(*r.Ctx, cacheKey).Result()
	if err == nil {
		var cached []zorgtechproduct.ZorgTechProduct
		if err := json.Unmarshal([]byte(cachedProducts), &cached); err == nil {
			c.JSON(http.StatusOK, gin.H{"data": cached})
			return
		}
	}

	// Zoek producten in database
	if err := r.DB.Where("is_actief = ?", true).Find(&producten).Error(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Kon producten niet ophalen"})
		return
	}

	// Sla in cache op
	serialized, err := json.Marshal(producten)
	if err == nil {
		r.RedisClient.Set(*r.Ctx, cacheKey, serialized, 30*time.Minute)
	}

	c.JSON(http.StatusOK, gin.H{"data": producten})
}

// ZoekOpNaam godoc
// @Summary Zoek producten op naam of beschrijving
// @Description Geeft een lijst van producten waarvan de naam of beschrijving overeenkomt met de zoekterm
// @Tags catalogus
// @Security ApiKeyAuth
// @Produce json
// @Param zoekterm query string true "Zoekterm"
// @Success 200 {array} ZorgTechProduct "Successfully retrieved products"
// @Router /catalogus/zoek [get]
func (r *zorgTechProductRepository) ZoekOpNaam(c *gin.Context) {
	var producten []zorgtechproduct.ZorgTechProduct
	zoekterm := c.Query("zoekterm")

	// Zoek in database (full-text search zou beter zijn)
	if err := r.DB.Where("is_actief = ? AND (naam ILIKE ? OR beschrijving ILIKE ?)",
		true, "%"+zoekterm+"%", "%"+zoekterm+"%").Find(&producten).Error(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Kon zoekopdracht niet uitvoeren"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": producten})
}
