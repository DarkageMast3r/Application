package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	amqp "github.com/rabbitmq/amqp091-go"
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

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()
	q, err := ch.QueueDeclare(
		"hello2", // name
		false,    // durable
		false,    // delete when unused
		false,    // exclusive
		false,    // no-wait
		nil,      // arguments
	)
	failOnError(err, "Failed to declare a queue")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	body := "Hello World!"
	err = ch.PublishWithContext(ctx,
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	failOnError(err, "Failed to publish a message")
	log.Printf(" [x] Sent %s\n", body)
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

	http.Post("localhost/Selection/Case/Create", "text/json", bytes.NewReader(requestJson))
	c.Redirect(http.StatusSeeOther, "/App")
}

func RequestBudget(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Aanvraag ingediend"})
}

func GenerateReport(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Rapport gegenereerd"})
}
