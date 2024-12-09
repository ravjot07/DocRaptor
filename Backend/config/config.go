package config

import (
	"fmt"
	"log"
	"os"

	"github.com/ravjot07/docraptor-backend/models"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var Db *gorm.DB

func InitDB() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, proceding with environment variables")
	}

	// Get environment variables
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	dbname := os.Getenv("DB_NAME")
	password := os.Getenv("DB_PASSWORD")

	//construct DNS
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		host, user, password, dbname, port)

	//connect to database
	DB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	log.Println("Connected to database")

	// auto-migrate the doc model

	err = DB.AutoMigrate(&models.Doc{})
	if err != nil {
		log.Fatalf("Error migrating the doc model: %v", err)
	}

	log.Println("Auto-migrated the doc model: database migration complete")

}
