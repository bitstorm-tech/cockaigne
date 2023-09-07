package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
)

var Conn *sql.DB

func Init() {
	pgHost := os.Getenv("PG_HOST")
	pgPort := os.Getenv("PG_PORT")
	pgDatabase := os.Getenv("PG_DATABASE")
	pgUser := os.Getenv("PG_USER")
	pgPassword := os.Getenv("PG_PASSWORD")

	connectionString := fmt.Sprintf("host=%s port=%s user=%s dbname=%s", pgHost, pgPort, pgUser, pgDatabase)
	log.Print("Connecting to database: ", connectionString)

	connectionString += " password=" + pgPassword
	var err error
	Conn, err = sql.Open("postgres", connectionString)

	if err != nil {
		log.Fatal("Can't open database connection", err)
	}

	log.Println("Database connection opened successfully")

	err = Conn.Ping()

	if err != nil {
		log.Fatal("Can't ping database", err)
	}

	log.Println("Database pinged successfully")
}
