package app

import (
	"context"

	"MyMoneyBackend/internal/domain"
)

// CategoryRepository define las operaciones para el repositorio de categorías
type CategoryRepository interface {
	// Create crea una nueva categoría
	Create(ctx context.Context, category *domain.Category) error

	// GetByID obtiene una categoría por su ID
	GetByID(ctx context.Context, id string) (*domain.Category, error)

	// GetByUserID obtiene todas las categorías de un usuario
	GetByUserID(ctx context.Context, userID string) ([]*domain.Category, error)

	// Update actualiza una categoría existente
	Update(ctx context.Context, category *domain.Category) error

	// Delete elimina una categoría
	Delete(ctx context.Context, id string) error
}
