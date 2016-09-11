package main

import (
	"log"

	_ "github.com/Pholey/bitAPI/db"
	"github.com/Pholey/bitAPI/resources"
	// config "github.com/Pholey/bitAPI/config"
)

func main() {
	// Start the server
	println("Listening at :8080")
	log.Fatal(resources.NewRouter().Run(":8080"))
}
