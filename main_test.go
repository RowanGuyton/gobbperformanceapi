package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// Helper function to measure request-response timing
// now takes a gin.HandlerFunc, returns a gin.HandlerFunc
func measureResponseTime(h gin.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		h(c)
		log.Printf("Request to %s took %s\n", c.Request.URL.Path, time.Since(start))
	}
}

func TestCreateExercise(t *testing.T) {
	reqBody, _ := json.Marshal(Exercise{Movement: "Push-ups", Reps: 10, Sets: 3, Weight: 0})
	req, _ := http.NewRequest("POST", "/exercises", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	// spin up Gin, register your endpoint
	r := gin.Default()
	r.POST("/exercises", measureResponseTime(createExercise))

	// response := executeRequest(measureResponseTime(createExercise), req)

	checkResponseCode(t, http.StatusCreated, response.Code)

	var m Exercise
	json.Unmarshal(response.Body.Bytes(), &m)
	assert.NotEmpty(t, m.ID)
	assert.Equal(t, "Push-ups", m.Movement)
}

func TestGetExercises(t *testing.T) {
	req, _ := http.NewRequest("GET", "/exercises", nil)
	response := executeRequest(measureResponseTime(getExercises), req)

	checkResponseCode(t, http.StatusOK, response.Code)
}

func TestUpdateExercise(t *testing.T) {
	reqBody, _ := json.Marshal(Exercise{Name: "Sit-ups", Reps: 15, Sets: 3, Weight: 0})
	req, _ := http.NewRequest("PUT", "/exercises/1", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	response := executeRequest(measureResponseTime(updateExercise), req)

	checkResponseCode(t, http.StatusOK, response.Code)
}

func TestDeleteExercise(t *testing.T) {
	req, _ := http.NewRequest("DELETE", "/exercises/1", nil)
	response := executeRequest(measureResponseTime(deleteExercise), req)

	checkResponseCode(t, http.StatusOK, response.Code)
}

func TestCreateMeal(t *testing.T) {
	reqBody, _ := json.Marshal(Meal{Name: "Breakfast", Calories: 500, Protein: 30, Carbs: 60, Fat: 20, Timestamp: time.Now().Unix()})
	req, _ := http.NewRequest("POST", "/meals", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	response := executeRequest(measureResponseTime(createMeal), req)

	checkResponseCode(t, http.StatusCreated, response.Code)

	var m Meal
	json.Unmarshal(response.Body.Bytes(), &m)
	assert.NotEmpty(t, m.ID)
	assert.Equal(t, "Breakfast", m.Name)
}

func TestGetMeals(t *testing.T) {
	req, _ := http.NewRequest("GET", "/meals", nil)
	response := executeRequest(measureResponseTime(getMeals), req)

	checkResponseCode(t, http.StatusOK, response.Code)
}

func TestUpdateMeal(t *testing.T) {
	reqBody, _ := json.Marshal(Meal{Name: "Lunch", Calories: 700, Protein: 40, Carbs: 80, Fat: 25, Timestamp: time.Now().Unix()})
	req, _ := http.NewRequest("PUT", "/meals/1", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	response := executeRequest(measureResponseTime(updateMeal), req)

	checkResponseCode(t, http.StatusOK, response.Code)
}

func TestDeleteMeal(t *testing.T) {
	req, _ := http.NewRequest("DELETE", "/meals/1", nil)
	response := executeRequest(measureResponseTime(deleteMeal), req)

	checkResponseCode(t, http.StatusOK, response.Code)
}

// Helper function to execute the request and return the response recorder
func executeRequest(handler http.HandlerFunc, req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	return rr
}

// Helper function to check the response code
func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

// Additional test cases for error scenarios can be added similarly
