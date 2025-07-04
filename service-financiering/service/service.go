package service

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"io"
	"log"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Config struct {
	Service_discovery_root string `json:"service_discovery_root"`
	Service_discovery_port int    `json:"service_discovery_port"`
	Allow_insecure         bool   `json:"allow_insecure"`
	RabbitMQ               string `json:"rabbitmq"`
}

func readConfig(path string) Config {
	jsonFile, err := os.Open(path)
	if err != nil {
		LogError(err)
	}
	defer jsonFile.Close()

	byteValue, _ := io.ReadAll(jsonFile)
	var config Config
	json.Unmarshal(byteValue, &config)
	return config
}

func getLocalPort() int {
	conn, err := net.Dial("udp", "0.0.0.0:80")
	if err != nil {
		log.Fatal(err)
	}
	return conn.LocalAddr().(*net.UDPAddr).Port
}

var config Config
var services map[string]string

const service_discovery string = "service_discovery"

var rabbitmq_connection *amqp.Connection
var rabbitmq_channel *amqp.Channel

func get_rabbitmq_connection() (*amqp.Connection, error) {
	if rabbitmq_connection != nil {
		return rabbitmq_connection, nil
	}
	conn, err := amqp.Dial(config.RabbitMQ)
	if err != nil {
		return nil, err
	}
	rabbitmq_connection = conn
	return rabbitmq_connection, nil
}

func get_rabbitmq_channel() (*amqp.Channel, error) {
	if rabbitmq_channel != nil {
		return rabbitmq_channel, nil
	}
	conn, err := get_rabbitmq_connection()
	if err != nil {
		return nil, err
	}

	channel, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	rabbitmq_connection = conn
	return channel, nil
}

func Queue_Write(name string, body []byte, contentType string) error {
	channel, err := get_rabbitmq_channel()
	if err != nil {
		return err
	}
	queue, err := channel.QueueDeclare(
		name,  // name
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = channel.PublishWithContext(ctx,
		"",         // exchange
		queue.Name, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType: contentType,
			Body:        body,
		},
	)

	return err
}

func Queue_Respond(d amqp.Delivery, body []byte, contentType string) error {
	channel, err := get_rabbitmq_channel()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = channel.PublishWithContext(ctx,
		"",        // exchange
		d.ReplyTo, // routing key
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			ContentType:   contentType,
			Body:          body,
			CorrelationId: d.CorrelationId,
		},
	)

	return err
}

func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}

func randomString(l int) string {
	bytes := make([]byte, l)
	for i := 0; i < l; i++ {
		bytes[i] = byte(randInt(65, 90))
	}
	return string(bytes)
}

func Queue_Call(name string, body []byte, contentType string) (string, error) {
	channel, err := get_rabbitmq_channel()
	if err != nil {
		return "", err
	}

	queue, err := channel.QueueDeclare(
		name,  // name
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // noWait
		nil,   // arguments
	)
	if err != nil {
		return "", err
	}

	corrId := randomString(32)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = channel.PublishWithContext(ctx,
		"",         // exchange
		queue.Name, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType:   contentType,
			CorrelationId: corrId,
			ReplyTo:       "Result",
			Body:          body,
		},
	)
	return corrId, err
}

func Queue_Call_Request(name string, request Request) (string, error) {
	payload, err := json.Marshal(&request)
	if err != nil {
		return "", err
	}

	return Queue_Call(name, payload, "text/json")
}

func Queue_Listen(name string, handler func(amqp.Delivery)) error {
	channel, err := get_rabbitmq_channel()
	if err != nil {
		return err
	}

	queue, err := channel.QueueDeclare(
		name,  // name
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		channel.Close()
		return err
	}
	msgs, err := channel.Consume(
		queue.Name, // queue
		"",         // consumer
		true,       // auto-ack
		false,      // exclusive
		false,      // no-local
		false,      // no-wait
		nil,        // args
	)
	if err != nil {
		channel.Close()
		return err
	}

	go func() {
		for d := range msgs {
			handler(d)
		}
	}()
	return nil
}

func Init() {
	config = readConfig("config.json")
	if config.Allow_insecure {
		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}
	services = make(map[string]string)
	services[service_discovery] = config.Service_discovery_root + ":" + strconv.Itoa(config.Service_discovery_port)
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

func Register(name string, handler func(http.ResponseWriter, *http.Request)) int {
	err := Queue_Listen(name, func(d amqp.Delivery) {
		var messageRequest Request
		err := json.Unmarshal(d.Body, &messageRequest)
		if err != nil {
			LogError(err)
			return
		}
		writer := NewServiceResponseWriter()
		request := new(http.Request)
		request.Method = messageRequest.Method
		request.URL = new(url.URL)
		request.URL.Path = messageRequest.Url
		request.Body = io.NopCloser(strings.NewReader(string(messageRequest.Body)))
		request.Header = messageRequest.Header
		handler(writer, request)
		messageResponse := Response{
			MimeType: writer.Header().Get("content-type"),
			Body:     writer.body,
		}
		body, err := json.Marshal(&messageResponse)
		if err != nil {
			LogError(err)
			return
		}
		err = Queue_Respond(d, body, "text/html")
		if err != nil {
			LogError(err)
		}
	})
	if err != nil {
		log.Fatal(err)
	}
	return getLocalPort()
}
