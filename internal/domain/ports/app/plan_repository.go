package app

import (
	"context"

	"mi-app-backend/internal/domain"
)

// PlanRepository es la interfaz que define los métodos para el repositorio de planes
type PlanRepository interface {
	// Create crea un nuevo plan en la base de datos
	Create(ctx context.Context, plan *domain.Plan) error

	// GetByID obtiene un plan por su ID
	GetByID(ctx context.Context, id string) (*domain.Plan, error)

	// GetAll obtiene todos los planes
	GetAll(ctx context.Context) ([]*domain.Plan, error)

	// GetAllPublic obtiene todos los planes públicos
	GetAllPublic(ctx context.Context) ([]*domain.Plan, error)

	// GetAllActive obtiene todos los planes activos
	GetAllActive(ctx context.Context) ([]*domain.Plan, error)

	// Update actualiza un plan existente
	Update(ctx context.Context, plan *domain.Plan) error

	// Delete elimina un plan por su ID
	Delete(ctx context.Context, id string) error
}
