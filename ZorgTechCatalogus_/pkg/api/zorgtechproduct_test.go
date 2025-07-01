package zorgtechproduct

import (
	"bytes"
	"context"
	"encoding/json"
	"ZorgTechCatalogus/pkg/cache"
	"ZorgTechCatalogus/pkg/database"
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

func TestNewZorgTechProductRepository(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := database.NewMockDatabase(ctrl)
	mockCache := cache.NewMockCache(ctrl)
	mockCtx := context.Background()

	repo := NewZorgTechProductRepository(mockDB, mockCache, &mockCtx)

	assert.NotNil(t, repo, "NewZorgTechProductRepository should return a non-nil instance")
	assert.Equal(t, mockDB, repo.DB, "DB should be set to the mock database instance")
	assert.Equal(t, mockCache, repo.RedisClient, "RedisClient should be set to the mock cache instance")
}

func TestHealthcheck(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	_, router := gin.CreateTestContext(recorder)

	mockRepo := NewMockZorgTechProductRepository(ctrl)
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

	repo := NewZorgTechProductRepository(mockDB, mockCache, &ctx)

	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.POST("/catalogus/product", repo.MaakZorgTechProduct)

	input := MaakZorgTechProductCommand{
		Naam:         "Test Product",
		Beschrijving: "Test Beschrijving",
		Categorie:    "valpreventie",
		Prijs:        99.99,
		Leverancier:  "Test Leverancier",
	}
	requestBody, _ := json.Marshal(input)

	// Mock database create
	mockDB.EXPECT().Create(gomock.Any()).DoAndReturn(func(product *ZorgTechProduct) *gorm.DB {
		assert.Equal(t, input.Naam, product.Naam)
		assert.Equal(t, input.Beschrijving, product.Beschrijving)
		assert.Equal(t, input.Categorie, product.Categorie)
		assert.Equal(t, input.Prijs, product.Prijs)
		assert.Equal(t, input.Leverancier, product.Leverancier)
		assert.True(t, product.IsActief)
		return &gorm.DB{Error: nil}
	})

	// Mock cache invalidation
	mockCache.EXPECT().Del(ctx, "alle_producten").Return(redis.NewIntResult(1, nil))
	mockCache.EXPECT().Del(ctx, "categorie_"+input.Categorie).Return(redis.NewIntResult(1, nil))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/catalogus/product", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response struct {
		Data ZorgTechProduct `json:"data"`
	}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, input.Naam, response.Data.Naam)
	assert.Equal(t, input.Categorie, response.Data.Categorie)
}

func TestWijzigZorgTechProduct(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := database.NewMockDatabase(ctrl)
	mockCache := cache.NewMockCache(ctrl)
	ctx := context.Background()

	repo := NewZorgTechProductRepository(mockDB, mockCache, &ctx)

	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.PUT("/catalogus/product", repo.WijzigZorgTechProduct)

	productID := uuid.New()
	newNaam := "Gewijzigde Naam"
	newPrijs := 129.99
	input := WijzigZorgTechProductCommand{
		ProductID: productID,
		Naam:      &newNaam,
		Prijs:     &newPrijs,
	}
	requestBody, _ := json.Marshal(input)

	// Mock database find
	mockDB.EXPECT().Where("id = ?", productID).Return(mockDB).Times(1)
	mockDB.EXPECT().First(gomock.Any()).DoAndReturn(func(product *ZorgTechProduct) *gorm.DB {
		product.ID = productID
		product.Naam = "Oorspronkelijke Naam"
		product.Prijs = 99.99
		product.Categorie = "valpreventie"
		product.IsActief = true
		return &gorm.DB{Error: nil}
	}).Times(1)

	// Mock database save
	mockDB.EXPECT().Save(gomock.Any()).Return(&gorm.DB{Error: nil}).Times(1)

	// Mock cache invalidation
	mockCache.EXPECT().Del(ctx, "product_"+productID.String()).Return(redis.NewIntResult(1, nil))
	mockCache.EXPECT().Del(ctx, "alle_producten").Return(redis.NewIntResult(1, nil))
	mockCache.EXPECT().Del(ctx, "categorie_valpreventie").Return(redis.NewIntResult(1, nil))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/catalogus/product", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response struct {
		Data ZorgTechProduct `json:"data"`
	}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, newNaam, response.Data.Naam)
	assert.Equal(t, newPrijs, response.Data.Prijs)
}

func TestVerwijderZorgTechProduct(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := database.NewMockDatabase(ctrl)
	mockCache := cache.NewMockCache(ctrl)
	ctx := context.Background()

	repo := NewZorgTechProductRepository(mockDB, mockCache, &ctx)

	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.DELETE("/catalogus/product", repo.VerwijderZorgTechProduct)

	productID := uuid.New()
	categorie := "valpreventie"
	input := VerwijderZorgTechProductCommand{
		ProductID: productID,
	}
	requestBody, _ := json.Marshal(input)

	// Mock database find
	mockDB.EXPECT().Where("id = ?", productID).Return(mockDB).Times(1)
	mockDB.EXPECT().First(gomock.Any()).DoAndReturn(func(product *ZorgTechProduct) *gorm.DB {
		product.ID = productID
		product.Naam = "Test Product"
		product.Categorie = categorie
		product.IsActief = true
		return &gorm.DB{Error: nil}
	}).Times(1)

	// Mock database save
	mockDB.EXPECT().Save(gomock.Any()).Return(&gorm.DB{Error: nil}).Times(1)

	// Mock cache invalidation
	mockCache.EXPECT().Del(ctx, "product_"+productID.String()).Return(redis.NewIntResult(1, nil))
	mockCache.EXPECT().Del(ctx, "alle_producten").Return(redis.NewIntResult(1, nil))
	mockCache.EXPECT().Del(ctx, "categorie_"+categorie).Return(redis.NewIntResult(1, nil))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/catalogus/product", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response struct {
		Data ZorgTechProduct `json:"data"`
	}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.False(t, response.Data.IsActief)
}

func TestVoegTechnischDetailToe(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := database.NewMockDatabase(ctrl)
	mockCache := cache.NewMockCache(ctrl)
	ctx := context.Background()

	repo := NewZorgTechProductRepository(mockDB, mockCache, &ctx)

	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.POST("/catalogus/product/technisch-detail", repo.VoegTechnischDetailToe)

	productID := uuid.New()
	input := VoegTechnischDetailToeCommand{
		ProductID: productID,
		Sleutel:   "compatibiliteit",
		Waarde:    "Windows 10+",
	}
	requestBody, _ := json.Marshal(input)

	// Mock database find
	mockDB.EXPECT().Where("id = ?", productID).Return(mockDB).Times(1)
	mockDB.EXPECT().First(gomock.Any()).DoAndReturn(func(product *ZorgTechProduct) *gorm.DB {
		product.ID = productID
		product.TechnischeDetails = []TechnischDetail{}
		return &gorm.DB{Error: nil}
	}).Times(1)

	// Mock database save
	mockDB.EXPECT().Save(gomock.Any()).Return(&gorm.DB{Error: nil}).Times(1)

	// Mock cache invalidation
	mockCache.EXPECT().Del(ctx, "product_"+productID.String()).Return(redis.NewIntResult(1, nil))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/catalogus/product/technisch-detail", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response struct {
		Data ZorgTechProduct `json:"data"`
	}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Len(t, response.Data.TechnischeDetails, 1)
	assert.Equal(t, input.Sleutel, response.Data.TechnischeDetails[0].Sleutel)
	assert.Equal(t, input.Waarde, response.Data.TechnischeDetails[0].Waarde)
}

func TestGetProductById(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := database.NewMockDatabase(ctrl)
	mockCache := cache.NewMockCache(ctrl)
	ctx := context.Background()

	repo := NewZorgTechProductRepository(mockDB, mockCache, &ctx)

	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.GET("/catalogus/product/:zorgtechId", repo.GetProductById)

	productID := uuid.New()
	product := ZorgTechProduct{
		ID:       productID,
		Naam:     "Test Product",
		Categorie: "valpreventie",
		IsActief: true,
	}

	// Mock cache get
	mockCache.EXPECT().Get(ctx, "product_"+productID.String()).
		Return(redis.NewStringResult("", redis.Nil)) // Cache miss

	// Mock database find
	mockDB.EXPECT().Where("id = ? AND is_actief = ?", productID, true).
		Return(mockDB).Times(1)
	mockDB.EXPECT().First(gomock.Any()).
		DoAndReturn(func(p *ZorgTechProduct) *gorm.DB {
			*p = product
			return &gorm.DB{Error: nil}
		}).Times(1)

	// Mock cache set
	serialized, _ := json.Marshal(product)
	mockCache.EXPECT().Set(ctx, "product_"+productID.String(), serialized, time.Hour).
		Return(redis.NewStatusResult("OK", nil))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/catalogus/product/"+productID.String(), nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response struct {
		Data ZorgTechProduct `json:"data"`
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

	repo := NewZorgTechProductRepository(mockDB, mockCache, &ctx)

	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.GET("/catalogus/categorie/:categorie", repo.FindByCategorie)

	categorie := "valpreventie"
	producten := []ZorgTechProduct{
		{
			ID:       uuid.New(),
			Naam:     "Product 1",
			Categorie: categorie,
			IsActief: true,
		},
		{
			ID:       uuid.New(),
			Naam:     "Product 2",
			Categorie: categorie,
			IsActief: true,
		},
	}

	// Mock cache get
	mockCache.EXPECT().Get(ctx, "categorie_"+categorie).
		Return(redis.NewStringResult("", redis.Nil)) // Cache miss

	// Mock database find
	mockDB.EXPECT().Where("categorie = ? AND is_actief = ?", categorie, true).
		Return(mockDB).Times(1)
	mockDB.EXPECT().Find(gomock.Any()).
		DoAndReturn(func(p *[]ZorgTechProduct) *gorm.DB {
			*p = producten
			return &gorm.DB{Error: nil}
		}).Times(1)

	// Mock cache set
	serialized, _ := json.Marshal(producten)
	mockCache.EXPECT().Set(ctx, "categorie_"+categorie, serialized, 30*time.Minute).
		Return(redis.NewStatusResult("OK", nil))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/catalogus/categorie/"+categorie, nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response struct {
		Data []ZorgTechProduct `json:"data"`
	}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Len(t, response.Data, 2)
	assert.Equal(t, categorie, response.Data[0].Categorie)
}

// Mock voor de ZorgTechProductRepository interface
type MockZorgTechProductRepository struct {
	ctrl     *gomock.Controller
	recorder *MockZorgTechProductRepositoryMockRecorder
}

func NewMockZorgTechProductRepository(ctrl *gomock.Controller) *MockZorgTechProductRepository {
	mock := &MockZorgTechProductRepository{ctrl: ctrl}
	mock.recorder = &MockZorgTechProductRepositoryMockRecorder{mock}
	return mock
}

func (m *MockZorgTechProductRepository) EXPECT() *MockZorgTechProductRepositoryMockRecorder {
	return m.recorder
}

type MockZorgTechProductRepositoryMockRecorder struct {
	mock *MockZorgTechProductRepository
}

func (m *MockZorgTechProductRepository) Healthcheck(c *gin.Context) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Healthcheck", c)
}

func (mr *MockZorgTechProductRepositoryMockRecorder) Healthcheck(c interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Healthcheck", reflect.TypeOf((*ZorgTechProductRepository)(nil).Healthcheck), c)
}

// Implementeer hier de overige interface methoden voor de mock
// ...