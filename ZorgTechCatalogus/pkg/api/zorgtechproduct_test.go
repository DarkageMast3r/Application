package api

import (
	"ZorgTechCatalogus/pkg/cache"
	"ZorgTechCatalogus/pkg/database"
	models "ZorgTechCatalogus/pkg/models"
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewCatalogusRepository(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := database.NewMockDatabase(ctrl)
	mockCache := cache.NewMockCache(ctrl)
	mockCtx := context.Background()

	repo := NewCatalogusRepository(mockDB, mockCache, &mockCtx)

	assert.NotNil(t, repo, "NewCatalogusRepository should return a non-nil instance")
	assert.Equal(t, mockDB, repo.DB, "DB should be set to the mock database instance")
	assert.Equal(t, mockCache, repo.RedisClient, "RedisClient should be set to the mock cache instance")
	assert.Equal(t, &mockCtx, repo.Ctx, "Ctx should be set to the mock context")
}

func TestHealthcheck(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	_, router := gin.CreateTestContext(recorder)

	mockRepo := NewMockCatalogusRepository(ctrl)
	mockRepo.EXPECT().Healthcheck(gomock.Any()).Do(func(c *gin.Context) {
		c.JSON(http.StatusOK, "ok")
	})

	router.GET("/healthcheck", mockRepo.Healthcheck)

	req, _ := http.NewRequest(http.MethodGet, "/healthcheck", nil)
	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Equal(t, "\"ok\"", recorder.Body.String())
}

func TestMaakZorgTechProduct(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := database.NewMockDatabase(ctrl)
	mockCache := cache.NewMockCache(ctrl)
	ctx := context.Background()

	repo := NewCatalogusRepository(mockDB, mockCache, &ctx)

	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.POST("/catalogus/product", repo.MaakZorgTechProduct)

	input := models.MaakZorgTechProductCommand{
		Naam:         "Test Product",
		Beschrijving: "Test Beschrijving",
		Categorie:    "valpreventie",
		Prijs:        99.99,
		Leverancier:  "Test Leverancier",
	}
	requestBody, _ := json.Marshal(input)

	// Mock database create voor product
	mockDB.EXPECT().Create(gomock.Any()).DoAndReturn(func(product *models.ZorgTechProduct) database.Database {
		assert.Equal(t, input.Naam, product.Naam)
		assert.Equal(t, input.Beschrijving, product.Beschrijving)
		assert.Equal(t, input.Categorie, product.Categorie)
		assert.Equal(t, input.Prijs, product.Prijs)
		assert.Equal(t, input.Leverancier, product.Leverancier)
		assert.True(t, product.IsActief)
		return mockDB
	}).Times(1)

	// Voeg deze toe: mock voor Error() na eerste Create
	mockDB.EXPECT().Error().Return(nil).Times(1)

	// Mock database create voor event
	mockDB.EXPECT().Create(gomock.Any()).DoAndReturn(func(event *models.ProductEvent) database.Database {
		assert.Equal(t, "ZorgTechProductAangemaakt", event.Type)
		assert.Equal(t, input.Naam, event.Payload)
		return mockDB
	}).Times(1)

	mockDB.EXPECT().Error().Return(nil).Times(1)

	// Mock cache invalidation
	mockCache.EXPECT().Del(ctx, "catalogus_alle_producten").Return(redis.NewIntResult(1, nil))
	mockCache.EXPECT().Del(ctx, "catalogus_categorie_"+input.Categorie).Return(redis.NewIntResult(1, nil))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/catalogus/product", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response struct {
		Data models.ZorgTechProduct `json:"data"`
	}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, input.Naam, response.Data.Naam)
	assert.Equal(t, input.Categorie, response.Data.Categorie)

}

// Mock voor de CatalogusRepository interface
type MockCatalogusRepository struct {
	ctrl     *gomock.Controller
	recorder *MockCatalogusRepositoryMockRecorder
}

func NewMockCatalogusRepository(ctrl *gomock.Controller) *MockCatalogusRepository {
	mock := &MockCatalogusRepository{ctrl: ctrl}
	mock.recorder = &MockCatalogusRepositoryMockRecorder{mock}
	return mock
}

func (m *MockCatalogusRepository) EXPECT() *MockCatalogusRepositoryMockRecorder {
	return m.recorder
}

func (m *MockCatalogusRepository) Healthcheck(c *gin.Context) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Healthcheck", c)
}

func (m *MockCatalogusRepository) MaakZorgTechProduct(c *gin.Context) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "MaakZorgTechProduct", c)
}

func (m *MockCatalogusRepository) WijzigZorgTechProduct(c *gin.Context) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "WijzigZorgTechProduct", c)
}

func (m *MockCatalogusRepository) VerwijderZorgTechProduct(c *gin.Context) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "VerwijderZorgTechProduct", c)
}

func (m *MockCatalogusRepository) VoegTechnischDetailToe(c *gin.Context) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "VoegTechnischDetailToe", c)
}

func (m *MockCatalogusRepository) VerwijderTechnischDetail(c *gin.Context) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "VerwijderTechnischDetail", c)
}

func (m *MockCatalogusRepository) GetProductById(c *gin.Context) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "GetProductById", c)
}

func (m *MockCatalogusRepository) FindByCategorie(c *gin.Context) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "FindByCategorie", c)
}

func (m *MockCatalogusRepository) ListAlleProducten(c *gin.Context) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "ListAlleProducten", c)
}

func (m *MockCatalogusRepository) ZoekOpNaam(c *gin.Context) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "ZoekOpNaam", c)
}

type MockCatalogusRepositoryMockRecorder struct {
	mock *MockCatalogusRepository
}

func (mr *MockCatalogusRepositoryMockRecorder) Healthcheck(c interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Healthcheck", reflect.TypeOf((*catalogusRepository)(nil).Healthcheck), c)
}

func (mr *MockCatalogusRepositoryMockRecorder) MaakZorgTechProduct(c interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MaakZorgTechProduct", reflect.TypeOf((*catalogusRepository)(nil).MaakZorgTechProduct), c)
}

func (mr *MockCatalogusRepositoryMockRecorder) WijzigZorgTechProduct(c interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WijzigZorgTechProduct", reflect.TypeOf((*catalogusRepository)(nil).WijzigZorgTechProduct), c)
}

func (mr *MockCatalogusRepositoryMockRecorder) VerwijderZorgTechProduct(c interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "VerwijderZorgTechProduct", reflect.TypeOf((*catalogusRepository)(nil).VerwijderZorgTechProduct), c)
}

func (mr *MockCatalogusRepositoryMockRecorder) VoegTechnischDetailToe(c interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "VoegTechnischDetailToe", reflect.TypeOf((*catalogusRepository)(nil).VoegTechnischDetailToe), c)
}

func (mr *MockCatalogusRepositoryMockRecorder) VerwijderTechnischDetail(c interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "VerwijderTechnischDetail", reflect.TypeOf((*catalogusRepository)(nil).VerwijderTechnischDetail), c)
}

func (mr *MockCatalogusRepositoryMockRecorder) GetProductById(c interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetProductById", reflect.TypeOf((*catalogusRepository)(nil).GetProductById), c)
}

func (mr *MockCatalogusRepositoryMockRecorder) FindByCategorie(c interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByCategorie", reflect.TypeOf((*catalogusRepository)(nil).FindByCategorie), c)
}

func (mr *MockCatalogusRepositoryMockRecorder) ListAlleProducten(c interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListAlleProducten", reflect.TypeOf((*catalogusRepository)(nil).ListAlleProducten), c)
}

func (mr *MockCatalogusRepositoryMockRecorder) ZoekOpNaam(c interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ZoekOpNaam", reflect.TypeOf((*catalogusRepository)(nil).ZoekOpNaam), c)
}

func TestGetProductById(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := database.NewMockDatabase(ctrl)
	mockCache := cache.NewMockCache(ctrl)
	ctx := context.Background()

	repo := NewCatalogusRepository(mockDB, mockCache, &ctx)

	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.GET("/catalogus/product/:zorgtechId", repo.GetProductById)

	productID := uuid.New()
	product := models.ZorgTechProduct{
		ID:       productID,
		Naam:     "Test Product",
		IsActief: true,
	}

	// Mock cache miss
	mockCache.EXPECT().Get(gomock.Any(), "catalogus_product_"+productID.String()).
		Return(redis.NewStringResult("", redis.Nil))

	// Mock database query
	mockDB.EXPECT().Where("id = ? AND is_actief = ?", productID.String(), true).
		Return(mockDB)
	mockDB.EXPECT().First(gomock.Any()).
		DoAndReturn(func(dest interface{}, conds ...interface{}) database.Database {
			if p, ok := dest.(*models.ZorgTechProduct); ok {
				*p = product
			}
			return mockDB
		})
	mockDB.EXPECT().Error().Return(nil)

	// Mock cache set
	serialized, _ := json.Marshal(product)
	mockCache.EXPECT().Set(gomock.Any(), "catalogus_product_"+productID.String(), serialized, time.Hour).
		Return(redis.NewStatusResult("OK", nil))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/catalogus/product/"+productID.String(), nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response struct {
		Data models.ZorgTechProduct `json:"data"`
	}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, productID, response.Data.ID)
	assert.Equal(t, "Test Product", response.Data.Naam)
}
func TestFindByCategorie(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := database.NewMockDatabase(ctrl)
	mockCache := cache.NewMockCache(ctrl)
	ctx := context.Background()

	repo := NewCatalogusRepository(mockDB, mockCache, &ctx)

	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.GET("/catalogus/categorie/:categorie", repo.FindByCategorie)

	categorie := "valpreventie"
	producten := []models.ZorgTechProduct{
		{
			ID:        uuid.New(),
			Naam:      "Product 1",
			Categorie: categorie,
			IsActief:  true,
		},
	}

	// Mock cache get - return cache miss
	mockCache.EXPECT().Get(gomock.Any(), "catalogus_categorie_"+categorie).
		Return(redis.NewStringResult("", redis.Nil)).Times(1)

	// Mock database chain
	mockDB.EXPECT().Where("categorie = ? AND is_actief = ?", categorie, true).
		Return(mockDB).Times(1)

	// Belangrijk: gebruik gomock.AssignableToTypeOf voor de Find parameter
	mockDB.EXPECT().Find(gomock.AssignableToTypeOf(&[]models.ZorgTechProduct{})).
		DoAndReturn(func(dest interface{}) database.Database {
			if products, ok := dest.(*[]models.ZorgTechProduct); ok {
				*products = producten
			}
			return mockDB
		}).Times(1)

	// Mock de Error() check - return nil voor succes
	mockDB.EXPECT().Error().
		Return(nil).Times(1)

	// Mock cache set
	serialized, _ := json.Marshal(producten)
	mockCache.EXPECT().Set(gomock.Any(), "catalogus_categorie_"+categorie, serialized, 30*time.Minute).
		Return(redis.NewStatusResult("OK", nil)).Times(1)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/catalogus/categorie/"+categorie, nil)
	r.ServeHTTP(w, req)

	// Debug output
	t.Logf("Response status: %d", w.Code)
	t.Logf("Response body: %s", w.Body.String())

	assert.Equal(t, http.StatusOK, w.Code, "Expected status code 200")

	var response struct {
		Data []models.ZorgTechProduct `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}
	assert.Len(t, response.Data, 1, "Expected 1 product in response")
	assert.Equal(t, categorie, response.Data[0].Categorie, "Category should match")
}

func TestWijzigZorgTechProduct(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := database.NewMockDatabase(ctrl)
	mockCache := cache.NewMockCache(ctrl)
	ctx := context.Background()

	repo := NewCatalogusRepository(mockDB, mockCache, &ctx)

	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.PUT("/catalogus/product", repo.WijzigZorgTechProduct)

	productID := uuid.New()
	newNaam := "Gewijzigde Naam"
	input := models.WijzigZorgTechProductCommand{
		ProductID: productID,
		Naam:      &newNaam,
	}
	requestBody, _ := json.Marshal(input)

	// Mock database find
	mockDB.EXPECT().Where("id = ?", productID).
		Return(mockDB)
	mockDB.EXPECT().First(gomock.Any()).
		DoAndReturn(func(dest interface{}, conds ...interface{}) database.Database {
			if p, ok := dest.(*models.ZorgTechProduct); ok {
				*p = models.ZorgTechProduct{
					ID:       productID,
					Naam:     "Oude Naam",
					IsActief: true,
				}
			}
			return mockDB
		})
	mockDB.EXPECT().Error().Return(nil)

	// Mock database save
	mockDB.EXPECT().Save(gomock.Any()).
		Return(mockDB)
	mockDB.EXPECT().Error().Return(nil)

	// Mock event creation
	mockDB.EXPECT().Create(gomock.Any()).
		Return(mockDB)
	mockDB.EXPECT().Error().Return(nil)

	// Mock cache invalidation
	mockCache.EXPECT().Del(ctx, "catalogus_product_"+productID.String()).Return(redis.NewIntResult(1, nil))
	mockCache.EXPECT().Del(ctx, "catalogus_alle_producten").Return(redis.NewIntResult(1, nil))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/catalogus/product", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response struct {
		Data models.ZorgTechProduct `json:"data"`
	}
	json.Unmarshal(w.Body.Bytes(), &response)
}

func TestVerwijderZorgTechProduct(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := database.NewMockDatabase(ctrl)
	mockCache := cache.NewMockCache(ctrl)
	ctx := context.Background()

	repo := NewCatalogusRepository(mockDB, mockCache, &ctx)

	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.DELETE("/catalogus/product", repo.VerwijderZorgTechProduct)

	productID := uuid.New()
	input := models.VerwijderZorgTechProductCommand{
		ProductID: productID,
	}
	requestBody, _ := json.Marshal(input)

	// Mock database find
	mockDB.EXPECT().Where("id = ?", productID).
		Return(mockDB)
	mockDB.EXPECT().First(gomock.Any()).
		DoAndReturn(func(dest interface{}, conds ...interface{}) database.Database {
			if p, ok := dest.(*models.ZorgTechProduct); ok {
				*p = models.ZorgTechProduct{
					ID:        productID,
					Naam:      "Test Product",
					Categorie: "valpreventie",
					IsActief:  true,
				}
			}
			return mockDB
		})
	mockDB.EXPECT().Error().Return(nil)

	// Mock database save
	mockDB.EXPECT().Save(gomock.Any()).
		Return(mockDB)
	mockDB.EXPECT().Error().Return(nil)

	// Mock event creation
	mockDB.EXPECT().Create(gomock.Any()).
		Return(mockDB)
	mockDB.EXPECT().Error().Return(nil)

	// Mock cache invalidation
	mockCache.EXPECT().Del(ctx, "catalogus_product_"+productID.String()).Return(redis.NewIntResult(1, nil))
	mockCache.EXPECT().Del(ctx, "catalogus_alle_producten").Return(redis.NewIntResult(1, nil))
	mockCache.EXPECT().Del(ctx, "catalogus_categorie_valpreventie").Return(redis.NewIntResult(1, nil))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/catalogus/product", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response struct {
		Data models.ZorgTechProduct `json:"data"`
	}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.False(t, response.Data.IsActief)
}
