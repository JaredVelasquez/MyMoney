package user_subscription

import (
	"context"
	"fmt"
	"time"

	"MyMoneyBackend/internal/domain"
	"MyMoneyBackend/internal/domain/ports/app"
)

// Service implementa la lógica de negocio para las suscripciones de usuarios
type Service struct {
	subscriptionRepo app.UserSubscriptionRepository
	planRepo         app.PlanRepository
	userRepo         app.UserRepository
}

// NewService crea una nueva instancia del servicio de suscripciones
func NewService(
	subscriptionRepo app.UserSubscriptionRepository,
	planRepo app.PlanRepository,
	userRepo app.UserRepository,
) *Service {
	return &Service{
		subscriptionRepo: subscriptionRepo,
		planRepo:         planRepo,
		userRepo:         userRepo,
	}
}

// CreateSubscription crea una nueva suscripción para un usuario
func (s *Service) CreateSubscription(
	ctx context.Context,
	userID string,
	planID string,
	startDate time.Time,
	endDate time.Time,
	paymentMethodID *string,
	metadata map[string]string,
) (*domain.UserSubscription, error) {
	// Verificar que el usuario exista
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, fmt.Errorf("error al verificar usuario: %w", err)
	}
	if user == nil {
		return nil, fmt.Errorf("usuario no encontrado con ID: %s", userID)
	}

	// Verificar que el plan exista
	plan, err := s.planRepo.GetByID(ctx, planID)
	if err != nil {
		return nil, fmt.Errorf("error al verificar plan: %w", err)
	}
	if plan == nil {
		return nil, fmt.Errorf("plan no encontrado con ID: %s", planID)
	}

	// Verificar si ya existe una suscripción activa
	activeSubscription, err := s.subscriptionRepo.GetActiveByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("error al verificar suscripción activa: %w", err)
	}

	// Si ya tiene una suscripción activa, cancelarla antes de crear la nueva
	if activeSubscription != nil {
		err = s.CancelSubscription(ctx, activeSubscription.ID, "Reemplazada por nueva suscripción")
		if err != nil {
			return nil, fmt.Errorf("error al cancelar suscripción existente: %w", err)
		}
	}

	// Calcular fecha de renovación
	var renewalDate *time.Time
	if plan.Interval == domain.PlanIntervalMonthly || plan.Interval == domain.PlanIntervalYearly {
		renewal := endDate.AddDate(0, -1, 0) // Por defecto, renovar un mes antes de finalizar
		if plan.Interval == domain.PlanIntervalYearly {
			renewal = endDate.AddDate(0, -1, 0) // Para planes anuales, renovar un mes antes
		}
		renewalDate = &renewal
	}

	// Para planes gratuitos, no necesitamos método de pago
	if s.isPlanFree(plan) && paymentMethodID != nil {
		paymentMethodID = nil // Ignorar método de pago para planes gratuitos
	}

	// Para planes pagos, requerimos método de pago
	if !s.isPlanFree(plan) && paymentMethodID == nil {
		return nil, fmt.Errorf("se requiere método de pago para planes no gratuitos")
	}

	// Crear la suscripción
	subscription := &domain.UserSubscription{
		UserID:          userID,
		PlanID:          planID,
		Status:          domain.SubscriptionStatusActive,
		StartDate:       startDate,
		EndDate:         endDate,
		RenewalDate:     renewalDate,
		PaymentMethodID: paymentMethodID,
		Metadata:        metadata,
	}

	// Validar la suscripción
	if err := subscription.Validate(); err != nil {
		return nil, fmt.Errorf("error de validación: %w", err)
	}

	// Guardar en la base de datos
	if err := s.subscriptionRepo.Create(ctx, subscription); err != nil {
		return nil, fmt.Errorf("error al crear suscripción: %w", err)
	}

	return subscription, nil
}

// GetSubscriptionByID obtiene una suscripción por su ID
func (s *Service) GetSubscriptionByID(ctx context.Context, id string) (*domain.UserSubscription, error) {
	subscription, err := s.subscriptionRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("error al obtener suscripción: %w", err)
	}
	return subscription, nil
}

// GetActiveSubscription obtiene la suscripción activa de un usuario
func (s *Service) GetActiveSubscription(ctx context.Context, userID string) (*domain.UserSubscription, error) {
	subscription, err := s.subscriptionRepo.GetActiveByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("error al obtener suscripción activa: %w", err)
	}
	return subscription, nil
}

// GetUserSubscriptions obtiene todas las suscripciones de un usuario
func (s *Service) GetUserSubscriptions(ctx context.Context, userID string) ([]*domain.UserSubscription, error) {
	subscriptions, err := s.subscriptionRepo.GetAllByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("error al obtener suscripciones del usuario: %w", err)
	}
	return subscriptions, nil
}

// GetSubscriptionsByStatus obtiene suscripciones por estado
func (s *Service) GetSubscriptionsByStatus(ctx context.Context, status domain.SubscriptionStatus) ([]*domain.UserSubscription, error) {
	subscriptions, err := s.subscriptionRepo.GetByStatus(ctx, status)
	if err != nil {
		return nil, fmt.Errorf("error al obtener suscripciones por estado: %w", err)
	}
	return subscriptions, nil
}

// GetExpiringSubscriptions obtiene suscripciones que expirarán pronto
func (s *Service) GetExpiringSubscriptions(ctx context.Context, daysFromNow int) ([]*domain.UserSubscription, error) {
	beforeDate := time.Now().AddDate(0, 0, daysFromNow)
	subscriptions, err := s.subscriptionRepo.GetExpiringSubscriptions(ctx, beforeDate)
	if err != nil {
		return nil, fmt.Errorf("error al obtener suscripciones por expirar: %w", err)
	}
	return subscriptions, nil
}

// GetPendingRenewals obtiene suscripciones pendientes de renovación
func (s *Service) GetPendingRenewals(ctx context.Context, daysFromNow int) ([]*domain.UserSubscription, error) {
	beforeDate := time.Now().AddDate(0, 0, daysFromNow)
	subscriptions, err := s.subscriptionRepo.GetPendingRenewals(ctx, beforeDate)
	if err != nil {
		return nil, fmt.Errorf("error al obtener renovaciones pendientes: %w", err)
	}
	return subscriptions, nil
}

// RenewSubscription renueva una suscripción existente
func (s *Service) RenewSubscription(ctx context.Context, id string, newEndDate time.Time) (*domain.UserSubscription, error) {
	subscription, err := s.subscriptionRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("error al obtener suscripción: %w", err)
	}

	// Verificar que la suscripción esté activa
	if subscription.Status != domain.SubscriptionStatusActive {
		return nil, fmt.Errorf("solo se pueden renovar suscripciones activas")
	}

	// Obtener el plan para calcular el intervalo de renovación
	plan, err := s.planRepo.GetByID(ctx, subscription.PlanID)
	if err != nil {
		return nil, fmt.Errorf("error al obtener plan: %w", err)
	}

	// Calcular nueva fecha de renovación
	var renewalDate *time.Time
	if plan.Interval == domain.PlanIntervalMonthly || plan.Interval == domain.PlanIntervalYearly {
		renewal := newEndDate.AddDate(0, -1, 0)
		if plan.Interval == domain.PlanIntervalYearly {
			renewal = newEndDate.AddDate(0, -1, 0)
		}
		renewalDate = &renewal
	}

	// Actualizar fechas
	subscription.EndDate = newEndDate
	subscription.RenewalDate = renewalDate
	subscription.LastPaymentDate = timePtr(time.Now())
	subscription.NextPaymentAttempt = nil // Resetear el próximo intento de pago

	// Guardar los cambios
	if err := s.subscriptionRepo.Update(ctx, subscription); err != nil {
		return nil, fmt.Errorf("error al actualizar suscripción: %w", err)
	}

	return subscription, nil
}

// CancelSubscription cancela una suscripción
func (s *Service) CancelSubscription(ctx context.Context, id string, reason string) error {
	subscription, err := s.subscriptionRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("error al obtener suscripción: %w", err)
	}

	if subscription.Status == domain.SubscriptionStatusCancelled {
		return fmt.Errorf("la suscripción ya está cancelada")
	}

	// Guardar el motivo de cancelación en los metadatos
	if subscription.Metadata == nil {
		subscription.Metadata = make(map[string]string)
	}
	subscription.Metadata["cancellation_reason"] = reason

	// Actualizar la suscripción
	cancellationDate := time.Now()
	err = s.subscriptionRepo.CancelSubscription(ctx, id, cancellationDate)
	if err != nil {
		return fmt.Errorf("error al cancelar suscripción: %w", err)
	}

	return nil
}

// ChangeSubscriptionPlan cambia el plan de una suscripción
func (s *Service) ChangeSubscriptionPlan(ctx context.Context, id string, newPlanID string) (*domain.UserSubscription, error) {
	subscription, err := s.subscriptionRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("error al obtener suscripción: %w", err)
	}

	if subscription.Status != domain.SubscriptionStatusActive {
		return nil, fmt.Errorf("solo se puede cambiar el plan de suscripciones activas")
	}

	// Verificar que el nuevo plan exista
	newPlan, err := s.planRepo.GetByID(ctx, newPlanID)
	if err != nil {
		return nil, fmt.Errorf("error al verificar nuevo plan: %w", err)
	}
	if newPlan == nil {
		return nil, fmt.Errorf("plan no encontrado con ID: %s", newPlanID)
	}

	// Obtener el plan actual
	currentPlan, err := s.planRepo.GetByID(ctx, subscription.PlanID)
	if err != nil {
		return nil, fmt.Errorf("error al obtener plan actual: %w", err)
	}

	// Si el nuevo plan es gratuito, no necesitamos método de pago
	if s.isPlanFree(newPlan) {
		subscription.PaymentMethodID = nil
	} else if !s.isPlanFree(newPlan) && subscription.PaymentMethodID == nil {
		// Si pasamos de un plan gratuito a uno de pago, necesitamos un método de pago
		return nil, fmt.Errorf("se requiere método de pago para cambiar a un plan no gratuito")
	}

	// Calcular nuevas fechas según el intervalo del plan
	now := time.Now()
	var newEndDate time.Time

	// Mantener la misma duración relativa
	remaining := subscription.EndDate.Sub(now)

	if currentPlan.Interval != newPlan.Interval {
		// Si cambiamos de tipo de intervalo, ajustar según el nuevo intervalo
		if newPlan.Interval == domain.PlanIntervalMonthly {
			newEndDate = now.AddDate(0, 1, 0)
		} else if newPlan.Interval == domain.PlanIntervalYearly {
			newEndDate = now.AddDate(1, 0, 0)
		} else {
			// Para otros intervalos, usar la fecha actual más un mes como predeterminado
			newEndDate = now.AddDate(0, 1, 0)
		}
	} else {
		// Si mantenemos el mismo tipo de intervalo, conservar tiempo restante
		newEndDate = now.Add(remaining)
	}

	// Calcular nueva fecha de renovación
	var renewalDate *time.Time
	if newPlan.Interval == domain.PlanIntervalMonthly || newPlan.Interval == domain.PlanIntervalYearly {
		renewal := newEndDate.AddDate(0, -1, 0)
		if newPlan.Interval == domain.PlanIntervalYearly {
			renewal = newEndDate.AddDate(0, -1, 0)
		}
		renewalDate = &renewal
	}

	// Actualizar la suscripción
	subscription.PlanID = newPlanID
	subscription.EndDate = newEndDate
	subscription.RenewalDate = renewalDate
	subscription.LastPaymentDate = timePtr(now)

	// Guardar los cambios
	if err := s.subscriptionRepo.Update(ctx, subscription); err != nil {
		return nil, fmt.Errorf("error al actualizar suscripción: %w", err)
	}

	return subscription, nil
}

// UpdatePaymentMethod actualiza el método de pago de una suscripción
func (s *Service) UpdatePaymentMethod(ctx context.Context, id string, paymentMethodID string) (*domain.UserSubscription, error) {
	subscription, err := s.subscriptionRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("error al obtener suscripción: %w", err)
	}

	// Verificar que la suscripción esté activa
	if subscription.Status != domain.SubscriptionStatusActive {
		return nil, fmt.Errorf("solo se puede actualizar el método de pago de suscripciones activas")
	}

	// Obtener el plan para verificar si es gratuito
	plan, err := s.planRepo.GetByID(ctx, subscription.PlanID)
	if err != nil {
		return nil, fmt.Errorf("error al obtener plan: %w", err)
	}

	// Si el plan es gratuito, no permitir establecer un método de pago
	if s.isPlanFree(plan) {
		return nil, fmt.Errorf("no se permite establecer método de pago para planes gratuitos")
	}

	// Actualizar el método de pago
	subscription.PaymentMethodID = &paymentMethodID

	// Guardar los cambios
	if err := s.subscriptionRepo.Update(ctx, subscription); err != nil {
		return nil, fmt.Errorf("error al actualizar suscripción: %w", err)
	}

	return subscription, nil
}

// UpdateSubscriptionStatus actualiza el estado de una suscripción
func (s *Service) UpdateSubscriptionStatus(ctx context.Context, id string, status domain.SubscriptionStatus) (*domain.UserSubscription, error) {
	// Verificar que la suscripción exista
	_, err := s.subscriptionRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("error al obtener suscripción: %w", err)
	}

	// No permitir cambiar a cancelado (hay un método específico para eso)
	if status == domain.SubscriptionStatusCancelled {
		return nil, fmt.Errorf("use CancelSubscription para cancelar suscripciones")
	}

	// Actualizar el estado
	err = s.subscriptionRepo.UpdateStatus(ctx, id, status)
	if err != nil {
		return nil, fmt.Errorf("error al actualizar estado: %w", err)
	}

	// Obtener la suscripción actualizada
	updatedSubscription, err := s.subscriptionRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("error al obtener suscripción actualizada: %w", err)
	}

	return updatedSubscription, nil
}

// isPlanFree determina si un plan es gratuito
func (s *Service) isPlanFree(plan *domain.Plan) bool {
	return plan.Price == 0
}

// timePtr devuelve un puntero a un tiempo
func timePtr(t time.Time) *time.Time {
	return &t
}
