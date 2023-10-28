package persistence

import (
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2/log"
	"github.com/jmoiron/sqlx"
)

var pgPort = os.Getenv("PGPORT")
var pgHost = os.Getenv("PGHOST")
var pgDatabase = os.Getenv("PGDATABASE")
var pgUser = os.Getenv("PGUSER")
var pgPassword = os.Getenv("PGPASSWORD")
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
