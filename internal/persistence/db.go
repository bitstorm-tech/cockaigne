package persistence

import (
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2/log"
	"github.com/jmoiron/sqlx"
)

var DB *sqlx.DB

func ConnectToDb() {
	var pgPort = os.Getenv("PG_PORT")
	var pgHost = os.Getenv("PG_HOST")
	var pgDatabase = os.Getenv("PG_DATABASE")
	var pgUser = os.Getenv("PG_USER")
	var pgPassword = os.Getenv("PG_PASSWORD")
	var pgSslMode = os.Getenv("PG_SSL_MODE")
	var pgSchema = os.Getenv("PG_SCHEMA")

	connectionString := fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s sslmode=%s search_path=public,%s",
		pgHost,
		pgPort,
		pgUser,
		pgDatabase,
		pgSslMode,
		pgSchema,
	)

	log.Debugf("Connecting to database: %s", connectionString)

	connectionString += " password=" + pgPassword

	var err error
	DB, err = sqlx.Connect("postgres", connectionString)
	if err != nil {
		log.Fatalf("Can't open database connection: %+v", err)
	}

	log.Info("Database connection opened successfully")
}
