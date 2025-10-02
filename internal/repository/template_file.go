package repository

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"memes-generator/internal/domain"
)

const (
	templatesDir = "./data/templates"
)

// TemplateFileRepository implements domain.TemplateRepository using file system
type TemplateFileRepository struct {
	dataPath string
}

// NewTemplateFileRepository creates a new file-based template repository
func NewTemplateFileRepository() *TemplateFileRepository {
	return &TemplateFileRepository{
		dataPath: templatesDir,
	}
}

// Create saves a new template to the file system
func (r *TemplateFileRepository) Create(template *domain.Template) error {
	// Create template directory
	templateDir := filepath.Join(r.dataPath, template.Name)
	if err := os.MkdirAll(templateDir, 0755); err != nil {
		return fmt.Errorf("failed to create template directory: %w", err)
	}

	// Save metadata
	metadataPath := filepath.Join(templateDir, "metadata.json")
	file, err := os.Create(metadataPath)
	if err != nil {
		return fmt.Errorf("failed to create metadata file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(template); err != nil {
		return fmt.Errorf("failed to encode template metadata: %w", err)
	}

	return nil
}

// GetByName retrieves a template by its name
func (r *TemplateFileRepository) GetByName(name string) (*domain.Template, error) {
	templateDir := filepath.Join(r.dataPath, name)
	metadataPath := filepath.Join(templateDir, "metadata.json")

	// Check if template exists
	if _, err := os.Stat(templateDir); os.IsNotExist(err) {
		return nil, fmt.Errorf("template with name %s not found", name)
	}

	// Read metadata
	file, err := os.Open(metadataPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open metadata file: %w", err)
	}
	defer file.Close()

	var template domain.Template
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&template); err != nil {
		return nil, fmt.Errorf("failed to decode template metadata: %w", err)
	}

	return &template, nil
}

// List returns all templates
func (r *TemplateFileRepository) List() ([]*domain.Template, error) {
	// Check if data directory exists
	if _, err := os.Stat(r.dataPath); os.IsNotExist(err) {
		return []*domain.Template{}, nil
	}

	entries, err := os.ReadDir(r.dataPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read data directory: %w", err)
	}

	var templates []*domain.Template
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		// Try to load template metadata
		template, err := r.GetByName(entry.Name())
		if err != nil {
			// Skip invalid template directories
			continue
		}

		templates = append(templates, template)
	}

	return templates, nil
}

// Delete removes a template by its name
func (r *TemplateFileRepository) Delete(name string) error {
	templateDir := filepath.Join(r.dataPath, name)

	// Check if template exists
	if _, err := os.Stat(templateDir); os.IsNotExist(err) {
		return fmt.Errorf("template with name %s not found", name)
	}

	// Remove template directory
	if err := os.RemoveAll(templateDir); err != nil {
		return fmt.Errorf("failed to delete template directory: %w", err)
	}

	return nil
}

// SaveImage saves an image file for a template
func (r *TemplateFileRepository) SaveImage(name string, imageData []byte, mimeType string) error {
	templateDir := filepath.Join(r.dataPath, name)

	// Create images directory
	imagesDir := filepath.Join(templateDir, "images")
	if err := os.MkdirAll(imagesDir, 0755); err != nil {
		return fmt.Errorf("failed to create images directory: %w", err)
	}

	// Determine file extension based on MIME type
	var ext string
	switch mimeType {
	case "image/jpeg":
		ext = ".jpg"
	case "image/png":
		ext = ".png"
	default:
		ext = ".jpg" // default to jpg
	}

	// Save image file
	imagePath := filepath.Join(imagesDir, "template"+ext)
	if err := os.WriteFile(imagePath, imageData, 0644); err != nil {
		return fmt.Errorf("failed to save template image: %w", err)
	}

	return nil
}

// GetImage retrieves the image file for a template
func (r *TemplateFileRepository) GetImage(name string) ([]byte, string, error) {
	templateDir := filepath.Join(r.dataPath, name)
	imagesDir := filepath.Join(templateDir, "images")

	// Check if images directory exists
	if _, err := os.Stat(imagesDir); os.IsNotExist(err) {
		return nil, "", fmt.Errorf("images directory not found for template %s", name)
	}

	// Walk the directory to find the first image file
	var imagePath string
	var mimeType string
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
			ext := strings.ToLower(filepath.Ext(path))
			switch ext {
			case ".jpg", ".jpeg":
				mimeType = "image/jpeg"
			case ".png":
				mimeType = "image/png"
			}
			return filepath.SkipDir // Stop walking after finding the first image
		}

		return nil
	})

	if err != nil || imagePath == "" {
		return nil, "", fmt.Errorf("no image found for template %s", name)
	}

	// Read image file
	imageData, err := os.ReadFile(imagePath)
	if err != nil {
		return nil, "", fmt.Errorf("failed to read template image: %w", err)
	}

	return imageData, mimeType, nil
}

// isImageFile checks if a file is an image based on its extension
func isImageFile(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	return ext == ".jpg" || ext == ".jpeg" || ext == ".png" || ext == ".gif"
}
