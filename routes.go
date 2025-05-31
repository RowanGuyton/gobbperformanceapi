package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// SetupRoutes configures all the routes for the application
func SetupRoutes() *gin.Engine {
	// Initialize Gin router
	r := gin.Default()
	r.Use(cors.Default())

	// Serve static files
	r.Static("/static", "./static")

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
