package domain

import (
	"time"
)

// Template represents a meme template entity
type Template struct {
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TemplateRepository defines the interface for template data operations
type TemplateRepository interface {
	Create(template *Template) error
	GetByName(name string) (*Template, error)
	List() ([]*Template, error)
	Delete(name string) error
}
