package plan

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	"MyMoneyBackend/internal/domain"
	"MyMoneyBackend/internal/domain/ports/app"
)

// Service encapsula la lógica de negocio relacionada con los planes
type Service struct {
	planRepo     app.PlanRepository
	currencyRepo app.CurrencyRepository
}

// NewService crea una nueva instancia del servicio de planes
func NewService(planRepo app.PlanRepository, currencyRepo app.CurrencyRepository) *Service {
	return &Service{
		planRepo:     planRepo,
		currencyRepo: currencyRepo,
	}
}

// CreatePlan crea un nuevo plan
func (s *Service) CreatePlan(
	ctx context.Context,
	name, description string,
	price float64,
	currencyID string,
	interval domain.PlanInterval,
	features []domain.PlanFeature,
	isActive, isPublic bool,
	sortOrder int,
) (*domain.Plan, error) {
	// Validar datos básicos
	if name == "" {
		return nil, errors.New("el nombre del plan es obligatorio")
	}
	if description == "" {
		return nil, errors.New("la descripción del plan es obligatoria")
	}
	if price < 0 {
		return nil, errors.New("el precio no puede ser negativo")
	}
	if currencyID == "" {
		return nil, errors.New("la moneda es obligatoria")
	}

	// Validar que la moneda exista
	currency, err := s.currencyRepo.GetByID(ctx, currencyID)
	if err != nil {
		return nil, errors.New("moneda no válida: " + err.Error())
	}
	if !currency.IsActive {
		return nil, errors.New("la moneda seleccionada no está activa")
	}

	// Validar intervalo
	if interval != domain.PlanIntervalMonthly && interval != domain.PlanIntervalYearly {
		return nil, errors.New("el intervalo de facturación debe ser mensual o anual")
	}

	// Crear el nuevo plan
	plan := &domain.Plan{
		ID:          uuid.New().String(),
		Name:        name,
		Description: description,
		Price:       price,
		CurrencyID:  currencyID,
		Interval:    interval,
		Features:    features,
		IsActive:    isActive,
		IsPublic:    isPublic,
		SortOrder:   sortOrder,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Guardar el plan en la base de datos
	if err := s.planRepo.Create(ctx, plan); err != nil {
		return nil, err
	}

	return plan, nil
}

// GetPlanByID obtiene un plan por su ID
func (s *Service) GetPlanByID(ctx context.Context, id string) (*domain.Plan, error) {
	return s.planRepo.GetByID(ctx, id)
}

// GetAllPlans obtiene todos los planes
func (s *Service) GetAllPlans(ctx context.Context) ([]*domain.Plan, error) {
	return s.planRepo.GetAll(ctx)
}

// GetPublicPlans obtiene todos los planes públicos
func (s *Service) GetPublicPlans(ctx context.Context) ([]*domain.Plan, error) {
	return s.planRepo.GetAllPublic(ctx)
}

// GetActivePlans obtiene todos los planes activos
func (s *Service) GetActivePlans(ctx context.Context) ([]*domain.Plan, error) {
	return s.planRepo.GetAllActive(ctx)
}

// UpdatePlan actualiza un plan existente
func (s *Service) UpdatePlan(
	ctx context.Context,
	id, name, description string,
	price float64,
	currencyID string,
	interval domain.PlanInterval,
	features []domain.PlanFeature,
	isActive, isPublic bool,
	sortOrder int,
) (*domain.Plan, error) {
	// Validar datos básicos
	if id == "" {
		return nil, errors.New("el ID del plan es obligatorio")
	}
	if name == "" {
		return nil, errors.New("el nombre del plan es obligatorio")
	}
	if description == "" {
		return nil, errors.New("la descripción del plan es obligatoria")
	}
	if price < 0 {
		return nil, errors.New("el precio no puede ser negativo")
	}
	if currencyID == "" {
		return nil, errors.New("la moneda es obligatoria")
	}

	// Validar que la moneda exista
	currency, err := s.currencyRepo.GetByID(ctx, currencyID)
	if err != nil {
		return nil, errors.New("moneda no válida: " + err.Error())
	}
	if !currency.IsActive {
		return nil, errors.New("la moneda seleccionada no está activa")
	}

	// Validar intervalo
	if interval != domain.PlanIntervalMonthly && interval != domain.PlanIntervalYearly {
		return nil, errors.New("el intervalo de facturación debe ser mensual o anual")
	}

	// Obtener el plan existente
	plan, err := s.planRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Actualizar los datos del plan
	plan.Name = name
	plan.Description = description
	plan.Price = price
	plan.CurrencyID = currencyID
	plan.Interval = interval
	plan.Features = features
	plan.IsActive = isActive
	plan.IsPublic = isPublic
	plan.SortOrder = sortOrder
	plan.UpdatedAt = time.Now()

	// Guardar los cambios en la base de datos
	if err := s.planRepo.Update(ctx, plan); err != nil {
		return nil, err
	}

	return plan, nil
}

// DeletePlan elimina un plan por su ID
func (s *Service) DeletePlan(ctx context.Context, id string) error {
	// Verificar si hay suscripciones activas con este plan
	// Esta verificación debería implementarse cuando se cree el repositorio de suscripciones

	return s.planRepo.Delete(ctx, id)
}

// IsFreePlan verifica si un plan es gratuito
func (s *Service) IsFreePlan(plan *domain.Plan) bool {
	return plan.Price == 0
}
