package domain

import (
	"time"
)

// Category representa una categoría en el sistema
type Category struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Color       string    `json:"color"`
	Icon        string    `json:"icon"`
	UserID      string    `json:"user_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Validate valida que los campos obligatorios estén presentes
func (c *Category) Validate() error {
	if c.Name == "" {
		return ErrEmptyName
	}
	if c.UserID == "" {
		return ErrEmptyUserID
	}
	return nil
}

// createCategoryRequest represents the create category request
type CreateCategoryRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
	Color       string `json:"color"`
}

// updateCategoryRequest represents the update category request
type UpdateCategoryRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
	Color       string `json:"color"`
}
