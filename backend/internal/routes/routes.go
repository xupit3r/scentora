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
	
	// Recipe repositories
	recipeRepo := repository.NewRecipeRepository(db)
	recipeVersionRepo := repository.NewRecipeVersionRepository(db)
	recipeIngredientRepo := repository.NewRecipeIngredientRepository(db)
	recipeNoteRepo := repository.NewRecipeNoteRepository(db)
	recipeTagRepo := repository.NewRecipeTagRepository(db)
	recipeCollectionRepo := repository.NewRecipeCollectionRepository(db)

	// Initialize services
	authService := services.NewAuthService(userRepo, tokenRepo, invitationRepo, cfg)
	invitationService := services.NewInvitationService(invitationRepo)
	accordService := services.NewAccordService(accordRepo, tagRepo)
	tagService := services.NewTagService(tagRepo)
	
	// Recipe services
	recipeService := services.NewRecipeService(recipeRepo, recipeVersionRepo, recipeIngredientRepo, recipeNoteRepo, recipeTagRepo, accordRepo)
	recipeVersionService := services.NewRecipeVersionService(recipeVersionRepo, recipeRepo, recipeIngredientRepo)
	recipeIngredientService := services.NewRecipeIngredientService(recipeIngredientRepo, recipeVersionRepo, recipeRepo, accordRepo, userRepo)
	recipeNoteService := services.NewRecipeNoteService(recipeNoteRepo, recipeRepo)
	recipeTagService := services.NewRecipeTagService(recipeTagRepo, recipeRepo)
	recipeCollectionService := services.NewRecipeCollectionService(recipeCollectionRepo, recipeRepo)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)
	invitationHandler := handlers.NewInvitationHandler(invitationService)
	accordHandler := handlers.NewAccordHandler(accordService)
	tagHandler := handlers.NewTagHandler(tagService)
	statsHandler := handlers.NewStatsHandler(accordService)
	exportHandler := handlers.NewExportHandler(accordService)
	
	// Recipe handlers
	recipeHandler := handlers.NewRecipeHandler(recipeService)
	recipeVersionHandler := handlers.NewRecipeVersionHandler(recipeVersionService)
	recipeIngredientHandler := handlers.NewRecipeIngredientHandler(recipeIngredientService)
	recipeNoteHandler := handlers.NewRecipeNoteHandler(recipeNoteService)
	recipeTagHandler := handlers.NewRecipeTagHandler(recipeTagService)
	recipeCollectionHandler := handlers.NewRecipeCollectionHandler(recipeCollectionService)

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

	// Statistics routes (protected)
	api.GET("/stats", statsHandler.GetStats, middleware.JWTAuth(cfg))
	
	// Export/Import routes (protected)
	exportRoutes := api.Group("/export")
	exportRoutes.Use(middleware.JWTAuth(cfg))
	exportRoutes.GET("", exportHandler.Export)
	exportRoutes.POST("/import", exportHandler.Import)
	
	// Recipe routes (protected)
	recipes := api.Group("/recipes")
	recipes.Use(middleware.JWTAuth(cfg))
	recipes.POST("", recipeHandler.CreateRecipe)
	recipes.GET("", recipeHandler.ListRecipes)
	recipes.GET("/search", recipeHandler.SearchRecipes)
	recipes.GET("/:id", recipeHandler.GetRecipe)
	recipes.PUT("/:id", recipeHandler.UpdateRecipe)
	recipes.DELETE("/:id", recipeHandler.DeleteRecipe)
	
	// Recipe version routes
	recipes.POST("/:id/versions", recipeVersionHandler.CreateVersion)
	recipes.GET("/:id/versions", recipeVersionHandler.ListVersions)
	recipes.GET("/:id/versions/:versionId", recipeVersionHandler.GetVersion)
	recipes.POST("/:id/versions/:versionNumber/activate", recipeVersionHandler.ActivateVersion)
	recipes.POST("/:id/versions/:versionId/duplicate", recipeVersionHandler.DuplicateVersion)
	
	// Recipe ingredient routes
	recipes.POST("/:id/versions/:versionId/ingredients", recipeIngredientHandler.AddIngredient)
	recipes.PUT("/:id/versions/:versionId/ingredients/:ingredientId", recipeIngredientHandler.UpdateIngredient)
	recipes.DELETE("/:id/versions/:versionId/ingredients/:ingredientId", recipeIngredientHandler.RemoveIngredient)
	
	// Recipe note routes
	recipes.POST("/:id/notes", recipeNoteHandler.CreateNote)
	recipes.GET("/:id/notes", recipeNoteHandler.ListNotes)
	recipes.PUT("/:id/notes/:noteId", recipeNoteHandler.UpdateNote)
	recipes.DELETE("/:id/notes/:noteId", recipeNoteHandler.DeleteNote)
	
	// Recipe tag routes
	recipes.POST("/:id/tags", recipeTagHandler.AddTag)
	recipes.DELETE("/:id/tags/:tag", recipeTagHandler.RemoveTag)
	api.GET("/recipes/tags/popular", recipeTagHandler.GetPopularTags, middleware.JWTAuth(cfg))
	
	// Collection routes (protected)
	collections := api.Group("/collections")
	collections.Use(middleware.JWTAuth(cfg))
	collections.POST("", recipeCollectionHandler.CreateCollection)
	collections.GET("", recipeCollectionHandler.ListCollections)
	collections.GET("/:id", recipeCollectionHandler.GetCollection)
	collections.PUT("/:id", recipeCollectionHandler.UpdateCollection)
	collections.DELETE("/:id", recipeCollectionHandler.DeleteCollection)
	collections.POST("/:id/recipes", recipeCollectionHandler.AddRecipeToCollection)
	collections.DELETE("/:id/recipes/:recipeId", recipeCollectionHandler.RemoveRecipeFromCollection)
}
