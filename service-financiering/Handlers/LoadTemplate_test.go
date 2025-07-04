package handlers_test

import (
	"log"
	httptest "net/http/httptest"
	"testing"
	h "Financiering/Handlers"
	U "Financiering/Utilities"
)

func Test_LoadTempateWithMissingPath(t *testing.T) {
	//Setup
	err := U.FindDir()
	if err != nil {
		log.Fatal(err)
	}
	rr := httptest.NewRecorder()

	//Execute
	err = h.LoadTemplate(rr, "", nil)

	//Assert
	if err == nil {
		t.Errorf("handler didnt return err \n")
	}
} 