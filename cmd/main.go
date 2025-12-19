package main

import (
	"log"

	"github.com/joho/godotenv"

	"quocbui.dev/m/internal/app"
	"quocbui.dev/m/internal/config"
)

// @title URL Shortener API
// @version 1.0
// @BasePath /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Load config
	cfg := config.Load()

	// Initialize app with all dependencies
	application, err := app.New(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize app: %v", err)
	}

	// Run server
	if err := application.Run(); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
