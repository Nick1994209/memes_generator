package http

import (
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"

	"memes-generator/internal/domain"
	"memes-generator/internal/usecase"
)

// MemeHandler represents the HTTP handler for memes
type MemeHandler struct {
	memeUsecase     domain.MemeUsecase
	templateUsecase *usecase.TemplateUsecase
	webRoot         string
}

// NewMemeHandler creates a new meme handler
func NewMemeHandler(memeUsecase domain.MemeUsecase, templateUsecase *usecase.TemplateUsecase, webRoot string) *MemeHandler {
	return &MemeHandler{
		memeUsecase:     memeUsecase,
		templateUsecase: templateUsecase,
		webRoot:         webRoot,
	}
}

// CreateMemeRequest represents the request body for creating a meme
type CreateMemeRequest struct {
	Template   string `json:"template" binding:"required"`
	TextTop    string `json:"text_top"`
	TextBottom string `json:"text_bottom"`
}

// MemeResponse represents the response body for a meme
type MemeResponse struct {
	ID         string `json:"id"`
	Template   string `json:"template"`
	TextTop    string `json:"text_top"`
	TextBottom string `json:"text_bottom"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}

// TemplateResponse represents the response body for a template
type TemplateResponse struct {
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// CreateMeme handles the creation of a new meme
func (h *MemeHandler) CreateMeme(c *gin.Context) {
	var req CreateMemeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	meme, err := h.memeUsecase.CreateMeme(req.Template, req.TextTop, req.TextBottom)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := MemeResponse{
		ID:         meme.ID,
		Template:   meme.Template,
		TextTop:    meme.TextTop,
		TextBottom: meme.TextBottom,
		CreatedAt:  meme.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:  meme.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}

	c.JSON(http.StatusCreated, response)
}

// GetMeme retrieves a specific meme by ID
func (h *MemeHandler) GetMeme(c *gin.Context) {
	id := c.Param("id")

	meme, err := h.memeUsecase.GetMemeByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Meme not found"})
		return
	}

	response := MemeResponse{
		ID:         meme.ID,
		Template:   meme.Template,
		TextTop:    meme.TextTop,
		TextBottom: meme.TextBottom,
		CreatedAt:  meme.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:  meme.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}

	c.JSON(http.StatusOK, response)
}

// ListMemes retrieves all memes
func (h *MemeHandler) ListMemes(c *gin.Context) {
	memes, err := h.memeUsecase.ListMemes()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Ensure we always return an array, even if empty
	if memes == nil {
		memes = []*domain.Meme{}
	}

	var response []MemeResponse
	for _, meme := range memes {
		response = append(response, MemeResponse{
			ID:         meme.ID,
			Template:   meme.Template,
			TextTop:    meme.TextTop,
			TextBottom: meme.TextBottom,
			CreatedAt:  meme.CreatedAt.Format("2006-01-02T15:04:05Z"),
			UpdatedAt:  meme.UpdatedAt.Format("2006-01-02T15:04:05Z"),
		})
	}

	c.JSON(http.StatusOK, response)
}

// DeleteMeme removes a meme by ID
func (h *MemeHandler) DeleteMeme(c *gin.Context) {
	id := c.Param("id")

	err := h.memeUsecase.DeleteMeme(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Meme not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Meme deleted successfully"})
}

// CreateTemplateRequest represents the request body for creating a template
type CreateTemplateRequest struct {
	Name string `json:"name" binding:"required"`
}

// CreateTemplate handles the creation of a new template
func (h *MemeHandler) CreateTemplate(c *gin.Context) {
	var req CreateTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	template, err := h.templateUsecase.CreateTemplate(req.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := TemplateResponse{
		Name:      template.Name,
		CreatedAt: template.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt: template.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}

	c.JSON(http.StatusCreated, response)
}

// ListTemplates retrieves all templates
func (h *MemeHandler) ListTemplates(c *gin.Context) {
	templates, err := h.templateUsecase.ListTemplates()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Ensure we always return an array, even if empty
	if templates == nil {
		templates = []*domain.Template{}
	}

	var response []TemplateResponse
	for _, template := range templates {
		response = append(response, TemplateResponse{
			Name:      template.Name,
			CreatedAt: template.CreatedAt.Format("2006-01-02T15:04:05Z"),
			UpdatedAt: template.UpdatedAt.Format("2006-01-02T15:04:05Z"),
		})
	}

	c.JSON(http.StatusOK, response)
}

// UploadTemplateImage handles uploading an image for a template
func (h *MemeHandler) UploadTemplateImage(c *gin.Context) {
	name := c.Param("name")

	// Get the file from the request
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No image file provided"})
		return
	}

	// Open the uploaded file
	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open uploaded file"})
		return
	}
	defer src.Close()

	// Read the file content
	fileBytes, err := io.ReadAll(src)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read uploaded file"})
		return
	}

	// Determine MIME type
	mimeType := "image/jpeg" // default
	if file.Header.Get("Content-Type") != "" {
		mimeType = file.Header.Get("Content-Type")
	}

	// Save the image using the template usecase
	if err := h.templateUsecase.SaveTemplateImage(name, fileBytes, mimeType); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to save template image: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Template image uploaded successfully"})
}

// ServeTemplateImage serves a template image by name
func (h *MemeHandler) ServeTemplateImage(c *gin.Context) {
	name := c.Param("name")

	// Get the image using the template usecase
	imageData, mimeType, err := h.templateUsecase.GetTemplateImage(name)
	if err != nil {
		// If no image found, serve a placeholder
		h.servePlaceholderImage(c)
		return
	}

	// Set content type header
	c.Header("Content-Type", mimeType)

	// Serve the image data
	c.Data(http.StatusOK, mimeType, imageData)
}

// ServeWebApp serves the React frontend application
func (h *MemeHandler) ServeWebApp(c *gin.Context) {
	// For SPA, serve index.html for all non-API routes
	path := filepath.Join(h.webRoot, "index.html")
	c.File(path)
}

// ServeMemeImage serves a meme image by ID
func (h *MemeHandler) ServeMemeImage(c *gin.Context) {
	id := c.Param("id")

	// Find the first image in the meme's images directory
	imagesDir := filepath.Join("./data/memes", id, "images")

	// Check if directory exists
	if _, err := os.Stat(imagesDir); os.IsNotExist(err) {
		// If images directory doesn't exist, serve a placeholder
		h.servePlaceholderImage(c)
		return
	}

	// Walk the directory to find the first image file
	var imagePath string
	err := filepath.WalkDir(imagesDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if d.IsDir() {
			return nil
		}

		// Check if file is an image
		if isImageFile(path) {
			imagePath = path
			return filepath.SkipDir // Stop walking after finding the first image
		}

		return nil
	})

	if err != nil || imagePath == "" {
		// If no image found, serve a placeholder
		h.servePlaceholderImage(c)
		return
	}

	// Determine content type based on file extension
	contentType := "image/jpeg"
	ext := filepath.Ext(imagePath)
	switch ext {
	case ".png":
		contentType = "image/png"
	case ".gif":
		contentType = "image/gif"
	}

	// Set content type header
	c.Header("Content-Type", contentType)

	// Serve the image file
	c.File(imagePath)
}

// servePlaceholderImage serves a placeholder image when no meme image is available
func (h *MemeHandler) servePlaceholderImage(c *gin.Context) {
	// SVG placeholder image
	placeholder := `<svg xmlns="http://www.w3.org/2000/svg" width="400" height="300" viewBox="0 0 400 300">
		<rect width="100%" height="100%" fill="#ddd"/>
		<text x="50%" y="50%" font-family="Arial" font-size="24" fill="#999" text-anchor="middle" dy=".3em">No Image Available</text>
	</svg>`

	c.Data(http.StatusOK, "image/svg+xml", []byte(placeholder))
}

// isImageFile checks if a file is an image based on its extension
func isImageFile(path string) bool {
	ext := filepath.Ext(path)
	switch ext {
	case ".jpg", ".jpeg", ".png", ".gif":
		return true
	default:
		return false
	}
}
