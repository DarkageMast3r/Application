package main

import (
	"ZorgTechCatalogus/pkg/auth"
	"fmt"
)

func main() {
	fmt.Println(auth.GenerateRandomKey())
}
