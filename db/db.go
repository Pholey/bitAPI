package db

import (
	"github.com/Pholey/bitAPI/config"
	"fmt"

	dbr "github.com/gocraft/dbr"
	_ "github.com/lib/pq"
)

type DatabaseConnection struct {
	user     string
	password string
	host     string
	database string
}

var Session *dbr.Session
var ConnectionInfo DatabaseConnection

func init() {
	info := config.File.GetStringMap("database")

	ConnectionInfo = DatabaseConnection{
		info["user"].(string),
		info["password"].(string),
		info["host"].(string),
		info["database"].(string),
	}

	connString := fmt.Sprintf(
		"postgres://%s:%s@%s/%s?sslmode=disable",
		info["user"],
		info["password"],
		info["host"],
		info["database"],
	)

	conn, err := dbr.Open("postgres", connString, nil)
	if err != nil {
		panic(err)
	}

	Session = conn.NewSession(nil)
	println("Connected to Postgres")
}
