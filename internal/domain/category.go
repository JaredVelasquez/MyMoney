package domain

import (
	"errors"
	"time"
)

// CategoryType representa el tipo de categoría (ingreso o gasto)
type CategoryType string

const (
	// CategoryTypeIncome representa una categoría de ingresos
	CategoryTypeIncome CategoryType = "INCOME"
	// CategoryTypeExpense representa una categoría de gastos
	CategoryTypeExpense CategoryType = "EXPENSE"
)

// Category representa una categoría en el sistema
type Category struct {
	ID          string       `json:"id"`
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Type        CategoryType `json:"type"`
	Color       string       `json:"color"`
	Icon        string       `json:"icon"`
	UserID      string       `json:"user_id"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
}

// Validate valida que los campos obligatorios estén presentes
func (c *Category) Validate() error {
	if c.Name == "" {
		return ErrEmptyName
	}
	if c.Type == "" {
		return errors.New("el tipo de categoría no puede estar vacío")
	}
	if c.Type != CategoryTypeIncome && c.Type != CategoryTypeExpense {
		return errors.New("tipo de categoría inválido")
	}
	if c.UserID == "" {
		return ErrEmptyUserID
	}
	return nil
}

// createCategoryRequest represents the create category request
type CreateCategoryRequest struct {
	Name        string       `json:"name" binding:"required"`
	Description string       `json:"description"`
	Type        CategoryType `json:"type" binding:"required"`
	Icon        string       `json:"icon"`
	Color       string       `json:"color"`
}

// updateCategoryRequest represents the update category request
type UpdateCategoryRequest struct {
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Type        CategoryType `json:"type"`
	Icon        string       `json:"icon"`
	Color       string       `json:"color"`
}
