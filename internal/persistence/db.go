package persistence

import (
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var DB *gorm.DB

func ConnectToDb() {
	pgHost := os.Getenv("PGHOST")
	pgPort := os.Getenv("PGPORT")
	pgDatabase := os.Getenv("PGDATABASE")
	pgUser := os.Getenv("PGUSER")
	pgPassword := os.Getenv("PGPASSWORD")

	connectionString := fmt.Sprintf("host=%s port=%s user=%s dbname=%s", pgHost, pgPort, pgUser, pgDatabase)
	log.Debugf("Connecting to database: %+v", connectionString)

	connectionString += " password=" + pgPassword
	var err error
	DB, err = gorm.Open(postgres.Open(connectionString), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "cockaigne.", // schema name
			SingularTable: false,
		}})

	if err != nil {
		log.Fatal("Can't open database connection", err)
	}

	log.Debug("Database connection opened successfully")
}
