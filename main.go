package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Struct for Movements for the ORM
type Exercise struct {
	gorm.Model
	ID       int     `json:"id"`
	Date     string  `json:"date"`
	Movement string  `json:"movement"`
	Sets     int     `json:"sets"`
	Reps     int     `json:"reps"`
	Weight   float64 `json:"weight"`
	Type     string  `json:"type"`
}

// Struct for Meals for the ORM
type Meal struct {
	gorm.Model
	ID       int    `json:"id"`
	Date     string `json:"date"`
	Name     string `json:"name"`
	Carbs    int    `json:"carbs"`
	Protein  int    `json:"protein"`
	Fats     int    `json:"fat"`
	Calories int    `json:"calories"`
}

// Struct for Weight entries for the ORM
type Weight struct {
	ID     int     `json:"id"`
	Date   string  `json:"date"`
	Weight float64 `json:"weight"`
}

var db *gorm.DB

func main() {

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

	// Start server
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

// ------------------
// Exercise functions
// ------------------
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

func getExercises(c *gin.Context) {
	var exercises []Exercise
	switch err := db.Order("created_at DESC").Find(&exercises).Error; err {
	case nil:
		c.JSON(http.StatusOK, exercises)
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch exercises"})
	}
}

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

// --------------
// Meal functions
// --------------
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

func getMeals(c *gin.Context) {
	var meals []Meal
	switch err := db.Order("created_at DESC").Find(&meals).Error; err {
	case nil:
		c.JSON(http.StatusOK, meals)
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch meals"})
	}
}

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

// ----------------
// Weight functions
// ----------------
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

func getWeightEntries(c *gin.Context) {
	var weights []Weight
	switch err := db.Order("created_at DESC").Find(&weights).Error; err {
	case nil:
		c.JSON(http.StatusOK, weights)
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch weight entries"})
	}
}

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
