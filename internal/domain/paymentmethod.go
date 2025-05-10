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
