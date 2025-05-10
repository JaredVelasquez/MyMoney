package transaction

import (
	"context"
	"time"

	"mi-app-backend/internal/domain"
	"mi-app-backend/internal/domain/ports/app"

	"github.com/google/uuid"
)

// Service maneja la lógica de negocio relacionada con transacciones
type Service struct {
	repo app.TransactionRepository
}

// NewService crea un nuevo servicio de transacciones
func NewService(repo app.TransactionRepository) *Service {
	return &Service{
		repo: repo,
	}
}

// CreateTransaction crea una nueva transacción
func (s *Service) CreateTransaction(ctx context.Context, amount float64, description string, date time.Time, categoryID, paymentMethodID, userID string) (*domain.Transaction, error) {
	transaction := &domain.Transaction{
		ID:              uuid.New().String(),
		Amount:          amount,
		Description:     description,
		Date:            date,
		CategoryID:      categoryID,
		PaymentMethodID: paymentMethodID,
		UserID:          userID,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	if err := transaction.Validate(); err != nil {
		return nil, err
	}

	if err := s.repo.Create(ctx, transaction); err != nil {
		return nil, err
	}

	return transaction, nil
}

// GetTransactionByID obtiene una transacción por su ID
func (s *Service) GetTransactionByID(ctx context.Context, id string) (*domain.Transaction, error) {
	return s.repo.GetByID(ctx, id)
}

// GetTransactionsByUserID obtiene todas las transacciones de un usuario
func (s *Service) GetTransactionsByUserID(ctx context.Context, userID string) ([]*domain.Transaction, error) {
	return s.repo.GetByUserID(ctx, userID)
}

// GetTransactionsByCategoryID obtiene todas las transacciones de una categoría
func (s *Service) GetTransactionsByCategoryID(ctx context.Context, categoryID string) ([]*domain.Transaction, error) {
	return s.repo.GetByCategoryID(ctx, categoryID)
}

// GetTransactionsByDateRange obtiene todas las transacciones de un usuario en un rango de fechas
func (s *Service) GetTransactionsByDateRange(ctx context.Context, userID string, startDate, endDate time.Time) ([]*domain.Transaction, error) {
	return s.repo.GetByDateRange(ctx, userID, startDate, endDate)
}

// UpdateTransaction actualiza una transacción existente
func (s *Service) UpdateTransaction(ctx context.Context, id string, amount float64, description string, date time.Time, categoryID, paymentMethodID string) (*domain.Transaction, error) {
	transaction, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if amount > 0 {
		transaction.Amount = amount
	}

	transaction.Description = description

	if !date.IsZero() {
		transaction.Date = date
	}

	if categoryID != "" {
		transaction.CategoryID = categoryID
	}

	if paymentMethodID != "" {
		transaction.PaymentMethodID = paymentMethodID
	}

	transaction.UpdatedAt = time.Now()

	if err := transaction.Validate(); err != nil {
		return nil, err
	}

	if err := s.repo.Update(ctx, transaction); err != nil {
		return nil, err
	}

	return transaction, nil
}

// DeleteTransaction elimina una transacción
func (s *Service) DeleteTransaction(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
