package validation

import (
	"fmt"
	"service-signalering/models"
	"strings"
)

// De grote functie die alles checkt
func ValidateBeoordelingRequest(request models.BeoordeelSituatieRequest) []string {
	var errors []string

	if conclusieError := validateConclusie(request.Conclusie); conclusieError != "" {
		errors = append(errors, conclusieError)
	}

	if !isValidUrgentie(request.Urgentie) {
		errors = append(errors, fmt.Sprintf("Ongeldige urgentie: %s", request.Urgentie))
	}

	if validatorError := validateValidator(request.GevalideerdDoor); validatorError != "" {
		errors = append(errors, validatorError)
	}

	return errors
}

func validateConclusie(conclusie string) string {
	if len(conclusie) == 0 {
		return "Conclusie is verplicht"
	}

	if len(conclusie) < 20 {
		return "Conclusie is te kort (minimaal 20 karakters)"
	}

	if len(conclusie) > 2000 {
		return "Conclusie is te lang (maximaal 2000 karakters)"
	}

	amongus := []string{"<script", "javascript:", "eval(", "exec(", "<iframe", "onload=", "onerror="}
	conclusieLower := strings.ToLower(conclusie)

	for _, sus := range amongus {
		if strings.Contains(conclusieLower, sus) {
			return "Conclusie bevat niet toegestane inhoud"
		}
	}

	return ""
}

func isValidUrgentie(urgentie string) bool {
	validLevels := []string{
		"laag",
		"normaal",
		"hoog",
		"urgent",
		"kritiek",
		"spoed",
	}

	for _, valid := range validLevels {
		if urgentie == valid {
			return true
		}
	}
	return false
}

func validateValidator(validator string) string {
	// Length checks
	if len(validator) == 0 {
		return "Validator naam is verplicht"
	}

	if len(validator) < 2 {
		return "Validator naam is te kort (minimaal 2 karakters)"
	}

	if len(validator) > 150 {
		return "Validator naam is te lang (maximaal 150 karakters)"
	}

	if !isReasonableName(validator) {
		return "Validator naam heeft een ongeldig formaat"
	}

	// Security checks
	suspicious := []string{"<", ">", "script", "javascript", "eval", "exec", "'", "\"", ";", "--"}
	validatorLower := strings.ToLower(validator)

	for _, sus := range suspicious {
		if strings.Contains(validatorLower, sus) {
			return "Validator naam bevat niet toegestane karakters"
		}
	}

	return ""
}

func isReasonableName(name string) bool {
	// Een naam moet bestaan uit letters, spaties, punten en streepjes
	allowedChars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ .-"

	for _, char := range name {
		found := false
		for _, allowed := range allowedChars {
			if char == allowed {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	// Een naam moet minstens één letter hebben
	hasLetter := false
	for _, char := range name {
		if (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') {
			hasLetter = true
			break
		}
	}

	return hasLetter
}
