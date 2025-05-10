package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"

	"mi-app-backend/internal/domain"
)

// CurrencyRepository implementa el puerto app.CurrencyRepository
type CurrencyRepository struct {
	db *sql.DB
}

// NewCurrencyRepository crea una nueva instancia de CurrencyRepository
func NewCurrencyRepository(db *sql.DB) *CurrencyRepository {
	return &CurrencyRepository{
		db: db,
	}
}

// Create crea una nueva moneda en la base de datos
func (r *CurrencyRepository) Create(ctx context.Context, currency *domain.Currency) error {
	// Generar un ID único si no se proporciona uno
	if currency.ID == "" {
		currency.ID = uuid.New().String()
	}

	// Establecer fechas de creación y actualización
	now := time.Now()
	currency.CreatedAt = now
	currency.UpdatedAt = now

	// Consulta SQL para insertar una nueva moneda
	query := `
		INSERT INTO currencies (id, code, name, symbol, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`

	// Ejecutar la consulta
	err := r.db.QueryRowContext(
		ctx,
		query,
		currency.ID,
		currency.Code,
		currency.Name,
		currency.Symbol,
		currency.IsActive,
		currency.CreatedAt,
		currency.UpdatedAt,
	).Scan(&currency.ID)

	if err != nil {
		return fmt.Errorf("error al crear moneda: %w", err)
	}

	return nil
}

// GetByID obtiene una moneda por su ID
func (r *CurrencyRepository) GetByID(ctx context.Context, id string) (*domain.Currency, error) {
	query := `
		SELECT id, code, name, symbol, is_active, created_at, updated_at
		FROM currencies
		WHERE id = $1
	`

	var currency domain.Currency
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&currency.ID,
		&currency.Code,
		&currency.Name,
		&currency.Symbol,
		&currency.IsActive,
		&currency.CreatedAt,
		&currency.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("moneda no encontrada con id: %s", id)
		}
		return nil, fmt.Errorf("error al obtener moneda por id: %w", err)
	}

	return &currency, nil
}

// GetByCode obtiene una moneda por su código
func (r *CurrencyRepository) GetByCode(ctx context.Context, code string) (*domain.Currency, error) {
	query := `
		SELECT id, code, name, symbol, is_active, created_at, updated_at
		FROM currencies
		WHERE code = $1
	`

	var currency domain.Currency
	err := r.db.QueryRowContext(ctx, query, code).Scan(
		&currency.ID,
		&currency.Code,
		&currency.Name,
		&currency.Symbol,
		&currency.IsActive,
		&currency.CreatedAt,
		&currency.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("moneda no encontrada con código: %s", code)
		}
		return nil, fmt.Errorf("error al obtener moneda por código: %w", err)
	}

	return &currency, nil
}

// GetAll obtiene todas las monedas
func (r *CurrencyRepository) GetAll(ctx context.Context) ([]*domain.Currency, error) {
	query := `
		SELECT id, code, name, symbol, is_active, created_at, updated_at
		FROM currencies
		ORDER BY code
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error al obtener monedas: %w", err)
	}
	defer rows.Close()

	var currencies []*domain.Currency
	for rows.Next() {
		var currency domain.Currency
		if err := rows.Scan(
			&currency.ID,
			&currency.Code,
			&currency.Name,
			&currency.Symbol,
			&currency.IsActive,
			&currency.CreatedAt,
			&currency.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("error al escanear moneda: %w", err)
		}
		currencies = append(currencies, &currency)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error al iterar sobre monedas: %w", err)
	}

	return currencies, nil
}

// GetAllActive obtiene todas las monedas activas
func (r *CurrencyRepository) GetAllActive(ctx context.Context) ([]*domain.Currency, error) {
	query := `
		SELECT id, code, name, symbol, is_active, created_at, updated_at
		FROM currencies
		WHERE is_active = true
		ORDER BY code
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error al obtener monedas activas: %w", err)
	}
	defer rows.Close()

	var currencies []*domain.Currency
	for rows.Next() {
		var currency domain.Currency
		if err := rows.Scan(
			&currency.ID,
			&currency.Code,
			&currency.Name,
			&currency.Symbol,
			&currency.IsActive,
			&currency.CreatedAt,
			&currency.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("error al escanear moneda: %w", err)
		}
		currencies = append(currencies, &currency)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error al iterar sobre monedas activas: %w", err)
	}

	return currencies, nil
}

// Update actualiza una moneda existente
func (r *CurrencyRepository) Update(ctx context.Context, currency *domain.Currency) error {
	// Actualizar fecha de modificación
	currency.UpdatedAt = time.Now()

	query := `
		UPDATE currencies
		SET code = $2, name = $3, symbol = $4, is_active = $5, updated_at = $6
		WHERE id = $1
	`

	result, err := r.db.ExecContext(
		ctx,
		query,
		currency.ID,
		currency.Code,
		currency.Name,
		currency.Symbol,
		currency.IsActive,
		currency.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("error al actualizar moneda: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error al obtener filas afectadas: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("moneda no encontrada con id: %s", currency.ID)
	}

	return nil
}

// Delete elimina una moneda por su ID
func (r *CurrencyRepository) Delete(ctx context.Context, id string) error {
	query := `
		DELETE FROM currencies
		WHERE id = $1
	`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("error al eliminar moneda: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error al obtener filas afectadas: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("moneda no encontrada con id: %s", id)
	}

	return nil
}
