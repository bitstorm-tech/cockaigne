package persistence

import (
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var pgPort = os.Getenv("PGPORT")
var pgHost = os.Getenv("PGHOST")
var pgDatabase = os.Getenv("PGDATABASE")
var pgUser = os.Getenv("PGUSER")
var pgPassword = os.Getenv("PGPASSWORD")
var ConnectionString = fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s", pgHost, pgPort, pgUser, pgDatabase, pgPassword)
var DB *gorm.DB

func ConnectToDb() {
	log.Debugf("Connecting to database: %s:*********@%s:%s/%s", pgUser, pgHost, pgPort, pgDatabase)

	ConnectionString += " password=" + pgPassword
	var err error
	DB, err = gorm.Open(postgres.Open(ConnectionString), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "cockaigne.",
			SingularTable: false,
		}})

	if err != nil {
		log.Fatal("Can't open database connection", err)
	}

	DB.Exec("set search_path = public,cockaigne")

	log.Debug("Database connection opened successfully")
}
