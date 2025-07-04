package test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"service-signalering/handlers"
	"service-signalering/models"
	"service-signalering/pkg/auth"
	"service-signalering/pkg/database"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

type IntegrationTestSuite struct {
	suite.Suite
	router *gin.Engine
	db     *sql.DB
	mock   sqlmock.Sqlmock
}

func (suite *IntegrationTestSuite) SetupSuite() {
	gin.SetMode(gin.TestMode)

	db, mock, err := sqlmock.New()
	suite.Require().NoError(err)
	suite.db = db
	suite.mock = mock

	database.DB = db

	suite.router = gin.New()
	suite.router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	api := suite.router.Group("/api/v1")
	api.Use(auth.AuthenticateKey())
	{
		clientRoutes := api.Group("/clients/:id")
		{
			clientRoutes.POST("/signals", handlers.AddSignals)
			clientRoutes.POST("/classify", handlers.ClassifyCondition)
			clientRoutes.POST("/assess", handlers.AssessCondition)
			clientRoutes.GET("/condition", handlers.GetClientCondition)
		}
	}
}

func (suite *IntegrationTestSuite) TearDownSuite() {
	suite.db.Close()
}

func (suite *IntegrationTestSuite) TestHealthEndpoint() {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/health", nil)

	suite.router.ServeHTTP(w, req)

	suite.Equal(http.StatusOK, w.Code)
	suite.Contains(w.Body.String(), "ok")
}

func (suite *IntegrationTestSuite) TestAuthenticationRequired() {
	clientID := uuid.New()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/clients/"+clientID.String()+"/condition", nil)

	suite.router.ServeHTTP(w, req)

	suite.Equal(http.StatusUnauthorized, w.Code)
}

func (suite *IntegrationTestSuite) TestCompleteWorkflow() {
	clientID := uuid.New()

	suite.mock.ExpectExec("INSERT INTO signals").
		WithArgs(clientID, "hartslag", 75.0, sqlmock.AnyArg(), "heart_monitor").
		WillReturnResult(sqlmock.NewResult(1, 1))

	suite.mock.ExpectExec("INSERT INTO signals").
		WithArgs(clientID, "bloeddruk_systolisch", 120.0, sqlmock.AnyArg(), "bp_cuff").
		WillReturnResult(sqlmock.NewResult(2, 1))

	suite.mock.ExpectExec("INSERT INTO classifications").
		WithArgs(clientID, "cardiovasculair", "normaal", "Vitale functies binnen normale waarden").
		WillReturnResult(sqlmock.NewResult(3, 1))

	suite.mock.ExpectExec("INSERT INTO assessments").
		WithArgs(clientID, sqlmock.AnyArg(), "laag", "Dr. Integration Test", sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(4, 1))

	signalRows := sqlmock.NewRows([]string{"type", "waarde", "tijdstip", "bron"}).
		AddRow("hartslag", 75.0, time.Now().Add(-5*time.Minute), "heart_monitor").
		AddRow("bloeddruk_systolisch", 120.0, time.Now().Add(-4*time.Minute), "bp_cuff")

	suite.mock.ExpectQuery("SELECT type, waarde, tijdstip, bron FROM signals").
		WithArgs(clientID).
		WillReturnRows(signalRows)

	classificationRows := sqlmock.NewRows([]string{"categorie", "ernst", "motivatie"}).
		AddRow("cardiovasculair", "normaal", "Vitale functies binnen normale waarden")

	suite.mock.ExpectQuery("SELECT categorie, ernst, motivatie FROM classifications").
		WithArgs(clientID).
		WillReturnRows(classificationRows)

	assessmentRows := sqlmock.NewRows([]string{"conclusie", "urgentie", "gevalideerd_door", "tijdstip"}).
		AddRow("Patient vertoont stabiele conditie", "laag", "Dr. Integration Test", time.Now())

	suite.mock.ExpectQuery("SELECT conclusie, urgentie, gevalideerd_door, tijdstip FROM assessments").
		WithArgs(clientID).
		WillReturnRows(assessmentRows)

	signalsRequest := models.RegistreerAchteruitgangRequest{
		Signalen: []models.Signaal{
			{
				Type:     "hartslag",
				Waarde:   75.0,
				Tijdstip: time.Now().Add(-5 * time.Minute),
				Bron:     "heart_monitor",
			},
			{
				Type:     "bloeddruk_systolisch",
				Waarde:   120.0,
				Tijdstip: time.Now().Add(-4 * time.Minute),
				Bron:     "bp_cuff",
			},
		},
	}

	body, _ := json.Marshal(signalsRequest)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/clients/"+clientID.String()+"/signals", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", "dev-key-123")

	suite.router.ServeHTTP(w, req)
	suite.Equal(http.StatusCreated, w.Code)

	classifyRequest := models.ClassificeerSituatieRequest{
		Classificatie: models.ToestandClassificatie{
			Categorie: "cardiovasculair",
			Ernst:     "normaal",
			Motivatie: "Vitale functies binnen normale waarden",
		},
	}

	body, _ = json.Marshal(classifyRequest)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/v1/clients/"+clientID.String()+"/classify", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", "dev-key-123")

	suite.router.ServeHTTP(w, req)
	suite.Equal(http.StatusOK, w.Code)

	assessRequest := models.BeoordeelSituatieRequest{
		Conclusie:       "Patient vertoont stabiele conditie zonder tekenen van achteruitgang",
		Urgentie:        "laag",
		GevalideerdDoor: "Dr. Integration Test",
	}

	body, _ = json.Marshal(assessRequest)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/v1/clients/"+clientID.String()+"/assess", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", "dev-key-123")

	suite.router.ServeHTTP(w, req)
	suite.Equal(http.StatusOK, w.Code)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/v1/clients/"+clientID.String()+"/condition", nil)
	req.Header.Set("X-API-Key", "dev-key-123")

	suite.router.ServeHTTP(w, req)
	suite.Equal(http.StatusOK, w.Code)

	var response models.ToestandResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	suite.NoError(err)

	suite.Equal(clientID, response.ClientID)
	suite.Len(response.Signalen, 2)
	suite.NotNil(response.Classificatie)
	suite.NotNil(response.Beoordeling)
	suite.Equal("cardiovasculair", response.Classificatie.Categorie)
	suite.Equal("laag", response.Beoordeling.Urgentie)

	suite.NoError(suite.mock.ExpectationsWereMet())
}

func (suite *IntegrationTestSuite) TestErrorHandling() {
	clientID := uuid.New()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/clients/invalid-uuid/condition", nil)
	req.Header.Set("X-API-Key", "dev-key-123")

	suite.router.ServeHTTP(w, req)
	suite.Equal(http.StatusBadRequest, w.Code)
	suite.Contains(w.Body.String(), "INVALID_CLIENT_ID")

	invalidSignalsRequest := models.RegistreerAchteruitgangRequest{
		Signalen: []models.Signaal{
			{
				Type:     "invalid_type",
				Waarde:   -100.0,
				Tijdstip: time.Now().Add(2 * time.Hour),
				Bron:     "<script>alert('hack')</script>",
			},
		},
	}

	body, _ := json.Marshal(invalidSignalsRequest)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/v1/clients/"+clientID.String()+"/signals", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", "dev-key-123")

	suite.router.ServeHTTP(w, req)
	suite.Equal(http.StatusBadRequest, w.Code)
	suite.Contains(w.Body.String(), "VALIDATION_FAILED")
}

func TestIntegrationSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}
