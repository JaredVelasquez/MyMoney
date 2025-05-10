package domain

import (
	"errors"
	"time"
)

// SubscriptionStatus define el estado de una suscripción
type SubscriptionStatus string

const (
	// SubscriptionStatusActive representa una suscripción activa
	SubscriptionStatusActive SubscriptionStatus = "active"
	// SubscriptionStatusCancelled representa una suscripción cancelada
	SubscriptionStatusCancelled SubscriptionStatus = "cancelled"
	// SubscriptionStatusExpired representa una suscripción expirada
	SubscriptionStatusExpired SubscriptionStatus = "expired"
	// SubscriptionStatusPending representa una suscripción pendiente de pago
	SubscriptionStatusPending SubscriptionStatus = "pending"
	// SubscriptionStatusFailed representa una suscripción con pago fallido
	SubscriptionStatusFailed SubscriptionStatus = "failed"
)

// UserSubscription representa una suscripción de un usuario a un plan
type UserSubscription struct {
	ID                 string             `json:"id"`
	UserID             string             `json:"user_id"`              // ID del usuario
	PlanID             string             `json:"plan_id"`              // ID del plan
	Status             SubscriptionStatus `json:"status"`               // Estado de la suscripción
	StartDate          time.Time          `json:"start_date"`           // Fecha de inicio
	EndDate            time.Time          `json:"end_date"`             // Fecha de finalización
	RenewalDate        *time.Time         `json:"renewal_date"`         // Fecha de renovación
	CancellationDate   *time.Time         `json:"cancellation_date"`    // Fecha de cancelación
	LastPaymentDate    *time.Time         `json:"last_payment_date"`    // Fecha del último pago
	NextPaymentAttempt *time.Time         `json:"next_payment_attempt"` // Fecha del próximo intento de pago
	PaymentMethodID    *string            `json:"payment_method_id"`    // ID del método de pago (null para planes gratuitos)
	Metadata           map[string]string  `json:"metadata"`             // Metadatos adicionales
	CreatedAt          time.Time          `json:"created_at"`           // Fecha de creación
	UpdatedAt          time.Time          `json:"updated_at"`           // Fecha de actualización
}

// Validate valida que la entidad UserSubscription tenga todos los campos requeridos
func (s *UserSubscription) Validate() error {
	if s.UserID == "" {
		return errors.New("el ID del usuario es obligatorio")
	}
	if s.PlanID == "" {
		return errors.New("el ID del plan es obligatorio")
	}
	if s.Status == "" {
		return errors.New("el estado de la suscripción es obligatorio")
	}
	if s.StartDate.IsZero() {
		return errors.New("la fecha de inicio es obligatoria")
	}
	if s.EndDate.IsZero() {
		return errors.New("la fecha de finalización es obligatoria")
	}

	// Validar estado de suscripción
	validStatus := map[SubscriptionStatus]bool{
		SubscriptionStatusActive:    true,
		SubscriptionStatusCancelled: true,
		SubscriptionStatusExpired:   true,
		SubscriptionStatusPending:   true,
		SubscriptionStatusFailed:    true,
	}

	if !validStatus[s.Status] {
		return errors.New("el estado de la suscripción no es válido")
	}

	// Validar fechas
	if s.EndDate.Before(s.StartDate) {
		return errors.New("la fecha de finalización debe ser posterior a la fecha de inicio")
	}

	if s.RenewalDate != nil && s.RenewalDate.Before(s.StartDate) {
		return errors.New("la fecha de renovación debe ser posterior a la fecha de inicio")
	}

	if s.CancellationDate != nil && (s.CancellationDate.Before(s.StartDate) || s.CancellationDate.After(s.EndDate)) {
		return errors.New("la fecha de cancelación debe estar entre la fecha de inicio y finalización")
	}

	return nil
}

// IsActive verifica si la suscripción está activa
func (s *UserSubscription) IsActive() bool {
	return s.Status == SubscriptionStatusActive && time.Now().Before(s.EndDate)
}

// HasExpired verifica si la suscripción ha expirado
func (s *UserSubscription) HasExpired() bool {
	return time.Now().After(s.EndDate)
}

// IsCancelled verifica si la suscripción ha sido cancelada
func (s *UserSubscription) IsCancelled() bool {
	return s.Status == SubscriptionStatusCancelled
}

// NeedsRenewal verifica si la suscripción necesita renovación
func (s *UserSubscription) NeedsRenewal() bool {
	if s.RenewalDate == nil {
		return false
	}
	return s.Status == SubscriptionStatusActive && time.Now().After(*s.RenewalDate)
}

// CreateSubscriptionRequest representa la solicitud para crear una suscripción
type CreateSubscriptionRequest struct {
	PlanID          string            `json:"plan_id" binding:"required"`
	StartDate       time.Time         `json:"start_date"`
	EndDate         time.Time         `json:"end_date"`
	PaymentMethodID *string           `json:"payment_method_id"`
	Metadata        map[string]string `json:"metadata"`
}

// UpdatePaymentMethodRequest representa la solicitud para actualizar un método de pago
type UpdatePaymentMethodRequest struct {
	PaymentMethodID string `json:"payment_method_id" binding:"required"`
}

// CancelSubscriptionRequest representa la solicitud para cancelar una suscripción
type CancelSubscriptionRequest struct {
	Reason string `json:"reason"`
}

// ChangePlanRequest representa la solicitud para cambiar de plan
type ChangePlanRequest struct {
	PlanID string `json:"plan_id" binding:"required"`
}

// SubscriptionResponse representa la respuesta de una suscripción
type SubscriptionResponse struct {
	ID                 string            `json:"id"`
	UserID             string            `json:"user_id"`
	PlanID             string            `json:"plan_id"`
	Status             string            `json:"status"`
	StartDate          time.Time         `json:"start_date"`
	EndDate            time.Time         `json:"end_date"`
	RenewalDate        *time.Time        `json:"renewal_date,omitempty"`
	CancellationDate   *time.Time        `json:"cancellation_date,omitempty"`
	LastPaymentDate    *time.Time        `json:"last_payment_date,omitempty"`
	NextPaymentAttempt *time.Time        `json:"next_payment_attempt,omitempty"`
	PaymentMethodID    *string           `json:"payment_method_id,omitempty"`
	Metadata           map[string]string `json:"metadata,omitempty"`
	CreatedAt          time.Time         `json:"created_at"`
	UpdatedAt          time.Time         `json:"updated_at"`
	IsActive           bool              `json:"is_active"`
}

// mapSubscriptionToResponse mapea una entidad de suscripción a una respuesta HTTP
func MapSubscriptionToResponse(subscription *UserSubscription) SubscriptionResponse {
	return SubscriptionResponse{
		ID:                 subscription.ID,
		UserID:             subscription.UserID,
		PlanID:             subscription.PlanID,
		Status:             string(subscription.Status),
		StartDate:          subscription.StartDate,
		EndDate:            subscription.EndDate,
		RenewalDate:        subscription.RenewalDate,
		CancellationDate:   subscription.CancellationDate,
		LastPaymentDate:    subscription.LastPaymentDate,
		NextPaymentAttempt: subscription.NextPaymentAttempt,
		PaymentMethodID:    subscription.PaymentMethodID,
		Metadata:           subscription.Metadata,
		CreatedAt:          subscription.CreatedAt,
		UpdatedAt:          subscription.UpdatedAt,
		IsActive:           subscription.IsActive(),
	}
}
