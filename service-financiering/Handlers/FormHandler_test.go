package handlers_test

import (
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	h "Financiering/Handlers"
	U "Financiering/Utilities"
)

// TestAddDossierMissingParams tests the case where parameters are missing
func Test_AddDossierWithMissingParams(t *testing.T) {
	// Setup
	err := U.FindDir()
	if err != nil {
		log.Fatal(err)
	}

	req := httptest.NewRequest(http.MethodPost, "/add-dossier", nil)
	rr := httptest.NewRecorder()

	// Execute
	h.AddDossier(rr, req)

	// Assert
	statusCode := rr.Code
	if statusCode != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v", statusCode, http.StatusInternalServerError)
	}
}
