package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"

	"mi-app-backend/internal/domain"
)

// PlanRepository implementa el puerto app.PlanRepository
type PlanRepository struct {
	db *sql.DB
}

// NewPlanRepository crea una nueva instancia de PlanRepository
func NewPlanRepository(db *sql.DB) *PlanRepository {
	return &PlanRepository{
		db: db,
	}
}

// Create crea un nuevo plan en la base de datos
func (r *PlanRepository) Create(ctx context.Context, plan *domain.Plan) error {
	// Generar un ID único si no se proporciona uno
	if plan.ID == "" {
		plan.ID = uuid.New().String()
	}

	// Establecer fechas de creación y actualización
	now := time.Now()
	plan.CreatedAt = now
	plan.UpdatedAt = now

	// Convertir las características a JSON
	featuresJSON, err := json.Marshal(plan.Features)
	if err != nil {
		return fmt.Errorf("error al serializar características: %w", err)
	}

	// Consulta SQL para insertar un nuevo plan
	query := `
		INSERT INTO plans (
			id, name, description, price, currency_id, interval, 
			features, is_active, is_public, sort_order, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12
		) RETURNING id
	`

	// Ejecutar la consulta
	err = r.db.QueryRowContext(
		ctx,
		query,
		plan.ID,
		plan.Name,
		plan.Description,
		plan.Price,
		plan.CurrencyID,
		plan.Interval,
		featuresJSON,
		plan.IsActive,
		plan.IsPublic,
		plan.SortOrder,
		plan.CreatedAt,
		plan.UpdatedAt,
	).Scan(&plan.ID)

	if err != nil {
		return fmt.Errorf("error al crear plan: %w", err)
	}

	return nil
}

// GetByID obtiene un plan por su ID
func (r *PlanRepository) GetByID(ctx context.Context, id string) (*domain.Plan, error) {
	query := `
		SELECT 
			id, name, description, price, currency_id, interval, 
			features, is_active, is_public, sort_order, created_at, updated_at
		FROM plans
		WHERE id = $1
	`

	var (
		plan         domain.Plan
		featuresJSON []byte
		intervalStr  string
	)

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&plan.ID,
		&plan.Name,
		&plan.Description,
		&plan.Price,
		&plan.CurrencyID,
		&intervalStr,
		&featuresJSON,
		&plan.IsActive,
		&plan.IsPublic,
		&plan.SortOrder,
		&plan.CreatedAt,
		&plan.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("plan no encontrado con id: %s", id)
		}
		return nil, fmt.Errorf("error al obtener plan por id: %w", err)
	}

	// Convertir el intervalo
	plan.Interval = domain.PlanInterval(intervalStr)

	// Decodificar las características
	if len(featuresJSON) > 0 {
		if err := json.Unmarshal(featuresJSON, &plan.Features); err != nil {
			return nil, fmt.Errorf("error al deserializar características: %w", err)
		}
	}

	return &plan, nil
}

// GetAll obtiene todos los planes
func (r *PlanRepository) GetAll(ctx context.Context) ([]*domain.Plan, error) {
	query := `
		SELECT 
			id, name, description, price, currency_id, interval, 
			features, is_active, is_public, sort_order, created_at, updated_at
		FROM plans
		ORDER BY sort_order ASC, name ASC
	`

	return r.queryPlans(ctx, query)
}

// GetAllPublic obtiene todos los planes públicos
func (r *PlanRepository) GetAllPublic(ctx context.Context) ([]*domain.Plan, error) {
	query := `
		SELECT 
			id, name, description, price, currency_id, interval, 
			features, is_active, is_public, sort_order, created_at, updated_at
		FROM plans
		WHERE is_public = true
		ORDER BY sort_order ASC, name ASC
	`

	return r.queryPlans(ctx, query)
}

// GetAllActive obtiene todos los planes activos
func (r *PlanRepository) GetAllActive(ctx context.Context) ([]*domain.Plan, error) {
	query := `
		SELECT 
			id, name, description, price, currency_id, interval, 
			features, is_active, is_public, sort_order, created_at, updated_at
		FROM plans
		WHERE is_active = true
		ORDER BY sort_order ASC, name ASC
	`

	return r.queryPlans(ctx, query)
}

// queryPlans ejecuta una consulta y devuelve una lista de planes
func (r *PlanRepository) queryPlans(ctx context.Context, query string, args ...interface{}) ([]*domain.Plan, error) {
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("error al ejecutar consulta: %w", err)
	}
	defer rows.Close()

	var plans []*domain.Plan
	for rows.Next() {
		var (
			plan         domain.Plan
			featuresJSON []byte
			intervalStr  string
		)

		if err := rows.Scan(
			&plan.ID,
			&plan.Name,
			&plan.Description,
			&plan.Price,
			&plan.CurrencyID,
			&intervalStr,
			&featuresJSON,
			&plan.IsActive,
			&plan.IsPublic,
			&plan.SortOrder,
			&plan.CreatedAt,
			&plan.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("error al escanear plan: %w", err)
		}

		// Convertir el intervalo
		plan.Interval = domain.PlanInterval(intervalStr)

		// Decodificar las características
		if len(featuresJSON) > 0 {
			if err := json.Unmarshal(featuresJSON, &plan.Features); err != nil {
				return nil, fmt.Errorf("error al deserializar características: %w", err)
			}
		}

		plans = append(plans, &plan)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error al iterar sobre planes: %w", err)
	}

	return plans, nil
}

// Update actualiza un plan existente
func (r *PlanRepository) Update(ctx context.Context, plan *domain.Plan) error {
	// Actualizar fecha de modificación
	plan.UpdatedAt = time.Now()

	// Convertir las características a JSON
	featuresJSON, err := json.Marshal(plan.Features)
	if err != nil {
		return fmt.Errorf("error al serializar características: %w", err)
	}

	query := `
		UPDATE plans
		SET 
			name = $2,
			description = $3,
			price = $4,
			currency_id = $5,
			interval = $6,
			features = $7,
			is_active = $8,
			is_public = $9,
			sort_order = $10,
			updated_at = $11
		WHERE id = $1
	`

	result, err := r.db.ExecContext(
		ctx,
		query,
		plan.ID,
		plan.Name,
		plan.Description,
		plan.Price,
		plan.CurrencyID,
		plan.Interval,
		featuresJSON,
		plan.IsActive,
		plan.IsPublic,
		plan.SortOrder,
		plan.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("error al actualizar plan: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error al obtener filas afectadas: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("plan no encontrado con id: %s", plan.ID)
	}

	return nil
}

// Delete elimina un plan por su ID
func (r *PlanRepository) Delete(ctx context.Context, id string) error {
	query := `
		DELETE FROM plans
		WHERE id = $1
	`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("error al eliminar plan: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error al obtener filas afectadas: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("plan no encontrado con id: %s", id)
	}

	return nil
}
