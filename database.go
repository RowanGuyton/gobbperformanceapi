package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

// InitDatabase initializes the database connection and performs migrations
func InitDatabase() {
	// Load .env file
	enverr := godotenv.Load()
	if enverr != nil {
		log.Fatalf("Error loading values from .env: %v", enverr)
	}

	username := os.Getenv("DB_USERNAME")
	password := os.Getenv("DB_PASSWORD")

	// Database connection
	dsn := fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/exercise_db?charset=utf8mb4&parseTime=True&loc=Local", username, password)
	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto migrate the schema
	db.AutoMigrate(&Exercise{})
	db.AutoMigrate(&Meal{})
	db.AutoMigrate(&Weight{})
}
