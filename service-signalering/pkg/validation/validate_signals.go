package validation

import (
	"fmt"
	"service-signalering/models"
	"strings"
	"time"
)

// De grote functie die alles checkt
func ValidateSignaalRequest(signaal models.Signaal) []string {
	var errors []string

	if !isValidSignalType(signaal.Type) {
		errors = append(errors, fmt.Sprintf("Ongeldig signaal type: %s", signaal.Type))
	}

	if valueError := validateSignalValue(signaal.Type, signaal.Waarde); valueError != "" {
		errors = append(errors, valueError)
	}

	if !isValidBron(signaal.Bron) {
		errors = append(errors, fmt.Sprintf("Ongeldige bron: %s", signaal.Bron))
	}

	if timestampError := validateTijdstip(signaal.Tijdstip); timestampError != "" {
		errors = append(errors, timestampError)
	}

	return errors
}

func ValidateSignalenRequest(request models.RegistreerAchteruitgangRequest) []string {
	var errors []string

	if len(request.Signalen) == 0 {
		errors = append(errors, "Geen signalen ontvangen")
		return errors
	}

	if len(request.Signalen) > 50 {
		errors = append(errors, "Te veel signalen in één request (max 50)")
	}

	for i, signaal := range request.Signalen {
		signaalErrors := ValidateSignaalRequest(signaal)
		for _, err := range signaalErrors {
			errors = append(errors, fmt.Sprintf("Signaal %d: %s", i+1, err))
		}
	}

	return errors
}

// Alleen requests van dit soort
// In de praktijk zouden we hier waarschijnlijk een meer modulair systeem voor moeten vinden
// Maar voor onze showcase API is het goed genoeg
func isValidSignalType(signalType string) bool {
	validTypes := []string{
		"hartslag",
		"bloeddruk_systolisch",
		"bloeddruk_diastolisch",
		"temperatuur",
		"saturatie",
		"glucose",
		"gewicht",
		"pijn_score",
		"ademhaling",
	}

	for _, validType := range validTypes {
		if signalType == validType {
			return true
		}
	}
	return false
}

func validateSignalValue(signalType string, waarde float64) string {
	if waarde < 0 {
		return "Waarde kan niet negatief zijn"
	}

	// Specifieke validatie, voeg meer toe/verander als nodig
	switch signalType {
	case "hartslag":
		if waarde < 20 || waarde > 300 {
			return fmt.Sprintf("Hartslag %.1f is onrealistisch (verwacht 20-300 bpm)", waarde)
		}
	case "bloeddruk_systolisch":
		if waarde < 50 || waarde > 300 {
			return fmt.Sprintf("Systolische bloeddruk %.1f is onrealistisch (verwacht 50-300 mmHg)", waarde)
		}
	case "bloeddruk_diastolisch":
		if waarde < 20 || waarde > 200 {
			return fmt.Sprintf("Diastolische bloeddruk %.1f is onrealistisch (verwacht 20-200 mmHg)", waarde)
		}
	case "temperatuur":
		if waarde < 30.0 || waarde > 50.0 {
			return fmt.Sprintf("Temperatuur %.1f is onrealistisch (verwacht 30-50°C)", waarde)
		}
	case "saturatie":
		if waarde < 0 || waarde > 100 {
			return fmt.Sprintf("Saturatie %.1f is ongeldig (verwacht 0-100%%)", waarde)
		}
	case "glucose":
		if waarde < 0 || waarde > 50 {
			return fmt.Sprintf("Glucose %.1f is onrealistisch (verwacht 0-50 mmol/L)", waarde)
		}
	case "gewicht":
		if waarde < 0.5 || waarde > 1000 {
			return fmt.Sprintf("Gewicht %.1f is onrealistisch (verwacht 0.5-1000 kg)", waarde)
		}
	case "pijn_score":
		if waarde < 0 || waarde > 10 {
			return fmt.Sprintf("Pijn score %.1f is ongeldig (verwacht 0-10)", waarde)
		}
	case "ademhaling":
		if waarde < 5 || waarde > 100 {
			return fmt.Sprintf("Ademhaling %.1f is onrealistisch (verwacht 5-100 per minuut)", waarde)
		}
	}

	return ""
}

func isValidBron(bron string) bool {
	// Geen te lange requests
	if len(bron) == 0 || len(bron) > 100 {
		return false
	}

	// XSS voorkomen
	amongus := []string{"<script", "javascript", "eval(", "exec(", "<iframe", "onload=", "onerror="}
	bronLower := strings.ToLower(bron)

	for _, sus := range amongus {
		if strings.Contains(bronLower, sus) {
			return false
		}
	}

	return true
}

// validateTijdstip checks if timestamp is reasonable
func validateTijdstip(tijdstip time.Time) string {
	now := time.Now()

	if tijdstip.After(now.Add(time.Hour)) {
		return "Tijdstip ligt te ver in de toekomst"
	}

	if tijdstip.Before(now.AddDate(-20, 0, 0)) {
		return "Tijdstip ligt te ver in het verleden"
	}

	return ""
}
