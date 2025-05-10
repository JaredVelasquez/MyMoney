package domain

import (
	"errors"
	"time"
)

// Currency representa una moneda en el sistema
type Currency struct {
	ID        string    `json:"id"`
	Code      string    `json:"code"`       // Código ISO de la moneda (USD, EUR, etc.)
	Name      string    `json:"name"`       // Nombre de la moneda
	Symbol    string    `json:"symbol"`     // Símbolo de la moneda ($, €, etc.)
	IsActive  bool      `json:"is_active"`  // Indica si la moneda está activa
	CreatedAt time.Time `json:"created_at"` // Fecha de creación
	UpdatedAt time.Time `json:"updated_at"` // Fecha de actualización
}

// Validate valida que la entidad Currency tenga todos los campos requeridos
func (c *Currency) Validate() error {
	if c.Code == "" {
		return errors.New("el código de la moneda es obligatorio")
	}
	if c.Name == "" {
		return errors.New("el nombre de la moneda es obligatorio")
	}
	if c.Symbol == "" {
		return errors.New("el símbolo de la moneda es obligatorio")
	}
	return nil
}

// CreateCurrencyRequest representa la solicitud para crear una nueva moneda
type CreateCurrencyRequest struct {
	Code     string `json:"code" binding:"required"`
	Name     string `json:"name" binding:"required"`
	Symbol   string `json:"symbol" binding:"required"`
	IsActive bool   `json:"is_active"`
}

// UpdateCurrencyRequest representa la solicitud para actualizar una moneda existente
type UpdateCurrencyRequest struct {
	Code     string `json:"code" binding:"required"`
	Name     string `json:"name" binding:"required"`
	Symbol   string `json:"symbol" binding:"required"`
	IsActive bool   `json:"is_active"`
}

// CurrencyResponse representa la respuesta con los datos de una moneda
type CurrencyResponse struct {
	ID        string `json:"id"`
	Code      string `json:"code"`
	Name      string `json:"name"`
	Symbol    string `json:"symbol"`
	IsActive  bool   `json:"is_active"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
