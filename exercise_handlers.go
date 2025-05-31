package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// createExercise handles POST /exercises
func createExercise(c *gin.Context) {
	var exercise Exercise

	log.Println("Received request to create exercise")

	switch {
	case c.ShouldBindJSON(&exercise) != nil:
		log.Println("JSON Bind Error")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	case exercise.Sets <= 0 || exercise.Reps <= 0 || exercise.Weight < 0:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input values"})
		return
	case db.Create(&exercise).Error != nil:
		log.Println("DB Insert Error")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create exercise"})
		return
	default:
		log.Printf("Parsed Data: %+v\n", exercise)
		c.JSON(http.StatusCreated, exercise)
	}
}

// getExercises handles GET /exercises
func getExercises(c *gin.Context) {
	var exercises []Exercise
	switch err := db.Order("created_at DESC").Find(&exercises).Error; err {
	case nil:
		c.JSON(http.StatusOK, exercises)
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch exercises"})
	}
}

// getExercise handles GET /exercises/:id
func getExercise(c *gin.Context) {
	id := c.Param("id")
	var exercise Exercise
	switch err := db.First(&exercise, id).Error; err {
	case nil:
		c.JSON(http.StatusOK, exercise)
	default:
		c.JSON(http.StatusNotFound, gin.H{"error": "Exercise not found"})
	}
}

// updateExercise handles PUT /exercises/:id
func updateExercise(c *gin.Context) {
	id := c.Param("id")
	var exercise Exercise
	
	switch {
	case db.First(&exercise, id).Error != nil:
		c.JSON(http.StatusNotFound, gin.H{"error": "Exercise not found"})
		return
	case c.ShouldBindJSON(&exercise) != nil:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	case db.Save(&exercise).Error != nil:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update exercise"})
		return
	default:
		c.JSON(http.StatusOK, exercise)
	}
}

// deleteExercise handles DELETE /exercises/:id
func deleteExercise(c *gin.Context) {
	log.Println("Received request to delete exercise")
	id := c.Param("id")
	switch err := db.Delete(&Exercise{}, id).Error; err {
	case nil:
		c.JSON(http.StatusOK, gin.H{"message": "Exercise deleted successfully"})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete exercise"})
	}
}
