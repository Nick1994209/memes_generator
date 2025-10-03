package usecase

import (
	"log"
	"path/filepath"
	"time"

	"memes-generator/internal/config"
	"memes-generator/internal/domain"
	"memes-generator/internal/service"
)

// MemeUsecase implements domain.MemeUsecase
type MemeUsecase struct {
	memeRepo domain.MemeRepository
}

// NewMemeUsecase creates a new meme usecase
func NewMemeUsecase(memeRepo domain.MemeRepository) *MemeUsecase {
	return &MemeUsecase{
		memeRepo: memeRepo,
	}
}

// CreateMeme creates a new meme
func (uc *MemeUsecase) CreateMeme(template, textTop, textBottom string) (*domain.Meme, error) {
	meme := &domain.Meme{
		Template:   template,
		TextTop:    textTop,
		TextBottom: textBottom,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	if err := uc.memeRepo.Create(meme); err != nil {
		return nil, err
	}

	// Handle different generation modes based on environment variable
	if config.IsBackgroundMode() {
		// Generate meme in background using goroutine
		go func() {
			imagesDir := filepath.Join(config.GetMemesDir(), meme.ID, "images")
			if err := meme.CreateMemeImage(textTop, textBottom, filepath.Join(imagesDir, "generated_meme.png")); err != nil {
				log.Printf("Failed to generate meme in background: %v", err)
			}
		}()
	} else if config.IsContainerAppJobMode() {
		// Generate meme using Cloud.ru Container App Job
		go func() {
			jobService := service.NewCloudRuJobService()
			memePath := filepath.Join(config.GetMemesDir(), meme.ID)
			if err := jobService.ExecuteJob("generate-meme-job", memePath); err != nil {
				log.Printf("Failed to execute Container App Job: %v", err)
			}
		}()
	} else {
		// Default behavior - generate meme synchronously
		imagesDir := filepath.Join(config.GetMemesDir(), meme.ID, "images")
		if err := meme.CreateMemeImage(textTop, textBottom, filepath.Join(imagesDir, "generated_meme.png")); err != nil {
			// If image generation fails, we still return the meme but log the error
			// In a production environment, you might want to handle this differently
			return meme, nil
		}
	}

	return meme, nil
}

// GetMemeByID retrieves a meme by its ID
func (uc *MemeUsecase) GetMemeByID(id string) (*domain.Meme, error) {
	return uc.memeRepo.GetByID(id)
}

// ListMemes returns all memes
func (uc *MemeUsecase) ListMemes() ([]*domain.Meme, error) {
	return uc.memeRepo.List()
}

// DeleteMeme removes a meme by its ID
func (uc *MemeUsecase) DeleteMeme(id string) error {
	return uc.memeRepo.Delete(id)
}
