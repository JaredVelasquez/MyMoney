package domain

import (
	"time"
)

// PaymentMethod representa un método de pago en el sistema
type PaymentMethod struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	IsActive    bool      `json:"is_active"`
	UserID      string    `json:"user_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Validate valida que los campos obligatorios estén presentes
func (p *PaymentMethod) Validate() error {
	if p.Name == "" {
		return ErrEmptyName
	}
	if p.UserID == "" {
		return ErrEmptyUserID
	}
	return nil
}

// CreatePaymentMethodRequest representa la solicitud para crear un método de pago
type CreatePaymentMethodRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

// UpdatePaymentMethodRequest representa la solicitud para actualizar un método de pago
type UpdatePaymentMethodRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	IsActive    bool   `json:"is_active"`
}

// PaymentMethodResponse representa la respuesta para un método de pago
type PaymentMethodResponse struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
