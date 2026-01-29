package routes

import (
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/yourusername/scentora-backend/internal/config"
	"github.com/yourusername/scentora-backend/internal/handlers"
	"github.com/yourusername/scentora-backend/internal/middleware"
	"github.com/yourusername/scentora-backend/internal/repository"
	"github.com/yourusername/scentora-backend/internal/services"
)

func SetupRoutes(e *echo.Echo, db *sqlx.DB, cfg *config.Config) {
	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	perfumeRepo := repository.NewPerfumeRepository(db)
	journalRepo := repository.NewJournalRepository(db)
	tokenRepo := repository.NewRefreshTokenRepository(db)

	// Initialize services
	authService := services.NewAuthService(userRepo, tokenRepo, cfg)
	perfumeService := services.NewPerfumeService(perfumeRepo)
	journalService := services.NewJournalService(journalRepo, perfumeRepo)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)
	perfumeHandler := handlers.NewPerfumeHandler(perfumeService)
	journalHandler := handlers.NewJournalHandler(journalService)
	notesHandler := handlers.NewNotesHandler(perfumeRepo)
	statsHandler := handlers.NewStatsHandler(perfumeRepo, journalRepo)
	exportHandler := handlers.NewExportHandler(perfumeRepo, journalRepo)

	// API group
	api := e.Group("/api")

	// Auth routes (public)
	auth := api.Group("/auth")
	auth.POST("/register", authHandler.Register)
	auth.POST("/login", authHandler.Login)
	auth.POST("/refresh", authHandler.Refresh)
	auth.POST("/logout", authHandler.Logout)

	// Protected auth routes
	authProtected := auth.Group("")
	authProtected.Use(middleware.JWTAuth(cfg))
	authProtected.GET("/me", authHandler.Me)

	// Perfume routes (protected)
	perfumes := api.Group("/perfumes")
	perfumes.Use(middleware.JWTAuth(cfg))
	perfumes.GET("", perfumeHandler.List)
	perfumes.GET("/:id", perfumeHandler.Get)
	perfumes.POST("", perfumeHandler.Create)
	perfumes.PUT("/:id", perfumeHandler.Update)
	perfumes.DELETE("/:id", perfumeHandler.Delete)

	// Journal routes (protected)
	perfumes.GET("/:perfumeId/journal", journalHandler.ListByPerfume)
	perfumes.POST("/:perfumeId/journal", journalHandler.Create)

	journal := api.Group("/journal")
	journal.Use(middleware.JWTAuth(cfg))
	journal.PUT("/:id", journalHandler.Update)
	journal.DELETE("/:id", journalHandler.Delete)

	// Other routes (protected)
	notes := api.Group("/notes")
	notes.Use(middleware.JWTAuth(cfg))
	notes.GET("", notesHandler.GetAll)

	stats := api.Group("/stats")
	stats.Use(middleware.JWTAuth(cfg))
	stats.GET("", statsHandler.Get)

	export := api.Group("/export")
	export.Use(middleware.JWTAuth(cfg))
	export.GET("", exportHandler.Export)
}
