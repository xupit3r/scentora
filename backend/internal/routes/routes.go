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
	tokenRepo := repository.NewRefreshTokenRepository(db)
	invitationRepo := repository.NewInvitationRepository(db)
	accordRepo := repository.NewAccordRepository(db)
	tagRepo := repository.NewPredefinedTagRepository(db)

	// Initialize services
	authService := services.NewAuthService(userRepo, tokenRepo, invitationRepo, cfg)
	invitationService := services.NewInvitationService(invitationRepo)
	accordService := services.NewAccordService(accordRepo, tagRepo)
	tagService := services.NewTagService(tagRepo)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)
	invitationHandler := handlers.NewInvitationHandler(invitationService)
	accordHandler := handlers.NewAccordHandler(accordService)
	tagHandler := handlers.NewTagHandler(tagService)

	// API group
	api := e.Group("/api")
	api.Use(middleware.GeneralRateLimiter())

	// Auth routes (public)
	auth := api.Group("/auth")
	auth.Use(middleware.AuthRateLimiter())
	auth.POST("/register", authHandler.Register)
	auth.POST("/login", authHandler.Login)
	auth.POST("/refresh", authHandler.Refresh)
	auth.POST("/logout", authHandler.Logout)

	// Protected auth routes
	authProtected := auth.Group("")
	authProtected.Use(middleware.JWTAuth(cfg))
	authProtected.GET("/me", authHandler.Me)
	authProtected.POST("/logout-all", authHandler.LogoutAll)

	// Invitation routes (protected)
	invitations := api.Group("/invitations")
	invitations.Use(middleware.JWTAuth(cfg))
	invitations.POST("", invitationHandler.Create)
	invitations.GET("", invitationHandler.List)
	invitations.DELETE("/:code", invitationHandler.Revoke)

	// Accord routes (protected)
	accords := api.Group("/accords")
	accords.Use(middleware.JWTAuth(cfg))
	accords.POST("", accordHandler.Create)
	accords.GET("", accordHandler.List)
	accords.GET("/:id", accordHandler.Get)
	accords.PUT("/:id", accordHandler.Update)
	accords.DELETE("/:id", accordHandler.Delete)
	accords.POST("/:id/tags", accordHandler.AddTag)
	accords.DELETE("/:id/tags/:tag", accordHandler.RemoveTag)

	// Tag routes (public - no auth required to view tags)
	tags := api.Group("/tags")
	tags.GET("", tagHandler.GetAll)
	tags.GET("/search", tagHandler.Search)
	tags.GET("/categories", tagHandler.GetCategories)
	tags.GET("/grouped", tagHandler.GetGrouped)
	tags.GET("/category/:category", tagHandler.GetByCategory)

	// TODO: Add stats/export routes in Phase 8.4
	// stats := api.Group("/stats")
	// stats.Use(middleware.JWTAuth(cfg))
	
	// export := api.Group("/export")
	// export.Use(middleware.JWTAuth(cfg))
}
