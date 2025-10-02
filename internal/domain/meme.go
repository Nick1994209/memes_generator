package domain

import (
	"time"

	"memes-generator/internal/meme"
)

// Meme represents a meme entity
type Meme struct {
	ID         string    `json:"id"`
	Template   string    `json:"template"`
	TextTop    string    `json:"text_top"`
	TextBottom string    `json:"text_bottom"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// CreateMemeImage generates a meme image for this meme entity
func (m *Meme) CreateMemeImage(textTop, textBottom, outputPath string) error {
	// Try to create meme from template
	if err := meme.CreateMemeFromTemplate(m.Template, textTop, textBottom, outputPath); err != nil {
		// If template not found, create a simple default template
		return meme.CreateMemeImage(textTop, textBottom, outputPath)
	}
	return nil
}

// MemeRepository defines the interface for meme data operations
type MemeRepository interface {
	Create(meme *Meme) error
	GetByID(id string) (*Meme, error)
	List() ([]*Meme, error)
	Delete(id string) error
}

// MemeUsecase defines the interface for meme business logic
type MemeUsecase interface {
	CreateMeme(template, textTop, textBottom string) (*Meme, error)
	GetMemeByID(id string) (*Meme, error)
	ListMemes() ([]*Meme, error)
	DeleteMeme(id string) error
}
