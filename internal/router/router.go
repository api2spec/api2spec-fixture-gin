package router

import (
	"github.com/gin-gonic/gin"
	"github.com/api2spec/api2spec-fixture-gin/internal/handlers"
	"github.com/api2spec/api2spec-fixture-gin/internal/store"
)

// Setup creates and configures the Gin router with all routes
func Setup() *gin.Engine {
	r := gin.Default()

	// Initialize store
	memStore := store.NewMemoryStore()

	// Initialize handlers
	teapotHandler := handlers.NewTeapotHandler(memStore)
	teaHandler := handlers.NewTeaHandler(memStore)
	brewHandler := handlers.NewBrewHandler(memStore)
	healthHandler := handlers.NewHealthHandler()

	// Health routes
	r.GET("/health", healthHandler.Health)
	r.GET("/health/live", healthHandler.Live)
	r.GET("/health/ready", healthHandler.Ready)
	r.GET("/brew", healthHandler.Brew)

	// Teapot routes
	teapots := r.Group("/teapots")
	{
		teapots.GET("", teapotHandler.List)
		teapots.POST("", teapotHandler.Create)
		teapots.GET("/:id", teapotHandler.Get)
		teapots.PUT("/:id", teapotHandler.Update)
		teapots.PATCH("/:id", teapotHandler.Patch)
		teapots.DELETE("/:id", teapotHandler.Delete)
		teapots.GET("/:id/brews", brewHandler.ListByTeapot)
	}

	// Tea routes
	teas := r.Group("/teas")
	{
		teas.GET("", teaHandler.List)
		teas.POST("", teaHandler.Create)
		teas.GET("/:id", teaHandler.Get)
		teas.PUT("/:id", teaHandler.Update)
		teas.PATCH("/:id", teaHandler.Patch)
		teas.DELETE("/:id", teaHandler.Delete)
	}

	// Brew routes
	brews := r.Group("/brews")
	{
		brews.GET("", brewHandler.List)
		brews.POST("", brewHandler.Create)
		brews.GET("/:id", brewHandler.Get)
		brews.PATCH("/:id", brewHandler.Patch)
		brews.DELETE("/:id", brewHandler.Delete)
		brews.GET("/:id/steeps", brewHandler.ListSteeps)
		brews.POST("/:id/steeps", brewHandler.CreateSteep)
	}

	return r
}

// SetupWithStore creates and configures the Gin router with a provided store (for testing)
func SetupWithStore(memStore *store.MemoryStore) *gin.Engine {
	r := gin.Default()

	// Initialize handlers
	teapotHandler := handlers.NewTeapotHandler(memStore)
	teaHandler := handlers.NewTeaHandler(memStore)
	brewHandler := handlers.NewBrewHandler(memStore)
	healthHandler := handlers.NewHealthHandler()

	// Health routes
	r.GET("/health", healthHandler.Health)
	r.GET("/health/live", healthHandler.Live)
	r.GET("/health/ready", healthHandler.Ready)
	r.GET("/brew", healthHandler.Brew)

	// Teapot routes
	teapots := r.Group("/teapots")
	{
		teapots.GET("", teapotHandler.List)
		teapots.POST("", teapotHandler.Create)
		teapots.GET("/:id", teapotHandler.Get)
		teapots.PUT("/:id", teapotHandler.Update)
		teapots.PATCH("/:id", teapotHandler.Patch)
		teapots.DELETE("/:id", teapotHandler.Delete)
		teapots.GET("/:id/brews", brewHandler.ListByTeapot)
	}

	// Tea routes
	teas := r.Group("/teas")
	{
		teas.GET("", teaHandler.List)
		teas.POST("", teaHandler.Create)
		teas.GET("/:id", teaHandler.Get)
		teas.PUT("/:id", teaHandler.Update)
		teas.PATCH("/:id", teaHandler.Patch)
		teas.DELETE("/:id", teaHandler.Delete)
	}

	// Brew routes
	brews := r.Group("/brews")
	{
		brews.GET("", brewHandler.List)
		brews.POST("", brewHandler.Create)
		brews.GET("/:id", brewHandler.Get)
		brews.PATCH("/:id", brewHandler.Patch)
		brews.DELETE("/:id", brewHandler.Delete)
		brews.GET("/:id/steeps", brewHandler.ListSteeps)
		brews.POST("/:id/steeps", brewHandler.CreateSteep)
	}

	return r
}
