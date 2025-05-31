package main

import "gorm.io/gorm"

// Exercise represents a workout exercise entry
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

// Meal represents a meal entry with nutritional information
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

// Weight represents a weight tracking entry
type Weight struct {
	gorm.Model
	ID     int     `json:"id"`
	Date   string  `json:"date"`
	Weight float64 `json:"weight"`
}
