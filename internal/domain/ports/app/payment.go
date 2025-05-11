package app

import (
	"context"

	"MyMoneyBackend/internal/domain"
)

// PaymentMethodRepository define las operaciones para el repositorio de métodos de pago
type PaymentMethodRepository interface {
	// Create crea un nuevo método de pago
	Create(ctx context.Context, paymentMethod *domain.PaymentMethod) error

	// GetByID obtiene un método de pago por su ID
	GetByID(ctx context.Context, id string) (*domain.PaymentMethod, error)

	// GetByUserID obtiene todos los métodos de pago de un usuario
	GetByUserID(ctx context.Context, userID string) ([]*domain.PaymentMethod, error)

	// Update actualiza un método de pago existente
	Update(ctx context.Context, paymentMethod *domain.PaymentMethod) error

	// Delete elimina un método de pago
	Delete(ctx context.Context, id string) error
}
