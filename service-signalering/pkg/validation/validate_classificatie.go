package validation

import (
	"fmt"
	"service-signalering/models"
	"strings"
)

func ValidateClassificatie(classificatie models.ToestandClassificatie) []string {
	var errors []string

	if !isValidCategorie(classificatie.Categorie) {
		errors = append(errors, fmt.Sprintf("Ongeldige categorie: %s", classificatie.Categorie))
	}

	if !isValidErnst(classificatie.Ernst) {
		errors = append(errors, fmt.Sprintf("Ongeldige ernst: %s", classificatie.Ernst))
	}

	if motivatieError := validateMotivatie(classificatie.Motivatie); motivatieError != "" {
		errors = append(errors, motivatieError)
	}

	return errors
}

func isValidCategorie(categorie string) bool {
	validCategories := []string{
		"cardiovasculair",
		"respiratoir",
		"neurologisch",
		"metabolisch",
		"infectieus",
		"psychiatrisch",
		"traumatisch",
		"algemeen",
	}

	for _, valid := range validCategories {
		if categorie == valid {
			return true
		}
	}
	return false
}

func isValidErnst(ernst string) bool {
	validLevels := []string{
		"laag",
		"normaal",
		"licht",
		"matig",
		"ernstig",
		"kritiek",
		"levensbedreignd",
	}

	for _, valid := range validLevels {
		if ernst == valid {
			return true
		}
	}
	return false
}

func validateMotivatie(motivatie string) string {
	if len(motivatie) == 0 {
		return "Motivatie is verplicht"
	}

	if len(motivatie) < 10 {
		return "Motivatie is te kort (minimaal 10 karakters)"
	}

	if len(motivatie) > 1000 {
		return "Motivatie is te lang (maximaal 1000 karakters)"
	}

	amongus := []string{"<script", "javascript:", "eval(", "exec(", "<iframe", "onload=", "onerror="}
	motivatieLower := strings.ToLower(motivatie)

	for _, sus := range amongus {
		if strings.Contains(motivatieLower, sus) {
			return "Motivatie bevat niet toegestane inhoud"
		}
	}

	return ""
}
