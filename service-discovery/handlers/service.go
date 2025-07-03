package handlers

import (
	"encoding/json"
	"fmt"
	"io"
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

type Request struct {
	Url    string
	Method string
	Body   []byte
	Header http.Header
}

type Response struct {
	MimeType string
	Body     []byte
}

func Send_Message(w http.ResponseWriter, r *http.Request) {
	queue := strings.ToLower(r.PathValue("queue"))
	if !slices.Contains(global.Config.Queues, queue) {
		http.Redirect(w, r, "/App", http.StatusSeeOther)
		return
	}

	requestBody, _ := io.ReadAll(r.Body)
	proxy := Request{
		Url:    r.URL.Path[len(queue)+1:],
		Method: r.Method,
		Body:   requestBody,
		Header: r.Header,
	}
	payload, err := json.Marshal(&proxy)

	if err != nil {
		fmt.Println(err)
		http.NotFound(w, r)
		return
	}
	corrId, err := service.Queue_Call(queue, payload, "text/json")
	if err != nil {
		fmt.Println(err)
		http.NotFound(w, r)
		return
	}

	channel := make(chan []byte)
	channels[corrId] = channel
	responseBody := <-channels[corrId]
	channels[corrId] = nil

	var response Response
	err = json.Unmarshal(responseBody, &response)
	if err != nil {
		fmt.Println(err)
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", response.MimeType)
	w.WriteHeader(http.StatusOK)
	w.Write(response.Body)
}

func Redirect_To_UI(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/App", http.StatusSeeOther)
}
