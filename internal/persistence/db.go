package persistence

import (
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2/log"
	"github.com/jmoiron/sqlx"
)

var pgPort = os.Getenv("PG_PORT")
var pgHost = os.Getenv("PG_HOST")
var pgDatabase = os.Getenv("PG_DATABASE")
var pgUser = os.Getenv("PG_USER")
var pgPassword = os.Getenv("PG_PASSWORD")
var DB *sqlx.DB

func ConnectToDb() {
	log.Debugf("Connecting to database: %s:*********@%s:%s/%s", pgUser, pgHost, pgPort, pgDatabase)
	connectionString := fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		pgHost,
		pgPort,
		pgUser,
		pgDatabase,
		pgPassword,
	)

	connectionString += " password=" + pgPassword

	var err error
	DB, err = sqlx.Connect("postgres", connectionString)
	if err != nil {
		log.Fatalf("Can't open database connection: %+v", err)
	}

	log.Info("Database connection opened successfully")
}
