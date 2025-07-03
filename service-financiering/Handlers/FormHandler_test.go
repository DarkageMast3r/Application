package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	h "Financiering/Handlers"
	U "Financiering/Utilities"
)

// TestAddDossierMissingParams tests the case where parameters are missing
func Test_AddDossierMissingParams(t *testing.T) {
	U.StartTest()

	// Setup
	req := httptest.NewRequest(http.MethodPost, "/add-dossier", nil)
	rr := httptest.NewRecorder()

	// Execute
	h.AddDossier(rr, req)

	// Assert
	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusInternalServerError)
	}
}
