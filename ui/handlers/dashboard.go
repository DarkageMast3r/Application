package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type DashboardData struct {
	Title       string
	CurrentPage string
	CSRFToken   string
	Stats       DashboardStats
	Activities  []Activity
}

type DashboardStats struct {
	ActiveClients  int     `json:"active_clients"`
	PendingSignals int     `json:"pending_signals"`
	ActiveTech     int     `json:"active_tech"`
	BudgetUsed     float64 `json:"budget_used"`
	TotalBudget    float64 `json:"total_budget"`
}

type Activity struct {
	ID          string    `json:"id"`
	Description string    `json:"description"`
	ClientID    string    `json:"client_id"`
	Timestamp   time.Time `json:"timestamp"`
	Type        string    `json:"type"`
}

func Dashboard(c *gin.Context) {
	data := DashboardData{
		Title:       "Dashboard",
		CurrentPage: "dashboard",
		CSRFToken:   generateCSRFToken(c),
		Stats:       getDashboardStats(),
		Activities:  getRecentActivities(),
	}
	c.HTML(http.StatusOK, "layout.html", data)
}

func DashboardPage(c *gin.Context) {
	data := DashboardData{
		Stats:      getDashboardStats(),
		Activities: getRecentActivities(),
	}
	c.HTML(http.StatusOK, "dashboard.html", data)
}

func GetDashboardStats(c *gin.Context) {
	stats := getDashboardStats()
	c.JSON(http.StatusOK, stats)
}

func GetRecentActivities(c *gin.Context) {
	activities := getRecentActivities()
	c.JSON(http.StatusOK, activities)
}

func getDashboardStats() DashboardStats {
	return DashboardStats{
		ActiveClients:  42,
		PendingSignals: 18,
		ActiveTech:     34,
		BudgetUsed:     24500,
		TotalBudget:    50000,
	}
}

func getRecentActivities() []Activity {
	return []Activity{
		{
			ID:          "A001",
			Description: "Valdetector toegevoegd aan zorgplan",
			ClientID:    "C001",
			Timestamp:   time.Now().Add(-2 * time.Hour),
			Type:        "technology_added",
		},
		{
			ID:          "A002",
			Description: "Budget aanvraag goedgekeurd voor slimme medicijndispenser",
			ClientID:    "C005",
			Timestamp:   time.Now().Add(-4 * time.Hour),
			Type:        "budget_approved",
		},
		{
			ID:          "A003",
			Description: "Achteruitgang signaal geregistreerd",
			ClientID:    "C007",
			Timestamp:   time.Now().Add(-6 * time.Hour),
			Type:        "signal_created",
		},
	}
}

func generateCSRFToken(c *gin.Context) string {
	return "csrf_token_placeholder"
}
