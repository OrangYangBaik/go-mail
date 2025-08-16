package db

import (
	"fmt"
	"go-mail/constants"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitPostgres() {
	dbUser := os.Getenv(constants.DB_USER)
	dbPass := os.Getenv(constants.DB_PASS)
	dbHost := os.Getenv(constants.DB_HOST)
	dbPort := os.Getenv(constants.DB_PORT)
	dbName := os.Getenv(constants.DB_NAME)

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		dbHost, dbUser, dbPass, dbName, dbPort)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	var db_name string
	DB.Raw("SELECT current_database()").Scan(&db_name)
	log.Println("Connected to DB:", db_name)

	log.Println("Postgres connected with GORM!")
}
