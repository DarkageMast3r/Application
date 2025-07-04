package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"

	"smartcare/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.SetFuncMap(template.FuncMap{
		"formatCurrency": formatCurrency,
		"formatTimeAgo":  formatTimeAgo,
	})
	r.LoadHTMLGlob("templates/*.html")
	r.Static("/static", "./static")
	setupRoutes(r)

	// Comment out these lines:
	// service.Init()
	// port := service.Register("app", func(w http.ResponseWriter, req *http.Request) {
	//     r.ServeHTTP(w, req)
	// })

	// Use a fixed port:
	port := 8080
	log.Println("SmartCare Assist server starting on :", port)
	http.ListenAndServe(":"+strconv.Itoa(port), r)
}

func setupRoutes(r *gin.Engine) {
	r.GET("/", handlers.Dashboard)

	api := r.Group("/api")
	{
		api.GET("/dashboard/stats", handlers.GetDashboardStats)
		api.GET("/dashboard/activities", handlers.GetRecentActivities)
		api.POST("/signals/create", handlers.CreateSignal)
		api.POST("/finance/request", handlers.RequestBudget)
		api.POST("/reports/generate", handlers.GenerateReport)
	}

	pages := r.Group("/pages")
	{
		pages.GET("/dashboard", handlers.DashboardPage)
		pages.GET("/signalen", handlers.SignalenPage)
		pages.GET("/zorgtechnologie", handlers.ZorgtechnologiePage)
		pages.GET("/financiering", handlers.FinancieringPage)
		pages.GET("/rapportage", handlers.RapportagePage)
	}
}

func formatCurrency(amount float64) string {
	return fmt.Sprintf("€%.1fk", amount/1000)
}

func formatTimeAgo(t time.Time) string {
	duration := time.Since(t)
	if duration < time.Hour {
		return fmt.Sprintf("%.0f minuten geleden", duration.Minutes())
	} else if duration < 24*time.Hour {
		return fmt.Sprintf("%.0f uur geleden", duration.Hours())
	}
	return fmt.Sprintf("%.0f dagen geleden", duration.Hours()/24)
}
