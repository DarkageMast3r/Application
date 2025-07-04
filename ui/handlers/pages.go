package handlers

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"smartcare/global"
	"smartcare/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

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
	c.Request.ParseForm()
	request := make(map[string]string)
	request["name"] = c.Request.Form.Get("name")
	request["client_id"] = c.Request.Form.Get("client_id")
	request["description"] = c.Request.Form.Get("description")
	requestJson, err := json.Marshal(request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": ""})
		return
	}
	url := "https://" + global.Config.Service_discovery_root + ":" + strconv.Itoa(global.Config.Service_discovery_port) + "/Selection/Case/Create"
	_, err = http.Post(url, "text/json", bytes.NewReader(requestJson))
	if err != nil {
		service.LogWarning("Could not send message to URL", url, err)
	}
	c.Redirect(http.StatusSeeOther, "/App")
}

func RequestBudget(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Aanvraag ingediend"})
}

func GenerateReport(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Rapport gegenereerd"})
}
