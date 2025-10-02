package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"

	"memes-generator/internal/delivery/http"
	"memes-generator/internal/repository"
	"memes-generator/internal/usecase"
)

const (
	defaultPort    = "8080"
	defaultWebRoot = "./web/build"
)

func main() {
	// Get port from environment variable or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	// Get web root from environment variable or use default
	webRoot := os.Getenv("WEB_ROOT")
	if webRoot == "" {
		webRoot = defaultWebRoot
	}

	// Ensure data directories exist
	dataDir := "./data/memes"
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		log.Fatalf("Failed to create data directory: %v", err)
	}

	templatesDir := "./data/templates"
	if err := os.MkdirAll(templatesDir, 0755); err != nil {
		log.Fatalf("Failed to create templates directory: %v", err)
	}

	// Initialize repositories
	memeRepo := repository.NewMemeFileRepository()
	templateRepo := repository.NewTemplateFileRepository()

	// Initialize usecases
	memeUsecase := usecase.NewMemeUsecase(memeRepo)
	templateUsecase := usecase.NewTemplateUsecase(templateRepo)

	// Initialize handler
	memeHandler := http.NewMemeHandler(memeUsecase, templateUsecase, webRoot)

	// Initialize Gin router
	router := gin.Default()

	// Serve static files (React app)
	router.Static("/static", filepath.Join(webRoot, "static"))
	router.StaticFile("/favicon.ico", filepath.Join(webRoot, "favicon.ico"))
	router.StaticFile("/favicon.svg", filepath.Join(webRoot, "favicon.svg"))

	// API routes
	api := router.Group("/api")
	{
		api.GET("/memes", memeHandler.ListMemes)
		api.POST("/memes", memeHandler.CreateMeme)
		api.GET("/memes/:id", memeHandler.GetMeme)
		api.DELETE("/memes/:id", memeHandler.DeleteMeme)

		// Template routes
		api.GET("/templates", memeHandler.ListTemplates)
		api.POST("/templates", memeHandler.CreateTemplate)
		api.POST("/templates/:name/image", memeHandler.UploadTemplateImage)
	}

	// Image routes
	router.GET("/memes/:id/image", memeHandler.ServeMemeImage)
	router.GET("/templates/:name/image", memeHandler.ServeTemplateImage)

	// Serve React app for all other routes (SPA)
	router.NoRoute(func(c *gin.Context) {
		// If it's an API route, return 404
		if len(c.Request.URL.Path) >= 4 && c.Request.URL.Path[:4] == "/api" {
			c.JSON(404, gin.H{"error": "API endpoint not found"})
			return
		}
		// Otherwise serve the React app
		memeHandler.ServeWebApp(c)
	})

	// Start server
	log.Printf("Server starting on port %s", port)
	log.Printf("Serving React app from %s", webRoot)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
