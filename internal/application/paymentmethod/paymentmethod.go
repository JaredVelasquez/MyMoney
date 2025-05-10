package payment_method

import (
	"context"
	"time"

	"mi-app-backend/internal/domain"
	"mi-app-backend/internal/domain/ports/app"

	"github.com/google/uuid"
)

// Service maneja la lógica de negocio relacionada con métodos de pago
type Service struct {
	repo app.PaymentMethodRepository
}

// NewService crea un nuevo servicio de métodos de pago
func NewService(repo app.PaymentMethodRepository) *Service {
	return &Service{
		repo: repo,
	}
}

// CreatePaymentMethod crea un nuevo método de pago
func (s *Service) CreatePaymentMethod(ctx context.Context, name, description string, userID string) (*domain.PaymentMethod, error) {
	paymentMethod := &domain.PaymentMethod{
		ID:          uuid.New().String(),
		Name:        name,
		Description: description,
		IsActive:    true,
		UserID:      userID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := paymentMethod.Validate(); err != nil {
		return nil, err
	}

	if err := s.repo.Create(ctx, paymentMethod); err != nil {
		return nil, err
	}

	return paymentMethod, nil
}

// GetPaymentMethodByID obtiene un método de pago por su ID
func (s *Service) GetPaymentMethodByID(ctx context.Context, id string) (*domain.PaymentMethod, error) {
	return s.repo.GetByID(ctx, id)
}

// GetPaymentMethodsByUserID obtiene todos los métodos de pago de un usuario
func (s *Service) GetPaymentMethodsByUserID(ctx context.Context, userID string) ([]*domain.PaymentMethod, error) {
	return s.repo.GetByUserID(ctx, userID)
}

// UpdatePaymentMethod actualiza un método de pago existente
func (s *Service) UpdatePaymentMethod(ctx context.Context, id, name, description string, isActive bool) (*domain.PaymentMethod, error) {
	paymentMethod, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if name != "" {
		paymentMethod.Name = name
	}

	paymentMethod.Description = description
	paymentMethod.IsActive = isActive
	paymentMethod.UpdatedAt = time.Now()

	if err := paymentMethod.Validate(); err != nil {
		return nil, err
	}

	if err := s.repo.Update(ctx, paymentMethod); err != nil {
		return nil, err
	}

	return paymentMethod, nil
}

// DeletePaymentMethod elimina un método de pago
func (s *Service) DeletePaymentMethod(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
