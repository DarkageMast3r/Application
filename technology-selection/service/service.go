package service

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
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
}

func readConfig(path string) Config {
	jsonFile, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
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

func get_rabbitmq_connection() (*amqp.Connection, error) {
	if rabbitmq_connection != nil {
		return rabbitmq_connection, nil
	}
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		return nil, err
	}
	rabbitmq_connection = conn
	return rabbitmq_connection, nil
}

func Queue_Write(name string, body []byte, contentType string) error {
	conn, err := get_rabbitmq_connection()
	if err != nil {
		return err
	}

	channel, err := conn.Channel()
	if err != nil {
		return err
	}
	defer channel.Close()
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

	log.Printf(" [x] Queue %s Sent %s\n", queue.Name, body)
	return err
}

func Queue_Listen(name string, handler func(amqp.Delivery)) error {
	conn, err := get_rabbitmq_connection()
	if err != nil {
		return err
	}

	channel, err := conn.Channel()
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

	log.Printf(" [x] Queue %s Listening\n", queue.Name)
	go func() {
		for d := range msgs {
			log.Println(string(d.Body))
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

func Register(name string) int {
	port := getLocalPort()
	Get(Route(Get_Uri(service_discovery), "Register", name, strconv.Itoa(port)))
	return port
}

func Route(host string, routeValues ...string) string {
	return fmt.Sprintf(
		"https://%s/%s",
		host,
		strings.Join(routeValues, "/"),
	)
}

func Get(route string) (string, error) {
	resp, err := http.Get(route)
	if err != nil {
		fmt.Println("Failure to GET: ", err)
		return "", err
	}
	body, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		fmt.Println("Failure to read all: ", err)
		return "", err
	}
	return string(body), nil
}

func Post(route string, contentType string, data io.Reader) (string, error) {
	resp, err := http.Post(route, contentType, data)
	if err != nil {
		fmt.Println("Failure to POST: ", err)
		return "", err
	}
	body, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		fmt.Println("Failure to read all: ", err)
		return "", err
	}
	return string(body), nil
}

func CallGet(service string, routeValues ...string) (string, error) {
	return Get(Route(Get_Uri(service), routeValues...))
}

func CallPost(contentType string, body io.Reader, service string, routeValues ...string) (string, error) {
	return Post(Route(Get_Uri(service), routeValues...), contentType, body)
}

func reload(name string) string {
	uri, err := Get(Route(Get_Uri(service_discovery), name))
	if err != nil {
		fmt.Println("Could not get uri for service ", name)
		log.Fatal(err)
	} else {
		services[name] = uri
	}
	return uri
}

func Get_Uri(name string) string {
	uri, exists := services[name]
	// Get service if not used before
	if !exists {
		uri = reload(name)
	}

	// Get new service is previous stopped
	_, err := Get(Route(uri))
	if err != nil {
		uri = reload(name)
	}

	return uri
}
