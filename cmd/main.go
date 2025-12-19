package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"quocbui.dev/m/docs"
	"quocbui.dev/m/internal/config"
	"quocbui.dev/m/internal/handlers"
	"quocbui.dev/m/internal/repository/postgres"
	"quocbui.dev/m/internal/router"
	"quocbui.dev/m/internal/service"
)

// @title URL Shortener API
// @version 1.0
// @BasePath /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

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

	// Initialize repositories
	userRepo := postgres.NewUserRepository(db)
	linkRepo := postgres.NewLinkRepository(db)
	clickRepo := postgres.NewClickRepository(db)

	// Initialize transaction manager
	txManager := postgres.NewTransactionManager(db)

	// Initialize services
	geoIPService := service.NewGeoIPService()
	authService := service.NewAuthService(userRepo)
	linkService := service.NewLinkService(linkRepo, clickRepo, txManager, geoIPService)
	analyticsService := service.NewAnalyticsService(clickRepo, linkRepo)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService, cfg.JWT.Secret, cfg.JWT.ExpiryHours)
	userHandler := handlers.NewUserHandler(userRepo)
	linkHandler := handlers.NewLinkHandler(linkService, analyticsService, authService, cfg.App.Domain, cfg.ShortCode.Length, cfg.JWT.Secret, cfg.JWT.ExpiryHours)

	// Setup router
	r := router.SetupRouter(cfg, authHandler, userHandler, linkHandler)

	// Swagger documentation - set host dynamically
	docs.SwaggerInfo.Host = cfg.App.Domain
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Create server
	srv := &http.Server{
		Addr:    cfg.App.Host + ":" + cfg.App.Port,
		Handler: r,
	}

	// Graceful shutdown
	go func() {
		log.Printf("Server starting on %s:%s...", cfg.App.Host, cfg.App.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exited gracefully")
}
