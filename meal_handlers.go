package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// createMeal handles POST /meals
func createMeal(c *gin.Context) {
	var meal Meal

	log.Println("Received request to create meal")

	switch {
	case c.ShouldBindJSON(&meal) != nil:
		log.Println("JSON Bind Error")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	case meal.Carbs < 0 || meal.Fats < 0 || meal.Protein < 0 || meal.Calories < 0:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input values"})
		return
	case db.Create(&meal).Error != nil:
		log.Println("DB Insert Error")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create meal"})
		return
	default:
		log.Printf("Parsed Data: %+v\n", meal)
		c.JSON(http.StatusCreated, meal)
	}
}

// getMeals handles GET /meals
func getMeals(c *gin.Context) {
	var meals []Meal
	switch err := db.Order("created_at DESC").Find(&meals).Error; err {
	case nil:
		c.JSON(http.StatusOK, meals)
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch meals"})
	}
}

// updateMeal handles PUT /meals/:id
func updateMeal(c *gin.Context) {
	id := c.Param("id")
	var meal Meal
	
	switch {
	case db.First(&meal, id).Error != nil:
		c.JSON(http.StatusNotFound, gin.H{"error": "Meal not found"})
		return
	case c.ShouldBindJSON(&meal) != nil:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	case db.Save(&meal).Error != nil:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update meal"})
		return
	default:
		c.JSON(http.StatusOK, meal)
	}
}

// deleteMeal handles DELETE /meals/:id
func deleteMeal(c *gin.Context) {
	id := c.Param("id")

	switch {
	case id == "" || id == "undefined":
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid meal ID"})
		return
	default:
		result := db.Where("id = ?", id).Delete(&Meal{})
		switch {
		case result.Error != nil:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete meal"})
		case result.RowsAffected == 0:
			c.JSON(http.StatusNotFound, gin.H{"error": "Meal not found"})
		default:
			c.JSON(http.StatusOK, gin.H{"message": "Meal deleted successfully"})
		}
	}
}
