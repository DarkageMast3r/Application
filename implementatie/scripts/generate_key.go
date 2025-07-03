package main

import (
	"fmt"
	"ZorgTechImplementatie/pkg/auth"
)

func main() {
	fmt.Println(auth.GenerateRandomKey())
}
