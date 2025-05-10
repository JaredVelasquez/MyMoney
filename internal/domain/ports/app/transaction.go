package app

import (
	"context"
	"time"

	"mi-app-backend/internal/domain"
)

// TransactionRepository define las operaciones para el repositorio de transacciones
type TransactionRepository interface {
	// Create crea una nueva transacción
	Create(ctx context.Context, transaction *domain.Transaction) error

	// GetByID obtiene una transacción por su ID
	GetByID(ctx context.Context, id string) (*domain.Transaction, error)

	// GetByUserID obtiene todas las transacciones de un usuario
	GetByUserID(ctx context.Context, userID string) ([]*domain.Transaction, error)

	// GetByCategoryID obtiene todas las transacciones de una categoría
	GetByCategoryID(ctx context.Context, categoryID string) ([]*domain.Transaction, error)

	// GetByDateRange obtiene todas las transacciones de un usuario en un rango de fechas
	GetByDateRange(ctx context.Context, userID string, startDate, endDate time.Time) ([]*domain.Transaction, error)

	// Update actualiza una transacción existente
	Update(ctx context.Context, transaction *domain.Transaction) error

	// Delete elimina una transacción
	Delete(ctx context.Context, id string) error
}
