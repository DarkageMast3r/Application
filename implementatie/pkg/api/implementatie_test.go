package api

import (
	"ZorgTechImplementatie/pkg/cache"
	"ZorgTechImplementatie/pkg/database"
	"ZorgTechImplementatie/pkg/models"
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestNewImplementatieRepository(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := database.NewMockDatabase(ctrl)
	mockCache := cache.NewMockCache(ctrl)
	mockCtx := context.Background()

	repo := NewImplementatieRepository(mockDB, mockCache, &mockCtx)

	assert.NotNil(t, repo, "NewImplementatieRepository should return a non-nil instance")
	assert.Equal(t, mockDB, repo.DB, "DB should be set to the mock database instance")
	assert.Equal(t, mockCache, repo.RedisClient, "RedisClient should be set to the mock cache instance")
}

func TestHealthcheck(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	_, router := gin.CreateTestContext(recorder)

	mockRepo := NewMockImplementatieRepository(ctrl)
	mockRepo.EXPECT().Healthcheck(gomock.Any()).Do(func(c *gin.Context) {
		c.JSON(http.StatusOK, "ok")
	})

	router.GET("/healthcheck", mockRepo.Healthcheck)

	req, _ := http.NewRequest(http.MethodGet, "/healthcheck", nil)
	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Equal(t, "\"ok\"", recorder.Body.String())
}

func TestAanvraagProduct(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := database.NewMockDatabase(ctrl)
	mockCache := cache.NewMockCache(ctrl)
	ctx := context.Background()

	repo := NewImplementatieRepository(mockDB, mockCache, &ctx)

	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.POST("/implementatie/aanvraag", repo.AanvraagProduct)

	clientID := uuid.New()
	zorgtechID := uuid.New()
	input := models.AanvraagProductCommand{
		ClientID:   clientID,
		ZorgtechID: zorgtechID,
	}
	requestBody, _ := json.Marshal(input)

	// Mock database create
	mockDB.EXPECT().Create(gomock.Any()).DoAndReturn(func(dossier *models.ImplementatieDossier) *gorm.DB {
		assert.Equal(t, clientID, dossier.ClientID)
		assert.Equal(t, zorgtechID, dossier.ZorgtechID)
		assert.Equal(t, models.StatusBesteld, dossier.Status)
		return &gorm.DB{Error: nil}
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/implementatie/aanvraag", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response struct {
		Data models.ImplementatieDossier `json:"data"`
	}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, clientID, response.Data.ClientID)
	assert.Equal(t, zorgtechID, response.Data.ZorgtechID)
}

func TestOntvangProduct(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := database.NewMockDatabase(ctrl)
	mockCache := cache.NewMockCache(ctrl)
	ctx := context.Background()

	repo := NewImplementatieRepository(mockDB, mockCache, &ctx)

	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.POST("/implementatie/ontvang", repo.OntvangProduct)

	clientID := uuid.New()
	zorgtechID := uuid.New()
	serienummer := "SN12345"
	input := models.OntvangProductCommand{
		ClientID:    clientID,
		ZorgtechID:  zorgtechID,
		Serienummer: serienummer,
	}
	requestBody, _ := json.Marshal(input)

	// Mock database find
	mockDB.EXPECT().Where("client_id = ? AND zorgtech_id = ?", clientID, zorgtechID).
		Return(mockDB).Times(1)
	mockDB.EXPECT().First(gomock.Any()).
		DoAndReturn(func(dossier *models.ImplementatieDossier) *gorm.DB {
			dossier.ImplementatieID = uuid.New()
			dossier.ClientID = clientID
			dossier.ZorgtechID = zorgtechID
			dossier.Status = models.StatusBesteld
			return &gorm.DB{Error: nil}
		}).Times(1)

	// Mock database save
	mockDB.EXPECT().Save(gomock.Any()).Return(&gorm.DB{Error: nil}).Times(1)

	// Mock cache invalidation
	mockCache.EXPECT().Del(ctx, "implementatie_status_"+clientID.String()).Return(redis.NewIntResult(1, nil))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/implementatie/ontvang", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response struct {
		Data models.ImplementatieDossier `json:"data"`
	}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, models.StatusGeleverd, response.Data.Status)
	assert.Equal(t, serienummer, response.Data.Serienummer)
}

func TestInstalleerProduct(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := database.NewMockDatabase(ctrl)
	mockCache := cache.NewMockCache(ctrl)
	ctx := context.Background()

	repo := NewImplementatieRepository(mockDB, mockCache, &ctx)

	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.POST("/implementatie/installeer", repo.InstalleerProduct)

	clientID := uuid.New()
	serienummer := "SN12345"
	input := models.InstalleerProductCommand{
		ClientID:    clientID,
		Serienummer: serienummer,
	}
	requestBody, _ := json.Marshal(input)

	// Mock database find
	mockDB.EXPECT().Where("client_id = ? AND serienummer = ?", clientID, serienummer).
		Return(mockDB).Times(1)
	mockDB.EXPECT().First(gomock.Any()).
		DoAndReturn(func(dossier *models.ImplementatieDossier) *gorm.DB {
			dossier.ImplementatieID = uuid.New()
			dossier.ClientID = clientID
			dossier.Serienummer = serienummer
			dossier.Status = models.StatusGeleverd
			return &gorm.DB{Error: nil}
		}).Times(1)

	// Mock database save
	mockDB.EXPECT().Save(gomock.Any()).Return(&gorm.DB{Error: nil}).Times(1)

	// Mock cache invalidation
	mockCache.EXPECT().Del(ctx, "implementatie_status_"+clientID.String()).Return(redis.NewIntResult(1, nil))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/implementatie/installeer", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response struct {
		Data models.ImplementatieDossier `json:"data"`
	}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, models.StatusGeinstalleerd, response.Data.Status)
	assert.NotNil(t, response.Data.InstallatieDatum)
}

func TestPersonaliseerProduct(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := database.NewMockDatabase(ctrl)
	mockCache := cache.NewMockCache(ctrl)
	ctx := context.Background()

	repo := NewImplementatieRepository(mockDB, mockCache, &ctx)

	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.POST("/implementatie/personaliseer", repo.PersonaliseerProduct)

	clientID := uuid.New()
	instellingen := models.PersonalisatieInstellingen{
		VolumeNiveau: 5,
		MeldingsType: "visueel",
		Schema:       "dag",
		Taal:         "nl",
	}
	input := models.PersonaliseerProductCommand{
		ClientID:     clientID,
		Instellingen: instellingen,
	}
	requestBody, _ := json.Marshal(input)

	// Mock database find
	mockDB.EXPECT().Where("client_id = ?", clientID).
		Return(mockDB).Times(1)
	mockDB.EXPECT().First(gomock.Any()).
		DoAndReturn(func(dossier *models.ImplementatieDossier) *gorm.DB {
			dossier.ImplementatieID = uuid.New()
			dossier.ClientID = clientID
			dossier.Status = models.StatusGeinstalleerd
			return &gorm.DB{Error: nil}
		}).Times(1)

	// Mock database save
	mockDB.EXPECT().Save(gomock.Any()).Return(&gorm.DB{Error: nil}).Times(1)

	// Mock cache invalidation
	mockCache.EXPECT().Del(ctx, "implementatie_status_"+clientID.String()).Return(redis.NewIntResult(1, nil))
	mockCache.EXPECT().Del(ctx, "persoonlijke_instellingen_"+clientID.String()).Return(redis.NewIntResult(1, nil))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/implementatie/personaliseer", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response struct {
		Data models.ImplementatieDossier `json:"data"`
	}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, models.StatusGepersonaliseerd, response.Data.Status)
	assert.Equal(t, instellingen.VolumeNiveau, response.Data.Personalisatie.VolumeNiveau)
}

func TestMarkeerAlsGeimplementeerd(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := database.NewMockDatabase(ctrl)
	mockCache := cache.NewMockCache(ctrl)
	ctx := context.Background()

	repo := NewImplementatieRepository(mockDB, mockCache, &ctx)

	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.POST("/implementatie/voltooi", repo.MarkeerAlsGeimplementeerd)

	clientID := uuid.New()
	input := models.MarkeerAlsGeimplementeerdCommand{
		ClientID: clientID,
	}
	requestBody, _ := json.Marshal(input)

	// Mock database find
	mockDB.EXPECT().Where("client_id = ?", clientID).
		Return(mockDB).Times(1)
	mockDB.EXPECT().First(gomock.Any()).
		DoAndReturn(func(dossier *models.ImplementatieDossier) *gorm.DB {
			dossier.ImplementatieID = uuid.New()
			dossier.ClientID = clientID
			dossier.Status = models.StatusGepersonaliseerd
			return &gorm.DB{Error: nil}
		}).Times(1)

	// Mock database save
	mockDB.EXPECT().Save(gomock.Any()).Return(&gorm.DB{Error: nil}).Times(1)

	// Mock cache invalidation
	mockCache.EXPECT().Del(ctx, "implementatie_status_"+clientID.String()).Return(redis.NewIntResult(1, nil))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/implementatie/voltooi", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response struct {
		Data models.ImplementatieDossier `json:"data"`
	}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, models.StatusVoltooid, response.Data.Status)
}

func TestGetProductInformatie(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := database.NewMockDatabase(ctrl)
	mockCache := cache.NewMockCache(ctrl)
	ctx := context.Background()

	repo := NewImplementatieRepository(mockDB, mockCache, &ctx)

	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.GET("/implementatie/product/:zorgtechId", repo.GetProductInformatie)

	zorgtechID := uuid.New()
	product := models.ZorgTechProduct{
		ZorgtechID: zorgtechID,
		Name:       "Test Product",
		Category:   "valpreventie",
	}

	// Mock cache get
	mockCache.EXPECT().Get(ctx, "product_info_"+zorgtechID.String()).
		Return(redis.NewStringResult("", redis.Nil)) // Cache miss

	// Mock database find
	mockDB.EXPECT().Where("zorgtech_id = ?", zorgtechID.String()).
		Return(mockDB).Times(1)
	mockDB.EXPECT().First(gomock.Any()).
		DoAndReturn(func(p *models.ZorgTechProduct) *gorm.DB {
			*p = product
			return &gorm.DB{Error: nil}
		}).Times(1)

	// Mock cache set
	serialized, _ := json.Marshal(product)
	mockCache.EXPECT().Set(ctx, "product_info_"+zorgtechID.String(), serialized, time.Hour).
		Return(redis.NewStatusResult("OK", nil))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/implementatie/product/"+zorgtechID.String(), nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response struct {
		Data models.ZorgTechProduct `json:"data"`
	}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, zorgtechID, response.Data.ZorgtechID)
	assert.Equal(t, "Test Product", response.Data.Name)
}

func TestGetImplementatieStatus(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := database.NewMockDatabase(ctrl)
	mockCache := cache.NewMockCache(ctrl)
	ctx := context.Background()

	repo := NewImplementatieRepository(mockDB, mockCache, &ctx)

	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.GET("/implementatie/status/:clientId", repo.GetImplementatieStatus)

	clientID := uuid.New()
	dossier := models.ImplementatieDossier{
		ImplementatieID: uuid.New(),
		ClientID:        clientID,
		Status:          models.StatusGeinstalleerd,
	}

	// Mock cache get
	mockCache.EXPECT().Get(ctx, "implementatie_status_"+clientID.String()).
		Return(redis.NewStringResult("", redis.Nil)) // Cache miss

	// Mock database find
	mockDB.EXPECT().Where("client_id = ?", clientID.String()).
		Return(mockDB).Times(1)
	mockDB.EXPECT().First(gomock.Any()).
		DoAndReturn(func(d *models.ImplementatieDossier) *gorm.DB {
			*d = dossier
			return &gorm.DB{Error: nil}
		}).Times(1)

	// Mock cache set
	serialized, _ := json.Marshal(dossier)
	mockCache.EXPECT().Set(ctx, "implementatie_status_"+clientID.String(), serialized, 30*time.Minute).
		Return(redis.NewStatusResult("OK", nil))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/implementatie/status/"+clientID.String(), nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response struct {
		Data models.ImplementatieDossier `json:"data"`
	}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, clientID, response.Data.ClientID)
	assert.Equal(t, models.StatusGeinstalleerd, response.Data.Status)
}
