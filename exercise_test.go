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
