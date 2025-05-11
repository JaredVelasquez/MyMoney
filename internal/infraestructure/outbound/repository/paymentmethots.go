package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"MyMoneyBackend/internal/domain"

	"github.com/google/uuid"
)

// PaymentMethodRepository implementa la interfaz app.PaymentMethodRepository
type PaymentMethodRepository struct {
	db *sql.DB
}

// NewPaymentMethodRepository crea un nuevo repositorio de métodos de pago
func NewPaymentMethodRepository(db *sql.DB) *PaymentMethodRepository {
	return &PaymentMethodRepository{
		db: db,
	}
}

// Create crea un nuevo método de pago en la base de datos
func (r *PaymentMethodRepository) Create(ctx context.Context, paymentMethod *domain.PaymentMethod) error {
	if paymentMethod.ID == "" {
		paymentMethod.ID = uuid.New().String()
	}

	now := time.Now()
	if paymentMethod.CreatedAt.IsZero() {
		paymentMethod.CreatedAt = now
	}
	paymentMethod.UpdatedAt = now

	query := `
		INSERT INTO payment_methods (id, name, description, is_active, user_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		paymentMethod.ID,
		paymentMethod.Name,
		paymentMethod.Description,
		paymentMethod.IsActive,
		paymentMethod.UserID,
		paymentMethod.CreatedAt,
		paymentMethod.UpdatedAt,
	)

	return err
}

// GetByID obtiene un método de pago por su ID
func (r *PaymentMethodRepository) GetByID(ctx context.Context, id string) (*domain.PaymentMethod, error) {
	query := `
		SELECT id, name, description, is_active, user_id, created_at, updated_at
		FROM payment_methods
		WHERE id = $1
	`

	var paymentMethod domain.PaymentMethod
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&paymentMethod.ID,
		&paymentMethod.Name,
		&paymentMethod.Description,
		&paymentMethod.IsActive,
		&paymentMethod.UserID,
		&paymentMethod.CreatedAt,
		&paymentMethod.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // No se encontró el método de pago
		}
		return nil, err
	}

	return &paymentMethod, nil
}

// GetByUserID obtiene todos los métodos de pago de un usuario
func (r *PaymentMethodRepository) GetByUserID(ctx context.Context, userID string) ([]*domain.PaymentMethod, error) {
	query := `
		SELECT id, name, description, is_active, user_id, created_at, updated_at
		FROM payment_methods
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var paymentMethods []*domain.PaymentMethod
	for rows.Next() {
		var paymentMethod domain.PaymentMethod
		if err := rows.Scan(
			&paymentMethod.ID,
			&paymentMethod.Name,
			&paymentMethod.Description,
			&paymentMethod.IsActive,
			&paymentMethod.UserID,
			&paymentMethod.CreatedAt,
			&paymentMethod.UpdatedAt,
		); err != nil {
			return nil, err
		}
		paymentMethods = append(paymentMethods, &paymentMethod)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return paymentMethods, nil
}

// Update actualiza un método de pago existente
func (r *PaymentMethodRepository) Update(ctx context.Context, paymentMethod *domain.PaymentMethod) error {
	paymentMethod.UpdatedAt = time.Now()

	query := `
		UPDATE payment_methods
		SET name = $1, description = $2, is_active = $3, updated_at = $4
		WHERE id = $5
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		paymentMethod.Name,
		paymentMethod.Description,
		paymentMethod.IsActive,
		paymentMethod.UpdatedAt,
		paymentMethod.ID,
	)

	return err
}

// Delete elimina un método de pago
func (r *PaymentMethodRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM payment_methods WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
