package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func SignalenPage(c *gin.Context) {
	data := DashboardData{
		Title:       "Signalen",
		CurrentPage: "signalen",
		CSRFToken:   generateCSRFToken(c),
	}

	if c.GetHeader("X-Requested-With") == "XMLHttpRequest" {
		c.HTML(http.StatusOK, "signalen.html", data)
	} else {
		c.HTML(http.StatusOK, "layout.html", data)
	}
}

func ZorgtechnologiePage(c *gin.Context) {
	data := DashboardData{
		Title:       "Zorgtechnologie",
		CurrentPage: "zorgtechnologie",
		CSRFToken:   generateCSRFToken(c),
	}

	if c.GetHeader("X-Requested-With") == "XMLHttpRequest" {
		c.HTML(http.StatusOK, "zorgtechnologie.html", data)
	} else {
		c.HTML(http.StatusOK, "layout.html", data)
	}
}

func FinancieringPage(c *gin.Context) {
	data := DashboardData{
		Title:       "Financiering",
		CurrentPage: "financiering",
		CSRFToken:   generateCSRFToken(c),
	}

	if c.GetHeader("X-Requested-With") == "XMLHttpRequest" {
		c.HTML(http.StatusOK, "financiering.html", data)
	} else {
		c.HTML(http.StatusOK, "layout.html", data)
	}
}

func RapportagePage(c *gin.Context) {
	data := DashboardData{
		Title:       "Rapportage",
		CurrentPage: "rapportage",
		CSRFToken:   generateCSRFToken(c),
	}

	if c.GetHeader("X-Requested-With") == "XMLHttpRequest" {
		c.HTML(http.StatusOK, "rapportage.html", data)
	} else {
		c.HTML(http.StatusOK, "layout.html", data)
	}
}

// obama giving himself a medal.jpg
func CreateSignal(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Signaal geregistreerd"})
}

func RequestBudget(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Aanvraag ingediend"})
}

func GenerateReport(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Rapport gegenereerd"})
}
