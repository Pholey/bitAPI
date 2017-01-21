package main

import (
	_ "github.com/Pholey/bitAPI/db"
	"github.com/Pholey/bitAPI/resources"
)

func main() {
	// Start the server
	println("Listening at :9000")
	resources.NewRouter().Start(":9000")
}
