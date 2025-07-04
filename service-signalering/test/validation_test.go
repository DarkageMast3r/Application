package test

import (
	"service-signalering/models"
	"service-signalering/pkg/validation"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestValidateSignaalRequest(t *testing.T) {
	tests := []struct {
		name        string
		signaal     models.Signaal
		expectError bool
		errorCount  int
	}{
		{
			name: "Valid heart rate signal",
			signaal: models.Signaal{
				Type:     "hartslag",
				Waarde:   75.0,
				Tijdstip: time.Now().Add(-5 * time.Minute),
				Bron:     "sensor",
			},
			expectError: false,
			errorCount:  0,
		},
		{
			name: "Heart rate too high",
			signaal: models.Signaal{
				Type:     "hartslag",
				Waarde:   500.0,
				Tijdstip: time.Now().Add(-5 * time.Minute),
				Bron:     "sensor",
			},
			expectError: true,
			errorCount:  1,
		},
		{
			name: "Invalid signal type",
			signaal: models.Signaal{
				Type:     "invalid_type",
				Waarde:   75.0,
				Tijdstip: time.Now().Add(-5 * time.Minute),
				Bron:     "sensor",
			},
			expectError: true,
			errorCount:  1,
		},
		{
			name: "XSS attempt in source",
			signaal: models.Signaal{
				Type:     "hartslag",
				Waarde:   75.0,
				Tijdstip: time.Now().Add(-5 * time.Minute),
				Bron:     "<script>alert('hack')</script>",
			},
			expectError: true,
			errorCount:  1,
		},
		{
			name: "Future timestamp",
			signaal: models.Signaal{
				Type:     "hartslag",
				Waarde:   75.0,
				Tijdstip: time.Now().Add(2 * time.Hour),
				Bron:     "sensor",
			},
			expectError: true,
			errorCount:  1,
		},
		{
			name: "Multiple errors",
			signaal: models.Signaal{
				Type:     "invalid_type",
				Waarde:   -100.0,
				Tijdstip: time.Now().Add(2 * time.Hour),
				Bron:     "<script>alert('hack')</script>",
			},
			expectError: true,
			errorCount:  4,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errors := validation.ValidateSignaalRequest(tt.signaal)

			if tt.expectError {
				assert.NotEmpty(t, errors)
				assert.Len(t, errors, tt.errorCount)
			} else {
				assert.Empty(t, errors)
			}
		})
	}
}

func TestBloodPressureValidation(t *testing.T) {
	tests := []struct {
		name        string
		signaalType string
		waarde      float64
		expectError bool
	}{
		{"Valid systolic BP", "bloeddruk_systolisch", 120.0, false},
		{"Valid diastolic BP", "bloeddruk_diastolisch", 80.0, false},
		{"High systolic BP", "bloeddruk_systolisch", 400.0, true},
		{"Low diastolic BP", "bloeddruk_diastolisch", 10.0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			signaal := models.Signaal{
				Type:     tt.signaalType,
				Waarde:   tt.waarde,
				Tijdstip: time.Now().Add(-5 * time.Minute),
				Bron:     "monitor",
			}

			errors := validation.ValidateSignaalRequest(signaal)

			if tt.expectError {
				assert.NotEmpty(t, errors)
			} else {
				assert.Empty(t, errors)
			}
		})
	}
}

func TestValidateClassificatie(t *testing.T) {
	tests := []struct {
		name           string
		classificatie  models.ToestandClassificatie
		expectError    bool
		expectedErrors int
	}{
		{
			name: "Valid classification",
			classificatie: models.ToestandClassificatie{
				Categorie: "cardiovasculair",
				Ernst:     "normaal",
				Motivatie: "Patient vertoont stabiele vitale functies",
			},
			expectError:    false,
			expectedErrors: 0,
		},
		{
			name: "Invalid category",
			classificatie: models.ToestandClassificatie{
				Categorie: "invalid_category",
				Ernst:     "normaal",
				Motivatie: "Patient vertoont stabiele vitale functies",
			},
			expectError:    true,
			expectedErrors: 1,
		},
		{
			name: "Invalid severity",
			classificatie: models.ToestandClassificatie{
				Categorie: "cardiovasculair",
				Ernst:     "super_critical",
				Motivatie: "Patient vertoont stabiele vitale functies",
			},
			expectError:    true,
			expectedErrors: 1,
		},
		{
			name: "Motivation too short",
			classificatie: models.ToestandClassificatie{
				Categorie: "cardiovasculair",
				Ernst:     "normaal",
				Motivatie: "Short",
			},
			expectError:    true,
			expectedErrors: 1,
		},
		{
			name: "XSS in motivation",
			classificatie: models.ToestandClassificatie{
				Categorie: "cardiovasculair",
				Ernst:     "normaal",
				Motivatie: "Patient has <script>alert('xss')</script> condition",
			},
			expectError:    true,
			expectedErrors: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errors := validation.ValidateClassificatie(tt.classificatie)

			if tt.expectError {
				assert.NotEmpty(t, errors)
				assert.Len(t, errors, tt.expectedErrors)
			} else {
				assert.Empty(t, errors)
			}
		})
	}
}

func TestValidateBeoordelingRequest(t *testing.T) {
	tests := []struct {
		name           string
		request        models.BeoordeelSituatieRequest
		expectError    bool
		expectedErrors int
	}{
		{
			name: "Valid assessment",
			request: models.BeoordeelSituatieRequest{
				Conclusie:       "Patient is in stable condition with normal vital signs",
				Urgentie:        "laag",
				GevalideerdDoor: "Dr. van der Berg",
			},
			expectError:    false,
			expectedErrors: 0,
		},
		{
			name: "Conclusion too short",
			request: models.BeoordeelSituatieRequest{
				Conclusie:       "Short",
				Urgentie:        "laag",
				GevalideerdDoor: "Dr. van der Berg",
			},
			expectError:    true,
			expectedErrors: 1,
		},
		{
			name: "Invalid urgency",
			request: models.BeoordeelSituatieRequest{
				Conclusie:       "Patient is in stable condition with normal vital signs",
				Urgentie:        "mega_urgent",
				GevalideerdDoor: "Dr. van der Berg",
			},
			expectError:    true,
			expectedErrors: 1,
		},
		{
			name: "Invalid validator name",
			request: models.BeoordeelSituatieRequest{
				Conclusie:       "Patient is in stable condition with normal vital signs",
				Urgentie:        "laag",
				GevalideerdDoor: "<script>alert('hack')</script>",
			},
			expectError:    true,
			expectedErrors: 1,
		},
		{
			name: "Empty validator",
			request: models.BeoordeelSituatieRequest{
				Conclusie:       "Patient is in stable condition with normal vital signs",
				Urgentie:        "laag",
				GevalideerdDoor: "",
			},
			expectError:    true,
			expectedErrors: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errors := validation.ValidateBeoordelingRequest(tt.request)

			if tt.expectError {
				assert.NotEmpty(t, errors)
				assert.Len(t, errors, tt.expectedErrors)
			} else {
				assert.Empty(t, errors)
			}
		})
	}
}
