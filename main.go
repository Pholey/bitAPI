package main

import (
	_ "github.com/Pholey/bitAPI/db"
	"github.com/Pholey/bitAPI/resources"
	// config "github.com/Pholey/bitAPI/config"
	"github.com/labstack/echo/engine/standard"
)

func main() {
	// Start the server
	println("Listening at :9000")
	resources.NewRouter().Run(standard.New(":9000"))
}
