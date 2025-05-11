package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"

	"MyMoneyBackend/internal/domain"
)

// UserSubscriptionRepository implementa el puerto app.UserSubscriptionRepository
type UserSubscriptionRepository struct {
	db *sql.DB
}

// NewUserSubscriptionRepository crea una nueva instancia de UserSubscriptionRepository
func NewUserSubscriptionRepository(db *sql.DB) *UserSubscriptionRepository {
	return &UserSubscriptionRepository{
		db: db,
	}
}

// Create crea una nueva suscripción en la base de datos
func (r *UserSubscriptionRepository) Create(ctx context.Context, subscription *domain.UserSubscription) error {
	// Generar un ID único si no se proporciona uno
	if subscription.ID == "" {
		subscription.ID = uuid.New().String()
	}

	// Establecer fechas de creación y actualización
	now := time.Now()
	subscription.CreatedAt = now
	subscription.UpdatedAt = now

	// Convertir los metadatos a JSON
	var metadataJSON []byte
	var err error
	if subscription.Metadata != nil {
		metadataJSON, err = json.Marshal(subscription.Metadata)
		if err != nil {
			return fmt.Errorf("error al serializar metadatos: %w", err)
		}
	}

	// Consulta SQL para insertar una nueva suscripción
	query := `
		INSERT INTO user_subscriptions (
			id, user_id, plan_id, status, start_date, end_date, 
			renewal_date, cancellation_date, last_payment_date, 
			next_payment_attempt, payment_method_id, metadata, 
			created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14
		) RETURNING id
	`

	// Ejecutar la consulta
	err = r.db.QueryRowContext(
		ctx,
		query,
		subscription.ID,
		subscription.UserID,
		subscription.PlanID,
		subscription.Status,
		subscription.StartDate,
		subscription.EndDate,
		subscription.RenewalDate,
		subscription.CancellationDate,
		subscription.LastPaymentDate,
		subscription.NextPaymentAttempt,
		subscription.PaymentMethodID,
		metadataJSON,
		subscription.CreatedAt,
		subscription.UpdatedAt,
	).Scan(&subscription.ID)

	if err != nil {
		return fmt.Errorf("error al crear suscripción: %w", err)
	}

	return nil
}

// GetByID obtiene una suscripción por su ID
func (r *UserSubscriptionRepository) GetByID(ctx context.Context, id string) (*domain.UserSubscription, error) {
	query := `
		SELECT 
			id, user_id, plan_id, status, start_date, end_date, 
			renewal_date, cancellation_date, last_payment_date, 
			next_payment_attempt, payment_method_id, metadata, 
			created_at, updated_at
		FROM user_subscriptions
		WHERE id = $1
	`

	var (
		subscription    domain.UserSubscription
		metadataJSON    []byte
		renewalDate     sql.NullTime
		cancelDate      sql.NullTime
		lastPayDate     sql.NullTime
		nextPayAttempt  sql.NullTime
		paymentMethodID sql.NullString
	)

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&subscription.ID,
		&subscription.UserID,
		&subscription.PlanID,
		&subscription.Status,
		&subscription.StartDate,
		&subscription.EndDate,
		&renewalDate,
		&cancelDate,
		&lastPayDate,
		&nextPayAttempt,
		&paymentMethodID,
		&metadataJSON,
		&subscription.CreatedAt,
		&subscription.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("suscripción no encontrada con id: %s", id)
		}
		return nil, fmt.Errorf("error al obtener suscripción por id: %w", err)
	}

	// Convertir campos nulos
	if renewalDate.Valid {
		subscription.RenewalDate = &renewalDate.Time
	}
	if cancelDate.Valid {
		subscription.CancellationDate = &cancelDate.Time
	}
	if lastPayDate.Valid {
		subscription.LastPaymentDate = &lastPayDate.Time
	}
	if nextPayAttempt.Valid {
		subscription.NextPaymentAttempt = &nextPayAttempt.Time
	}
	if paymentMethodID.Valid {
		subscription.PaymentMethodID = &paymentMethodID.String
	}

	// Decodificar metadatos
	if len(metadataJSON) > 0 {
		var metadata map[string]string
		if err := json.Unmarshal(metadataJSON, &metadata); err != nil {
			return nil, fmt.Errorf("error al deserializar metadatos: %w", err)
		}
		subscription.Metadata = metadata
	}

	return &subscription, nil
}

// GetActiveByUserID obtiene la suscripción activa de un usuario
func (r *UserSubscriptionRepository) GetActiveByUserID(ctx context.Context, userID string) (*domain.UserSubscription, error) {
	query := `
		SELECT 
			id, user_id, plan_id, status, start_date, end_date, 
			renewal_date, cancellation_date, last_payment_date, 
			next_payment_attempt, payment_method_id, metadata, 
			created_at, updated_at
		FROM user_subscriptions
		WHERE user_id = $1 AND status = $2 AND end_date > $3
		ORDER BY end_date DESC
		LIMIT 1
	`

	var (
		subscription    domain.UserSubscription
		metadataJSON    []byte
		renewalDate     sql.NullTime
		cancelDate      sql.NullTime
		lastPayDate     sql.NullTime
		nextPayAttempt  sql.NullTime
		paymentMethodID sql.NullString
	)

	err := r.db.QueryRowContext(
		ctx,
		query,
		userID,
		domain.SubscriptionStatusActive,
		time.Now(),
	).Scan(
		&subscription.ID,
		&subscription.UserID,
		&subscription.PlanID,
		&subscription.Status,
		&subscription.StartDate,
		&subscription.EndDate,
		&renewalDate,
		&cancelDate,
		&lastPayDate,
		&nextPayAttempt,
		&paymentMethodID,
		&metadataJSON,
		&subscription.CreatedAt,
		&subscription.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No hay suscripción activa
		}
		return nil, fmt.Errorf("error al obtener suscripción activa del usuario: %w", err)
	}

	// Convertir campos nulos
	if renewalDate.Valid {
		subscription.RenewalDate = &renewalDate.Time
	}
	if cancelDate.Valid {
		subscription.CancellationDate = &cancelDate.Time
	}
	if lastPayDate.Valid {
		subscription.LastPaymentDate = &lastPayDate.Time
	}
	if nextPayAttempt.Valid {
		subscription.NextPaymentAttempt = &nextPayAttempt.Time
	}
	if paymentMethodID.Valid {
		subscription.PaymentMethodID = &paymentMethodID.String
	}

	// Decodificar metadatos
	if len(metadataJSON) > 0 {
		var metadata map[string]string
		if err := json.Unmarshal(metadataJSON, &metadata); err != nil {
			return nil, fmt.Errorf("error al deserializar metadatos: %w", err)
		}
		subscription.Metadata = metadata
	}

	return &subscription, nil
}

// GetAllByUserID obtiene todas las suscripciones de un usuario
func (r *UserSubscriptionRepository) GetAllByUserID(ctx context.Context, userID string) ([]*domain.UserSubscription, error) {
	query := `
		SELECT 
			id, user_id, plan_id, status, start_date, end_date, 
			renewal_date, cancellation_date, last_payment_date, 
			next_payment_attempt, payment_method_id, metadata, 
			created_at, updated_at
		FROM user_subscriptions
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	return r.querySubscriptions(ctx, query, userID)
}

// GetByStatus obtiene suscripciones por estado
func (r *UserSubscriptionRepository) GetByStatus(ctx context.Context, status domain.SubscriptionStatus) ([]*domain.UserSubscription, error) {
	query := `
		SELECT 
			id, user_id, plan_id, status, start_date, end_date, 
			renewal_date, cancellation_date, last_payment_date, 
			next_payment_attempt, payment_method_id, metadata, 
			created_at, updated_at
		FROM user_subscriptions
		WHERE status = $1
		ORDER BY created_at DESC
	`

	return r.querySubscriptions(ctx, query, status)
}

// GetExpiringSubscriptions obtiene suscripciones que expirarán pronto
func (r *UserSubscriptionRepository) GetExpiringSubscriptions(ctx context.Context, beforeDate time.Time) ([]*domain.UserSubscription, error) {
	query := `
		SELECT 
			id, user_id, plan_id, status, start_date, end_date, 
			renewal_date, cancellation_date, last_payment_date, 
			next_payment_attempt, payment_method_id, metadata, 
			created_at, updated_at
		FROM user_subscriptions
		WHERE status = $1 AND end_date <= $2 AND end_date > $3
		ORDER BY end_date ASC
	`

	return r.querySubscriptions(ctx, query, domain.SubscriptionStatusActive, beforeDate, time.Now())
}

// GetPendingRenewals obtiene suscripciones pendientes de renovación
func (r *UserSubscriptionRepository) GetPendingRenewals(ctx context.Context, beforeDate time.Time) ([]*domain.UserSubscription, error) {
	query := `
		SELECT 
			id, user_id, plan_id, status, start_date, end_date, 
			renewal_date, cancellation_date, last_payment_date, 
			next_payment_attempt, payment_method_id, metadata, 
			created_at, updated_at
		FROM user_subscriptions
		WHERE status = $1 AND renewal_date IS NOT NULL AND renewal_date <= $2
		ORDER BY renewal_date ASC
	`

	return r.querySubscriptions(ctx, query, domain.SubscriptionStatusActive, beforeDate)
}

// querySubscriptions ejecuta una consulta y devuelve una lista de suscripciones
func (r *UserSubscriptionRepository) querySubscriptions(ctx context.Context, query string, args ...interface{}) ([]*domain.UserSubscription, error) {
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("error al ejecutar consulta: %w", err)
	}
	defer rows.Close()

	var subscriptions []*domain.UserSubscription
	for rows.Next() {
		var (
			subscription    domain.UserSubscription
			metadataJSON    []byte
			renewalDate     sql.NullTime
			cancelDate      sql.NullTime
			lastPayDate     sql.NullTime
			nextPayAttempt  sql.NullTime
			paymentMethodID sql.NullString
		)

		if err := rows.Scan(
			&subscription.ID,
			&subscription.UserID,
			&subscription.PlanID,
			&subscription.Status,
			&subscription.StartDate,
			&subscription.EndDate,
			&renewalDate,
			&cancelDate,
			&lastPayDate,
			&nextPayAttempt,
			&paymentMethodID,
			&metadataJSON,
			&subscription.CreatedAt,
			&subscription.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("error al escanear suscripción: %w", err)
		}

		// Convertir campos nulos
		if renewalDate.Valid {
			subscription.RenewalDate = &renewalDate.Time
		}
		if cancelDate.Valid {
			subscription.CancellationDate = &cancelDate.Time
		}
		if lastPayDate.Valid {
			subscription.LastPaymentDate = &lastPayDate.Time
		}
		if nextPayAttempt.Valid {
			subscription.NextPaymentAttempt = &nextPayAttempt.Time
		}
		if paymentMethodID.Valid {
			subscription.PaymentMethodID = &paymentMethodID.String
		}

		// Decodificar metadatos
		if len(metadataJSON) > 0 {
			var metadata map[string]string
			if err := json.Unmarshal(metadataJSON, &metadata); err != nil {
				return nil, fmt.Errorf("error al deserializar metadatos: %w", err)
			}
			subscription.Metadata = metadata
		}

		subscriptions = append(subscriptions, &subscription)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error al iterar sobre suscripciones: %w", err)
	}

	return subscriptions, nil
}

// Update actualiza una suscripción existente
func (r *UserSubscriptionRepository) Update(ctx context.Context, subscription *domain.UserSubscription) error {
	// Actualizar fecha de modificación
	subscription.UpdatedAt = time.Now()

	// Convertir los metadatos a JSON
	var metadataJSON []byte
	var err error
	if subscription.Metadata != nil {
		metadataJSON, err = json.Marshal(subscription.Metadata)
		if err != nil {
			return fmt.Errorf("error al serializar metadatos: %w", err)
		}
	}

	query := `
		UPDATE user_subscriptions
		SET 
			user_id = $2,
			plan_id = $3,
			status = $4,
			start_date = $5,
			end_date = $6,
			renewal_date = $7,
			cancellation_date = $8,
			last_payment_date = $9,
			next_payment_attempt = $10,
			payment_method_id = $11,
			metadata = $12,
			updated_at = $13
		WHERE id = $1
	`

	result, err := r.db.ExecContext(
		ctx,
		query,
		subscription.ID,
		subscription.UserID,
		subscription.PlanID,
		subscription.Status,
		subscription.StartDate,
		subscription.EndDate,
		subscription.RenewalDate,
		subscription.CancellationDate,
		subscription.LastPaymentDate,
		subscription.NextPaymentAttempt,
		subscription.PaymentMethodID,
		metadataJSON,
		subscription.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("error al actualizar suscripción: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error al obtener filas afectadas: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("suscripción no encontrada con id: %s", subscription.ID)
	}

	return nil
}

// UpdateStatus actualiza el estado de una suscripción
func (r *UserSubscriptionRepository) UpdateStatus(ctx context.Context, id string, status domain.SubscriptionStatus) error {
	query := `
		UPDATE user_subscriptions
		SET 
			status = $2,
			updated_at = $3
		WHERE id = $1
	`

	result, err := r.db.ExecContext(
		ctx,
		query,
		id,
		status,
		time.Now(),
	)

	if err != nil {
		return fmt.Errorf("error al actualizar estado de suscripción: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error al obtener filas afectadas: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("suscripción no encontrada con id: %s", id)
	}

	return nil
}

// CancelSubscription cancela una suscripción
func (r *UserSubscriptionRepository) CancelSubscription(ctx context.Context, id string, cancellationDate time.Time) error {
	query := `
		UPDATE user_subscriptions
		SET 
			status = $2,
			cancellation_date = $3,
			updated_at = $4
		WHERE id = $1
	`

	result, err := r.db.ExecContext(
		ctx,
		query,
		id,
		domain.SubscriptionStatusCancelled,
		cancellationDate,
		time.Now(),
	)

	if err != nil {
		return fmt.Errorf("error al cancelar suscripción: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error al obtener filas afectadas: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("suscripción no encontrada con id: %s", id)
	}

	return nil
}

// Delete elimina una suscripción por su ID
func (r *UserSubscriptionRepository) Delete(ctx context.Context, id string) error {
	query := `
		DELETE FROM user_subscriptions
		WHERE id = $1
	`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("error al eliminar suscripción: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error al obtener filas afectadas: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("suscripción no encontrada con id: %s", id)
	}

	return nil
}
