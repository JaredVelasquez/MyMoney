package app

import (
	"context"

	"MyMoneyBackend/internal/domain"
)

// CurrencyRepository es la interfaz que define los métodos para el repositorio de monedas
type CurrencyRepository interface {
	// Create crea una nueva moneda en la base de datos
	Create(ctx context.Context, currency *domain.Currency) error

	// GetByID obtiene una moneda por su ID
	GetByID(ctx context.Context, id string) (*domain.Currency, error)

	// GetByCode obtiene una moneda por su código
	GetByCode(ctx context.Context, code string) (*domain.Currency, error)

	// GetAll obtiene todas las monedas
	GetAll(ctx context.Context) ([]*domain.Currency, error)

	// GetAllActive obtiene todas las monedas activas
	GetAllActive(ctx context.Context) ([]*domain.Currency, error)

	// Update actualiza una moneda existente
	Update(ctx context.Context, currency *domain.Currency) error

	// Delete elimina una moneda por su ID
	Delete(ctx context.Context, id string) error
}
