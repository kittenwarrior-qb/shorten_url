package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "quocbui.dev/m/docs"
	"quocbui.dev/m/internal/config"
	"quocbui.dev/m/internal/repository/postgres"
)

// @title           URL Shortener API
// @version         1.0
// @description     API service for URL shortening with analytics
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.email  support@quocbui.dev

// @license.name  MIT
// @license.url   http://opensource.org/licenses/MIT

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

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

	r := gin.Default()

	// Swagger documentation
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

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
