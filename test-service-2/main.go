package main

import (
	"fmt"
	"log"
	"test_service_2/service"
)

func main() {
	service.Init()

	message, err := service.CallGet("Test")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(message)
}
