package usecase

import (
	"time"

	"memes-generator/internal/domain"
	"memes-generator/internal/repository"
)

// TemplateUsecase implements template business logic
type TemplateUsecase struct {
	templateRepo domain.TemplateRepository
}

// NewTemplateUsecase creates a new template usecase
func NewTemplateUsecase(templateRepo domain.TemplateRepository) *TemplateUsecase {
	return &TemplateUsecase{
		templateRepo: templateRepo,
	}
}

// CreateTemplate creates a new template
func (uc *TemplateUsecase) CreateTemplate(name string) (*domain.Template, error) {
	template := &domain.Template{
		Name:      name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := uc.templateRepo.Create(template); err != nil {
		return nil, err
	}

	return template, nil
}

// GetTemplateByName retrieves a template by its name
func (uc *TemplateUsecase) GetTemplateByName(name string) (*domain.Template, error) {
	return uc.templateRepo.GetByName(name)
}

// ListTemplates returns all templates
func (uc *TemplateUsecase) ListTemplates() ([]*domain.Template, error) {
	return uc.templateRepo.List()
}

// DeleteTemplate removes a template by its name
func (uc *TemplateUsecase) DeleteTemplate(name string) error {
	return uc.templateRepo.Delete(name)
}

// SaveTemplateImage saves an image for a template
func (uc *TemplateUsecase) SaveTemplateImage(name string, imageData []byte, mimeType string) error {
	// First verify that the template exists
	if _, err := uc.templateRepo.GetByName(name); err != nil {
		return err
	}

	// Cast to the file repository to access SaveImage method
	if fileRepo, ok := uc.templateRepo.(*repository.TemplateFileRepository); ok {
		return fileRepo.SaveImage(name, imageData, mimeType)
	}

	return nil
}

// GetTemplateImage retrieves the image for a template
func (uc *TemplateUsecase) GetTemplateImage(name string) ([]byte, string, error) {
	// Cast to the file repository to access GetImage method
	if fileRepo, ok := uc.templateRepo.(*repository.TemplateFileRepository); ok {
		return fileRepo.GetImage(name)
	}

	return nil, "", nil
}
