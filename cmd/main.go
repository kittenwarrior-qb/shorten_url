package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"quocbui.dev/m/internal/config"
	"quocbui.dev/m/internal/repository/postgres"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Load config
	cfg := config.Load()

	// Connect to database
	db, err := postgres.NewDB(&cfg.DB)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto migrate
	if err := postgres.AutoMigrate(db); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	// Get underlying SQL DB for defer close
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("Failed to get underlying DB:", err)
	}
	defer sqlDB.Close()

	// Create Gin router
	r := gin.Default()

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "URL Shortener Service is running",
		})
	})

	log.Printf("Server starting on %s:%s...", cfg.App.Host, cfg.App.Port)
	if err := r.Run(cfg.App.Host + ":" + cfg.App.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
