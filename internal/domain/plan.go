package domain

import (
	"errors"
	"time"
)

// PlanInterval define el intervalo de facturación de un plan
type PlanInterval string

const (
	// PlanIntervalMonthly representa un intervalo mensual
	PlanIntervalMonthly PlanInterval = "monthly"
	// PlanIntervalYearly representa un intervalo anual
	PlanIntervalYearly PlanInterval = "yearly"
)

// PlanFeature representa una característica disponible en un plan
type PlanFeature struct {
	Name        string `json:"name"`        // Nombre de la característica
	Description string `json:"description"` // Descripción de la característica
	Value       string `json:"value"`       // Valor de la característica (e.g. "Ilimitado", "5 GB", etc.)
	Included    bool   `json:"included"`    // Indica si la característica está incluida en el plan
}

// Plan representa un plan de suscripción
type Plan struct {
	ID          string        `json:"id"`
	Name        string        `json:"name"`        // Nombre del plan (e.g. "Gratis", "Pro")
	Description string        `json:"description"` // Descripción del plan
	Price       float64       `json:"price"`       // Precio del plan
	CurrencyID  string        `json:"currency_id"` // ID de la moneda
	Interval    PlanInterval  `json:"interval"`    // Intervalo de facturación
	Features    []PlanFeature `json:"features"`    // Características del plan
	IsActive    bool          `json:"is_active"`   // Indica si el plan está activo
	IsPublic    bool          `json:"is_public"`   // Indica si el plan es visible públicamente
	SortOrder   int           `json:"sort_order"`  // Orden de visualización
	CreatedAt   time.Time     `json:"created_at"`  // Fecha de creación
	UpdatedAt   time.Time     `json:"updated_at"`  // Fecha de actualización
}

// Validate valida que la entidad Plan tenga todos los campos requeridos
func (p *Plan) Validate() error {
	if p.Name == "" {
		return errors.New("el nombre del plan es obligatorio")
	}
	if p.Description == "" {
		return errors.New("la descripción del plan es obligatoria")
	}
	if p.Price < 0 {
		return errors.New("el precio no puede ser negativo")
	}
	if p.CurrencyID == "" {
		return errors.New("la moneda es obligatoria")
	}
	if p.Interval == "" {
		return errors.New("el intervalo de facturación es obligatorio")
	}
	if p.Interval != PlanIntervalMonthly && p.Interval != PlanIntervalYearly {
		return errors.New("el intervalo de facturación debe ser mensual o anual")
	}
	return nil
}

// PlanFeatureRequest representa una característica en la solicitud
type PlanFeatureRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Value       string `json:"value"`
	Included    bool   `json:"included"`
}

// CreatePlanRequest representa la solicitud para crear un nuevo plan
type CreatePlanRequest struct {
	Name        string               `json:"name" binding:"required"`
	Description string               `json:"description" binding:"required"`
	Price       float64              `json:"price"`
	CurrencyID  string               `json:"currency_id" binding:"required"`
	Interval    string               `json:"interval" binding:"required"`
	Features    []PlanFeatureRequest `json:"features"`
	IsActive    bool                 `json:"is_active"`
	IsPublic    bool                 `json:"is_public"`
	SortOrder   int                  `json:"sort_order"`
}

// UpdatePlanRequest representa la solicitud para actualizar un plan existente
type UpdatePlanRequest struct {
	Name        string               `json:"name" binding:"required"`
	Description string               `json:"description" binding:"required"`
	Price       float64              `json:"price"`
	CurrencyID  string               `json:"currency_id" binding:"required"`
	Interval    string               `json:"interval" binding:"required"`
	Features    []PlanFeatureRequest `json:"features"`
	IsActive    bool                 `json:"is_active"`
	IsPublic    bool                 `json:"is_public"`
	SortOrder   int                  `json:"sort_order"`
}

// PlanFeatureResponse representa una característica en la respuesta
type PlanFeatureResponse struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Value       string `json:"value"`
	Included    bool   `json:"included"`
}

// PlanResponse representa la respuesta con los datos de un plan
type PlanResponse struct {
	ID          string                `json:"id"`
	Name        string                `json:"name"`
	Description string                `json:"description"`
	Price       float64               `json:"price"`
	CurrencyID  string                `json:"currency_id"`
	Interval    string                `json:"interval"`
	Features    []PlanFeatureResponse `json:"features"`
	IsActive    bool                  `json:"is_active"`
	IsPublic    bool                  `json:"is_public"`
	SortOrder   int                   `json:"sort_order"`
	CreatedAt   string                `json:"created_at"`
	UpdatedAt   string                `json:"updated_at"`
}
