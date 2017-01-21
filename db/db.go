package db

import (
	"github.com/Pholey/bitAPI/config"
	_ "database/sql"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type DatabaseConnection struct {
	user     string
	password string
	host     string
	database string
}

var Client *sqlx.DB
var Sq sq.StatementBuilderType
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

	Sq = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	client, err := sqlx.Connect("postgres", connString)
	if err != nil {
		panic(err)
	}

	Client = client
	println("Connected to Postgres")
}
