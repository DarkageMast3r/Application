package main

import (
	"fmt"
	"log"
	"main/global"
	"main/handlers"
	"main/service"
	"net/http"
	"strconv"
)

func main() {
	config := global.ReadConfig("config.json")
	var port_str = strconv.Itoa(config.Service_discovery_port)
	service.Init()

	service.Queue_Listen("Result", handlers.Message_Respond)
	http.HandleFunc("/{queue}/", handlers.Send_Message)

	fmt.Printf("Listening on %s\n", ":"+port_str)
	err := http.ListenAndServeTLS(":"+port_str, "../server.crt", "../server.key", nil)
	if err != nil {
		log.Println(err)
	}
}
