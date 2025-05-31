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
