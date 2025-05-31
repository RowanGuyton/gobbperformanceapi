package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// setupTestDB creates a fresh in-memory SQLite database for testing
func setupTestDB() *gorm.DB {
	// Use a unique database for each test to avoid data persistence
	dsn := ":memory:"
	testDB, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to test database:", err)
	}
	// Auto-migrate the schema for test models
	testDB.AutoMigrate(&Exercise{}, &Meal{}, &Weight{})
	return testDB
}

// setupRouter creates a test router with all routes configured
func setupRouter() *gin.Engine {
	r := gin.Default()
	// Routes for exercises
	r.POST("/exercises", createExercise)
	r.GET("/exercises", getExercises)
	r.PUT("/exercises/:id", updateExercise)
	r.DELETE("/exercises/:id", deleteExercise)

	// Routes for meals
	r.POST("/meals", createMeal)
	r.GET("/meals", getMeals)
	r.PUT("/meals/:id", updateMeal)
	r.DELETE("/meals/:id", deleteMeal)

	// Routes for weight entries
	r.POST("/weights", createWeightEntry)
	r.GET("/weights", getWeightEntries)
	r.PUT("/weights/:id", updateWeightEntry)
	r.DELETE("/weights/:id", deleteWeightEntry)

	return r
}
