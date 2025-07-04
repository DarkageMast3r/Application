package test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"service-signalering/handlers"
	"service-signalering/models"
	"service-signalering/pkg/database"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func createTestContext(method, path string, body interface{}) (*gin.Context, *httptest.ResponseRecorder) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	var reqBody []byte
	if body != nil {
		reqBody, _ = json.Marshal(body)
	}

	req, _ := http.NewRequest(method, path, bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", "dev-key-123")
	c.Request = req

	return c, w
}

func TestAddSignals_ValidInput(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	originalDB := database.DB
	database.DB = db
	defer func() { database.DB = originalDB }()

	clientID := uuid.New()
	mock.ExpectExec("INSERT INTO signals").
		WithArgs(clientID, "hartslag", 75.0, sqlmock.AnyArg(), "sensor").
		WillReturnResult(sqlmock.NewResult(1, 1))

	request := models.RegistreerAchteruitgangRequest{
		Signalen: []models.Signaal{
			{
				Type:     "hartslag",
				Waarde:   75.0,
				Tijdstip: time.Now().Add(-5 * time.Minute),
				Bron:     "sensor",
			},
		},
	}

	c, w := createTestContext("POST", "/api/v1/clients/"+clientID.String()+"/signals", request)
	c.Params = []gin.Param{{Key: "id", Value: clientID.String()}}

	handlers.AddSignals(c)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestAddSignals_InvalidUUID(t *testing.T) {
	c, w := createTestContext("POST", "/api/v1/clients/invalid-uuid/signals", nil)
	c.Params = []gin.Param{{Key: "id", Value: "invalid-uuid"}}

	handlers.AddSignals(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "INVALID_CLIENT_ID")
}

func TestAddSignals_ValidationFailure(t *testing.T) {
	clientID := uuid.New()

	request := models.RegistreerAchteruitgangRequest{
		Signalen: []models.Signaal{
			{
				Type:     "hartslag",
				Waarde:   999.0,
				Tijdstip: time.Now().Add(-5 * time.Minute),
				Bron:     "sensor",
			},
		},
	}

	c, w := createTestContext("POST", "/api/v1/clients/"+clientID.String()+"/signals", request)
	c.Params = []gin.Param{{Key: "id", Value: clientID.String()}}

	handlers.AddSignals(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "VALIDATION_FAILED")
}

func TestClassifyCondition_ValidInput(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	originalDB := database.DB
	database.DB = db
	defer func() { database.DB = originalDB }()

	clientID := uuid.New()
	mock.ExpectExec("INSERT INTO classifications").
		WithArgs(clientID, "cardiovasculair", "normaal", "Patient vertoont stabiele vitale functies").
		WillReturnResult(sqlmock.NewResult(1, 1))

	request := models.ClassificeerSituatieRequest{
		Classificatie: models.ToestandClassificatie{
			Categorie: "cardiovasculair",
			Ernst:     "normaal",
			Motivatie: "Patient vertoont stabiele vitale functies",
		},
	}

	c, w := createTestContext("POST", "/api/v1/clients/"+clientID.String()+"/classify", request)
	c.Params = []gin.Param{{Key: "id", Value: clientID.String()}}

	handlers.ClassifyCondition(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestAssessCondition_ValidInput(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	originalDB := database.DB
	database.DB = db
	defer func() { database.DB = originalDB }()

	clientID := uuid.New()
	expectedConclusion := "Patient is stable and shows no signs of deterioration"
	mock.ExpectExec("INSERT INTO assessments").
		WithArgs(clientID, expectedConclusion, "laag", "Dr. Test", sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	request := models.BeoordeelSituatieRequest{
		Conclusie:       expectedConclusion,
		Urgentie:        "laag",
		GevalideerdDoor: "Dr. Test",
	}

	c, w := createTestContext("POST", "/api/v1/clients/"+clientID.String()+"/assess", request)
	c.Params = []gin.Param{{Key: "id", Value: clientID.String()}}

	handlers.AssessCondition(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetClientCondition_NoData(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	originalDB := database.DB
	database.DB = db
	defer func() { database.DB = originalDB }()

	clientID := uuid.New()

	mock.ExpectQuery("SELECT type, waarde, tijdstip, bron FROM signals").
		WithArgs(clientID).
		WillReturnRows(sqlmock.NewRows([]string{"type", "waarde", "tijdstip", "bron"}))

	mock.ExpectQuery("SELECT categorie, ernst, motivatie FROM classifications").
		WithArgs(clientID).
		WillReturnError(sql.ErrNoRows)

	mock.ExpectQuery("SELECT conclusie, urgentie, gevalideerd_door, tijdstip FROM assessments").
		WithArgs(clientID).
		WillReturnError(sql.ErrNoRows)

	c, w := createTestContext("GET", "/api/v1/clients/"+clientID.String()+"/condition", nil)
	c.Params = []gin.Param{{Key: "id", Value: clientID.String()}}

	handlers.GetClientCondition(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.ToestandResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Empty(t, response.Signalen)
	assert.Nil(t, response.Classificatie)
	assert.Nil(t, response.Beoordeling)
	assert.Equal(t, clientID, response.ClientID)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestHandlers_InvalidJSON(t *testing.T) {
	clientID := uuid.New()

	tests := []struct {
		name    string
		handler gin.HandlerFunc
		path    string
	}{
		{"AddSignals invalid JSON", handlers.AddSignals, "/signals"},
		{"ClassifyCondition invalid JSON", handlers.ClassifyCondition, "/classify"},
		{"AssessCondition invalid JSON", handlers.AssessCondition, "/assess"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			req, _ := http.NewRequest("POST", tt.path, bytes.NewBufferString("{invalid json"))
			req.Header.Set("Content-Type", "application/json")
			c.Request = req
			c.Params = []gin.Param{{Key: "id", Value: clientID.String()}}

			tt.handler(c)

			assert.Equal(t, http.StatusBadRequest, w.Code)
			assert.Contains(t, w.Body.String(), "INVALID_REQUEST")
		})
	}
}
