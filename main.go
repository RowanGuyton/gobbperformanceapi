package main

import (
	"log"
)

func main() {
	// Initialize database connection
	InitDatabase()

	// Setup routes
	r := SetupRoutes()

	// Start server
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
