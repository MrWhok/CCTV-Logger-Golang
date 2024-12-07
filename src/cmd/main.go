package main

import (
	"log"

	"CCTV-Logger-Golang/src/app"
	"CCTV-Logger-Golang/src/db"
	"CCTV-Logger-Golang/src/internal/config"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Connect to the database
	_, err := db.ConnectDB(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Unable to connect to the database: %v", err)
	}

	// Initialize the app
	r := app.SetupRouter()

	// Start the server
	port := cfg.Port
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Unable to start the server: %v", err)
	}
}
