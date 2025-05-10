package category

import (
	"context"
	"time"

	"mi-app-backend/internal/domain"
	"mi-app-backend/internal/domain/ports/app"

	"github.com/google/uuid"
)

// Service maneja la lógica de negocio relacionada con categorías
type Service struct {
	repo app.CategoryRepository
}

// NewService crea un nuevo servicio de categorías
func NewService(repo app.CategoryRepository) *Service {
	return &Service{
		repo: repo,
	}
}

// CreateCategory crea una nueva categoría
func (s *Service) CreateCategory(ctx context.Context, name, description string, categoryType domain.CategoryType, icon, color, userID string) (*domain.Category, error) {
	category := &domain.Category{
		ID:          uuid.New().String(),
		Name:        name,
		Description: description,
		Type:        categoryType,
		Icon:        icon,
		Color:       color,
		UserID:      userID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := category.Validate(); err != nil {
		return nil, err
	}

	if err := s.repo.Create(ctx, category); err != nil {
		return nil, err
	}

	return category, nil
}

// GetCategoryByID obtiene una categoría por su ID
func (s *Service) GetCategoryByID(ctx context.Context, id string) (*domain.Category, error) {
	return s.repo.GetByID(ctx, id)
}

// GetCategoriesByUserID obtiene todas las categorías de un usuario
func (s *Service) GetCategoriesByUserID(ctx context.Context, userID string) ([]*domain.Category, error) {
	return s.repo.GetByUserID(ctx, userID)
}

// GetCategoriesByType obtiene todas las categorías de un usuario por tipo
func (s *Service) GetCategoriesByType(ctx context.Context, userID string, categoryType domain.CategoryType) ([]*domain.Category, error) {
	return s.repo.GetByType(ctx, userID, categoryType)
}

// UpdateCategory actualiza una categoría existente
func (s *Service) UpdateCategory(ctx context.Context, id, name, description string, categoryType domain.CategoryType, icon, color string) (*domain.Category, error) {
	category, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if name != "" {
		category.Name = name
	}

	category.Description = description

	if categoryType != "" {
		category.Type = categoryType
	}

	if icon != "" {
		category.Icon = icon
	}

	if color != "" {
		category.Color = color
	}

	category.UpdatedAt = time.Now()

	if err := category.Validate(); err != nil {
		return nil, err
	}

	if err := s.repo.Update(ctx, category); err != nil {
		return nil, err
	}

	return category, nil
}

// DeleteCategory elimina una categoría
func (s *Service) DeleteCategory(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
