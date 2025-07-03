package handlers

import (
	"net/http"
	"slices"
	"strings"

	"main/global"
	"main/service"

	amqp "github.com/rabbitmq/amqp091-go"
)

var channels map[string]chan []byte = make(map[string]chan []byte)

func Message_Respond(d amqp.Delivery) {
	channel, exists := channels[d.CorrelationId]
	if exists {
		channel <- d.Body
	}
}

func Send_Message(w http.ResponseWriter, r *http.Request) {
	queue := strings.ToLower(r.PathValue("queue"))
	if !slices.Contains(global.Config.Queues, queue) {
		http.Redirect(w, r, "/App", http.StatusSeeOther)
		return
	}

	path := r.URL.Path[len(queue)+1:]
	corrId, err := service.Queue_Call(queue, []byte(path), "plain/text")
	if err != nil {
		http.NotFound(w, r)
		return
	}

	channel := make(chan []byte)
	channels[corrId] = channel
	body := <-channels[corrId]
	channels[corrId] = nil
	w.Write(body)
	w.Header().Set("content-type", "text/html")
	w.WriteHeader(http.StatusOK)
}
