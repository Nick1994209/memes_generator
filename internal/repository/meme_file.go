package repository

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"memes-generator/internal/config"
	"memes-generator/internal/domain"
)

const (
	dataDir = ""
)

// MemeFileRepository implements domain.MemeRepository using file system
type MemeFileRepository struct {
	dataPath string
}

// NewMemeFileRepository creates a new file-based meme repository
func NewMemeFileRepository() *MemeFileRepository {
	return &MemeFileRepository{
		dataPath: config.GetMemesDir(),
	}
}

// Create saves a new meme to the file system
func (r *MemeFileRepository) Create(meme *domain.Meme) error {
	// Generate unique ID if not set
	if meme.ID == "" {
		meme.ID = r.GenerateID()
	}

	// Create meme directory
	memeDir := filepath.Join(r.dataPath, meme.ID)
	if err := os.MkdirAll(memeDir, 0755); err != nil {
		return fmt.Errorf("failed to create meme directory: %w", err)
	}

	// Save metadata
	metadataPath := filepath.Join(memeDir, "metadata.json")
	file, err := os.Create(metadataPath)
	if err != nil {
		return fmt.Errorf("failed to create metadata file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(meme); err != nil {
		return fmt.Errorf("failed to encode meme metadata: %w", err)
	}

	return nil
}

// GetByID retrieves a meme by its ID
func (r *MemeFileRepository) GetByID(id string) (*domain.Meme, error) {
	memeDir := filepath.Join(r.dataPath, id)
	metadataPath := filepath.Join(memeDir, "metadata.json")

	// Check if meme exists
	if _, err := os.Stat(memeDir); os.IsNotExist(err) {
		return nil, fmt.Errorf("meme with ID %s not found", id)
	}

	// Read metadata
	file, err := os.Open(metadataPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open metadata file: %w", err)
	}
	defer file.Close()

	var meme domain.Meme
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&meme); err != nil {
		return nil, fmt.Errorf("failed to decode meme metadata: %w", err)
	}

	return &meme, nil
}

// List returns all memes
func (r *MemeFileRepository) List() ([]*domain.Meme, error) {
	// Check if data directory exists
	if _, err := os.Stat(r.dataPath); os.IsNotExist(err) {
		return []*domain.Meme{}, nil
	}

	entries, err := os.ReadDir(r.dataPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read data directory: %w", err)
	}

	var memes []*domain.Meme
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		// Try to load meme metadata
		meme, err := r.GetByID(entry.Name())
		if err != nil {
			// Skip invalid meme directories
			continue
		}

		memes = append(memes, meme)
	}

	// Ensure we always return an array, even if empty
	if memes == nil {
		memes = []*domain.Meme{}
	}

	return memes, nil
}

// Delete removes a meme by its ID
func (r *MemeFileRepository) Delete(id string) error {
	memeDir := filepath.Join(r.dataPath, id)

	// Check if meme exists
	if _, err := os.Stat(memeDir); os.IsNotExist(err) {
		return fmt.Errorf("meme with ID %s not found", id)
	}

	// Remove meme directory
	if err := os.RemoveAll(memeDir); err != nil {
		return fmt.Errorf("failed to delete meme directory: %w", err)
	}

	return nil
}

// GenerateID creates a unique ID for a meme
func (r *MemeFileRepository) GenerateID() string {
	// In a real application, you would use a proper UUID generator
	// For simplicity, we'll use timestamp-based ID
	// Ensure uniqueness by checking if ID already exists
	for {
		id := fmt.Sprintf("meme_%d", time.Now().UnixNano())
		memeDir := filepath.Join(r.dataPath, id)
		if _, err := os.Stat(memeDir); os.IsNotExist(err) {
			return id
		}
		// If ID exists, generate a new one
		time.Sleep(time.Nanosecond)
	}
}
