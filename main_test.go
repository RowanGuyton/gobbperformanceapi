package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB() *gorm.DB {
	dsn := "file::memory:?cache=shared"
	testDB, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to test database:", err)
	}
	// Auto-migrate the schema for test models
	testDB.AutoMigrate(&Exercise{}, &Meal{}, &Weight{})
	return testDB
}

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

// ------------------
// Exercise Tests
// ------------------
func TestCreateExercise_Success(t *testing.T) {
	db = setupTestDB()

	r := setupRouter()

	reqBody, _ := json.Marshal(Exercise{
		Date:     "2023-10-01",
		Movement: "Push-ups",
		Reps:     10,
		Sets:     3,
		Weight:   0,
		Type:     "Bodyweight",
	})
	req, _ := http.NewRequest("POST", "/exercises", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var m Exercise
	err := json.Unmarshal(w.Body.Bytes(), &m)
	assert.NoError(t, err)
	assert.NotEmpty(t, m.ID)
	assert.Equal(t, "Push-ups", m.Movement)
}

func TestCreateExercise_InvalidInput(t *testing.T) {
	db = setupTestDB()

	r := setupRouter()

	reqBody, _ := json.Marshal(Exercise{
		Date:     "2023-10-01",
		Movement: "Push-ups",
		Reps:     0, // Invalid reps
		Sets:     3,
		Weight:   0,
		Type:     "Bodyweight",
	})
	req, _ := http.NewRequest("POST", "/exercises", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Invalid input values", response["error"])
}

func TestGetExercises_Success(t *testing.T) {
	db = setupTestDB()

	// Add test data
	exercise := Exercise{
		Date:     "2023-10-01",
		Movement: "Push-ups",
		Reps:     10,
		Sets:     3,
		Weight:   0,
		Type:     "Bodyweight",
	}
	db.Create(&exercise)

	r := setupRouter()

	req, _ := http.NewRequest("GET", "/exercises", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var exercises []Exercise
	err := json.Unmarshal(w.Body.Bytes(), &exercises)
	assert.NoError(t, err)
	assert.Len(t, exercises, 1)
	assert.Equal(t, "Push-ups", exercises[0].Movement)
}

func TestUpdateExercise_Success(t *testing.T) {
	db = setupTestDB()

	// Create a test exercise
	exercise := Exercise{
		Date:     "2023-10-01",
		Movement: "Push-ups",
		Reps:     10,
		Sets:     3,
		Weight:   0,
		Type:     "Bodyweight",
	}
	db.Create(&exercise)

	r := setupRouter()

	reqBody, _ := json.Marshal(Exercise{
		Date:     "2023-10-01",
		Movement: "Sit-ups",
		Reps:     15,
		Sets:     3,
		Weight:   0,
		Type:     "Bodyweight",
	})
	req, _ := http.NewRequest("PUT", fmt.Sprintf("/exercises/%d", exercise.ID), bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var updatedExercise Exercise
	err := json.Unmarshal(w.Body.Bytes(), &updatedExercise)
	assert.NoError(t, err)
	assert.Equal(t, "Sit-ups", updatedExercise.Movement)
}

func TestUpdateExercise_NotFound(t *testing.T) {
	db = setupTestDB()

	r := setupRouter()

	reqBody, _ := json.Marshal(Exercise{
		Date:     "2023-10-01",
		Movement: "Sit-ups",
		Reps:     15,
		Sets:     3,
		Weight:   0,
		Type:     "Bodyweight",
	})
	req, _ := http.NewRequest("PUT", "/exercises/999", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestDeleteExercise_Success(t *testing.T) {
	db = setupTestDB()

	// Create a test exercise
	exercise := Exercise{
		Date:     "2023-10-01",
		Movement: "Push-ups",
		Reps:     10,
		Sets:     3,
		Weight:   0,
		Type:     "Bodyweight",
	}
	db.Create(&exercise)

	r := setupRouter()

	req, _ := http.NewRequest("DELETE", fmt.Sprintf("/exercises/%d", exercise.ID), nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Exercise deleted successfully", response["message"])
}

// ------------------
// Meal Tests
// ------------------
func TestCreateMeal_Success(t *testing.T) {
	db = setupTestDB()

	r := setupRouter()

	reqBody, _ := json.Marshal(Meal{
		Date:     "2023-10-01",
		Name:     "Breakfast",
		Carbs:    60,
		Protein:  30,
		Fats:     20,
		Calories: 500,
	})
	req, _ := http.NewRequest("POST", "/meals", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var m Meal
	err := json.Unmarshal(w.Body.Bytes(), &m)
	assert.NoError(t, err)
	assert.NotEmpty(t, m.ID)
	assert.Equal(t, "Breakfast", m.Name)
}

func TestCreateMeal_InvalidInput(t *testing.T) {
	db = setupTestDB()

	r := setupRouter()

	reqBody, _ := json.Marshal(Meal{
		Date:     "2023-10-01",
		Name:     "Breakfast",
		Carbs:    -10, // Invalid carbs
		Protein:  30,
		Fats:     20,
		Calories: 500,
	})
	req, _ := http.NewRequest("POST", "/meals", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Invalid input values", response["error"])
}

func TestGetMeals_Success(t *testing.T) {
	db = setupTestDB()

	// Add test data
	meal := Meal{
		Date:     "2023-10-01",
		Name:     "Breakfast",
		Carbs:    60,
		Protein:  30,
		Fats:     20,
		Calories: 500,
	}
	db.Create(&meal)

	r := setupRouter()

	req, _ := http.NewRequest("GET", "/meals", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var meals []Meal
	err := json.Unmarshal(w.Body.Bytes(), &meals)
	assert.NoError(t, err)
	assert.Len(t, meals, 1)
	assert.Equal(t, "Breakfast", meals[0].Name)
}

func TestUpdateMeal_Success(t *testing.T) {
	db = setupTestDB()

	// Create a test meal
	meal := Meal{
		Date:     "2023-10-01",
		Name:     "Breakfast",
		Carbs:    60,
		Protein:  30,
		Fats:     20,
		Calories: 500,
	}
	db.Create(&meal)

	r := setupRouter()

	reqBody, _ := json.Marshal(Meal{
		Date:     "2023-10-01",
		Name:     "Lunch",
		Carbs:    80,
		Protein:  40,
		Fats:     25,
		Calories: 700,
	})
	req, _ := http.NewRequest("PUT", fmt.Sprintf("/meals/%d", meal.ID), bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var updatedMeal Meal
	err := json.Unmarshal(w.Body.Bytes(), &updatedMeal)
	assert.NoError(t, err)
	assert.Equal(t, "Lunch", updatedMeal.Name)
}

func TestDeleteMeal_Success(t *testing.T) {
	db = setupTestDB()

	// Create a test meal
	meal := Meal{
		Date:     "2023-10-01",
		Name:     "Breakfast",
		Carbs:    60,
		Protein:  30,
		Fats:     20,
		Calories: 500,
	}
	db.Create(&meal)

	r := setupRouter()

	req, _ := http.NewRequest("DELETE", fmt.Sprintf("/meals/%d", meal.ID), nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Meal deleted successfully", response["message"])
}

// ------------------
// Weight Tests
// ------------------
func TestCreateWeightEntry_Success(t *testing.T) {
	db = setupTestDB()

	r := setupRouter()

	reqBody, _ := json.Marshal(Weight{
		Date:   "2023-10-01",
		Weight: 75.5,
	})
	req, _ := http.NewRequest("POST", "/weights", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var wEntry Weight
	err := json.Unmarshal(w.Body.Bytes(), &wEntry)
	assert.NoError(t, err)
	assert.NotEmpty(t, wEntry.ID)
	assert.Equal(t, 75.5, wEntry.Weight)
}

func TestCreateWeightEntry_InvalidInput(t *testing.T) {
	db = setupTestDB()

	r := setupRouter()

	reqBody, _ := json.Marshal(Weight{
		Date:   "2023-10-01",
		Weight: 0, // Invalid weight
	})
	req, _ := http.NewRequest("POST", "/weights", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Invalid weight input value", response["error"])
}

func TestGetWeightEntries_Success(t *testing.T) {
	db = setupTestDB()

	// Add test data
	weight := Weight{
		Date:   "2023-10-01",
		Weight: 75.5,
	}
	db.Create(&weight)

	r := setupRouter()

	req, _ := http.NewRequest("GET", "/weights", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var weights []Weight
	err := json.Unmarshal(w.Body.Bytes(), &weights)
	assert.NoError(t, err)
	assert.Len(t, weights, 1)
	assert.Equal(t, 75.5, weights[0].Weight)
}

func TestUpdateWeightEntry_Success(t *testing.T) {
	db = setupTestDB()

	// Create a test weight entry
	weight := Weight{
		Date:   "2023-10-01",
		Weight: 75.5,
	}
	db.Create(&weight)

	r := setupRouter()

	reqBody, _ := json.Marshal(Weight{
		Date:   "2023-10-02",
		Weight: 76.0,
	})
	req, _ := http.NewRequest("PUT", fmt.Sprintf("/weights/%d", weight.ID), bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var updatedWeight Weight
	err := json.Unmarshal(w.Body.Bytes(), &updatedWeight)
	assert.NoError(t, err)
	assert.Equal(t, 76.0, updatedWeight.Weight)
}

func TestDeleteWeightEntry_Success(t *testing.T) {
	db = setupTestDB()

	// Create a test weight entry
	weight := Weight{
		Date:   "2023-10-01",
		Weight: 75.5,
	}
	db.Create(&weight)

	r := setupRouter()

	req, _ := http.NewRequest("DELETE", fmt.Sprintf("/weights/%d", weight.ID), nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Weight entry deleted successfully", response["message"])
}
