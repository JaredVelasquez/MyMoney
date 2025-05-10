package currency

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	"mi-app-backend/internal/domain"
	"mi-app-backend/internal/domain/ports/app"
)

// Service encapsula la lógica de negocio relacionada con las monedas
type Service struct {
	repo app.CurrencyRepository
}

// NewService crea una nueva instancia del servicio de monedas
func NewService(repo app.CurrencyRepository) *Service {
	return &Service{
		repo: repo,
	}
}

// CreateCurrency crea una nueva moneda
func (s *Service) CreateCurrency(ctx context.Context, code, name, symbol string, isActive bool) (*domain.Currency, error) {
	// Validar que los datos requeridos estén presentes
	if code == "" {
		return nil, errors.New("el código de la moneda es obligatorio")
	}
	if name == "" {
		return nil, errors.New("el nombre de la moneda es obligatorio")
	}
	if symbol == "" {
		return nil, errors.New("el símbolo de la moneda es obligatorio")
	}

	// Verificar si ya existe una moneda con el mismo código
	existingCurrency, err := s.repo.GetByCode(ctx, code)
	if err == nil && existingCurrency != nil {
		return nil, errors.New("ya existe una moneda con este código")
	}

	// Crear la nueva moneda
	currency := &domain.Currency{
		ID:        uuid.New().String(),
		Code:      code,
		Name:      name,
		Symbol:    symbol,
		IsActive:  isActive,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Guardar la moneda en la base de datos
	if err := s.repo.Create(ctx, currency); err != nil {
		return nil, err
	}

	return currency, nil
}

// GetCurrencyByID obtiene una moneda por su ID
func (s *Service) GetCurrencyByID(ctx context.Context, id string) (*domain.Currency, error) {
	return s.repo.GetByID(ctx, id)
}

// GetCurrencyByCode obtiene una moneda por su código
func (s *Service) GetCurrencyByCode(ctx context.Context, code string) (*domain.Currency, error) {
	return s.repo.GetByCode(ctx, code)
}

// GetAllCurrencies obtiene todas las monedas
func (s *Service) GetAllCurrencies(ctx context.Context) ([]*domain.Currency, error) {
	return s.repo.GetAll(ctx)
}

// GetActiveCurrencies obtiene todas las monedas activas
func (s *Service) GetActiveCurrencies(ctx context.Context) ([]*domain.Currency, error) {
	return s.repo.GetAllActive(ctx)
}

// UpdateCurrency actualiza una moneda existente
func (s *Service) UpdateCurrency(ctx context.Context, id, code, name, symbol string, isActive bool) (*domain.Currency, error) {
	// Validar que los datos requeridos estén presentes
	if id == "" {
		return nil, errors.New("el ID de la moneda es obligatorio")
	}
	if code == "" {
		return nil, errors.New("el código de la moneda es obligatorio")
	}
	if name == "" {
		return nil, errors.New("el nombre de la moneda es obligatorio")
	}
	if symbol == "" {
		return nil, errors.New("el símbolo de la moneda es obligatorio")
	}

	// Obtener la moneda existente
	currency, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Verificar si el nuevo código ya existe en otra moneda
	if code != currency.Code {
		existingCurrency, err := s.repo.GetByCode(ctx, code)
		if err == nil && existingCurrency != nil && existingCurrency.ID != id {
			return nil, errors.New("ya existe otra moneda con este código")
		}
	}

	// Actualizar los datos de la moneda
	currency.Code = code
	currency.Name = name
	currency.Symbol = symbol
	currency.IsActive = isActive
	currency.UpdatedAt = time.Now()

	// Guardar los cambios en la base de datos
	if err := s.repo.Update(ctx, currency); err != nil {
		return nil, err
	}

	return currency, nil
}

// DeleteCurrency elimina una moneda por su ID
func (s *Service) DeleteCurrency(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
