package app

import (
	"context"
	"time"

	"mi-app-backend/internal/domain"
)

// UserSubscriptionRepository es la interfaz que define los métodos para el repositorio de suscripciones
type UserSubscriptionRepository interface {
	// Create crea una nueva suscripción en la base de datos
	Create(ctx context.Context, subscription *domain.UserSubscription) error

	// GetByID obtiene una suscripción por su ID
	GetByID(ctx context.Context, id string) (*domain.UserSubscription, error)

	// GetActiveByUserID obtiene la suscripción activa de un usuario
	GetActiveByUserID(ctx context.Context, userID string) (*domain.UserSubscription, error)

	// GetAllByUserID obtiene todas las suscripciones de un usuario
	GetAllByUserID(ctx context.Context, userID string) ([]*domain.UserSubscription, error)

	// GetByStatus obtiene suscripciones por estado
	GetByStatus(ctx context.Context, status domain.SubscriptionStatus) ([]*domain.UserSubscription, error)

	// GetExpiringSubscriptions obtiene suscripciones que expirarán pronto
	GetExpiringSubscriptions(ctx context.Context, beforeDate time.Time) ([]*domain.UserSubscription, error)

	// GetPendingRenewals obtiene suscripciones pendientes de renovación
	GetPendingRenewals(ctx context.Context, beforeDate time.Time) ([]*domain.UserSubscription, error)

	// Update actualiza una suscripción existente
	Update(ctx context.Context, subscription *domain.UserSubscription) error

	// UpdateStatus actualiza el estado de una suscripción
	UpdateStatus(ctx context.Context, id string, status domain.SubscriptionStatus) error

	// CancelSubscription cancela una suscripción
	CancelSubscription(ctx context.Context, id string, cancellationDate time.Time) error

	// Delete elimina una suscripción por su ID
	Delete(ctx context.Context, id string) error
}
