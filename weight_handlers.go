package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// createWeightEntry handles POST /weights
func createWeightEntry(c *gin.Context) {
	var weight Weight

	switch {
	case c.ShouldBindBodyWithJSON(&weight) != nil:
		log.Println("JSON Bind Error")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	case weight.Weight <= 0:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid weight input value"})
		return
	case db.Create(&weight).Error != nil:
		log.Println("DB Insert Error")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create weight entry"})
		return
	default:
		c.JSON(http.StatusCreated, weight)
	}
}

// getWeightEntries handles GET /weights
func getWeightEntries(c *gin.Context) {
	var weights []Weight
	switch err := db.Order("created_at DESC").Find(&weights).Error; err {
	case nil:
		c.JSON(http.StatusOK, weights)
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch weight entries"})
	}
}

// updateWeightEntry handles PUT /weights/:id
func updateWeightEntry(c *gin.Context) {
	id := c.Param("id")
	var weight Weight
	
	switch {
	case db.First(&weight, id).Error != nil:
		c.JSON(http.StatusNotFound, gin.H{"error": "Weight entry not found"})
		return
	case c.ShouldBindJSON(&weight) != nil:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	case db.Save(&weight).Error != nil:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update weight entry"})
		return
	default:
		c.JSON(http.StatusOK, weight)
	}
}

// deleteWeightEntry handles DELETE /weights/:id
func deleteWeightEntry(c *gin.Context) {
	id := c.Param("id")

	switch {
	case id == "" || id == "undefined":
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid weight entry ID"})
		return
	default:
		result := db.Where("id = ?", id).Delete(&Weight{})
		switch {
		case result.Error != nil:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete weight entry"})
		case result.RowsAffected == 0:
			c.JSON(http.StatusNotFound, gin.H{"error": "Weight entry not found"})
		default:
			c.JSON(http.StatusOK, gin.H{"message": "Weight entry deleted successfully"})
		}
	}
}
