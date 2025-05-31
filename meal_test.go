package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
