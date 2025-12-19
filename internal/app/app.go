package app

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"

	"quocbui.dev/m/docs"
	"quocbui.dev/m/internal/config"
	"quocbui.dev/m/internal/handlers"
	"quocbui.dev/m/internal/middleware"
	"quocbui.dev/m/internal/repository"
	"quocbui.dev/m/internal/repository/postgres"
	"quocbui.dev/m/internal/service"
)

type App struct {
	Config *config.Config
	DB     *gorm.DB
	Router *gin.Engine
	Server *http.Server

	UserRepo  repository.UserRepository
	LinkRepo  repository.LinkRepository
	ClickRepo repository.ClickRepository
	TxManager repository.TransactionManager

	AuthService      *service.AuthService
	LinkService      *service.LinkService
	AnalyticsService *service.AnalyticsService
	GeoIPService     *service.GeoIPService

	AuthHandler *handlers.AuthHandler
	UserHandler *handlers.UserHandler
	LinkHandler *handlers.LinkHandler
}

func New(cfg *config.Config) (*App, error) {
	app := &App{Config: cfg}

	if err := app.initDB(); err != nil {
		return nil, err
	}

	app.initRepositories()
	app.initServices()
	app.initHandlers()
	app.initRouter()
	app.initServer()

	return app, nil
}

func (a *App) initDB() error {
	db, err := postgres.NewDB(&a.Config.DB)
	if err != nil {
		return err
	}
	a.DB = db
	return postgres.AutoMigrate(db)
}

func (a *App) initRepositories() {
	a.UserRepo = postgres.NewUserRepository(a.DB)
	a.LinkRepo = postgres.NewLinkRepository(a.DB)
	a.ClickRepo = postgres.NewClickRepository(a.DB)
	a.TxManager = postgres.NewTransactionManager(a.DB)
}

func (a *App) initServices() {
	a.GeoIPService = service.NewGeoIPService()
	a.AuthService = service.NewAuthService(a.UserRepo)
	a.LinkService = service.NewLinkService(a.LinkRepo, a.ClickRepo, a.TxManager, a.GeoIPService)
	a.AnalyticsService = service.NewAnalyticsService(a.ClickRepo, a.LinkRepo)
}

func (a *App) initHandlers() {
	a.AuthHandler = handlers.NewAuthHandler(a.AuthService, a.Config.JWT.Secret, a.Config.JWT.ExpiryHours)
	a.UserHandler = handlers.NewUserHandler(a.UserRepo)
	a.LinkHandler = handlers.NewLinkHandler(
		a.LinkService, a.AnalyticsService, a.AuthService,
		a.Config.App.Domain, a.Config.ShortCode.Length,
		a.Config.JWT.Secret, a.Config.JWT.ExpiryHours,
	)
}

func (a *App) initRouter() {
	if a.Config.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()
	r.Use(middleware.CORSMiddleware())
	r.Use(middleware.RateLimitMiddleware(a.Config.RateLimit.Requests, a.Config.RateLimit.Window))

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	docs.SwaggerInfo.Host = a.Config.App.Domain
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	a.registerRoutes(r)
	a.Router = r
}

func (a *App) registerRoutes(r *gin.Engine) {
	api := r.Group("/api/v1")
	{
		auth := api.Group("/auth")
		auth.POST("/register", a.AuthHandler.Register)
		auth.POST("/login", a.AuthHandler.Login)

		api.POST("/shorten", a.LinkHandler.Shorten)
	}

	protected := r.Group("/api/v1/me")
	protected.Use(middleware.AuthMiddleware(a.Config.JWT.Secret))
	{
		protected.GET("", a.UserHandler.GetMe)
		protected.GET("/links", a.LinkHandler.GetMyLinks)
		protected.GET("/links/:code", a.LinkHandler.GetMyLinkDetail)
		protected.DELETE("/links/:code", a.LinkHandler.DeleteMyLink)
	}

	r.GET("/:code", a.LinkHandler.Redirect)
}

func (a *App) initServer() {
	a.Server = &http.Server{
		Addr:         a.Config.App.Host + ":" + a.Config.App.Port,
		Handler:      a.Router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
}

func (a *App) Run() error {
	go func() {
		log.Printf("Server starting on %s:%s", a.Config.App.Host, a.Config.App.Port)
		if err := a.Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down...")
	return a.Shutdown()
}

func (a *App) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := a.Server.Shutdown(ctx); err != nil {
		return err
	}

	if sqlDB, err := a.DB.DB(); err == nil {
		sqlDB.Close()
	}

	log.Println("Server stopped")
	return nil
}
