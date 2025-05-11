package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"MyMoneyBackend/internal/domain"
)

// TransactionRepository implements domain.TransactionRepository for PostgreSQL
type TransactionRepository struct {
	db *sql.DB
}

// NewTransactionRepository creates a new TransactionRepository
func NewTransactionRepository(db *sql.DB) *TransactionRepository {
	return &TransactionRepository{
		db: db,
	}
}

// Create inserts a new transaction into the database
func (r *TransactionRepository) Create(ctx context.Context, transaction *domain.Transaction) error {
	query := `
		INSERT INTO transactions (
			id, user_id, amount, description, category_id, type, payment_method_id, 
			currency_id, date, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
		)
	`

	now := time.Now()
	_, err := r.db.ExecContext(
		ctx,
		query,
		transaction.ID,
		transaction.UserID,
		transaction.Amount,
		transaction.Description,
		transaction.CategoryID,
		transaction.Type,
		transaction.PaymentMethodID,
		transaction.CurrencyID,
		transaction.Date,
		now,
		now,
	)

	if err != nil {
		return fmt.Errorf("error creating transaction: %w", err)
	}

	transaction.CreatedAt = now
	transaction.UpdatedAt = now

	return nil
}

// GetByID retrieves a transaction by ID
func (r *TransactionRepository) GetByID(ctx context.Context, id string) (*domain.Transaction, error) {
	query := `
		SELECT 
			id, user_id, amount, description, category_id, type,
			payment_method_id, currency_id, date, created_at, updated_at
		FROM transactions
		WHERE id = $1
	`

	var transaction domain.Transaction
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&transaction.ID,
		&transaction.UserID,
		&transaction.Amount,
		&transaction.Description,
		&transaction.CategoryID,
		&transaction.Type,
		&transaction.PaymentMethodID,
		&transaction.CurrencyID,
		&transaction.Date,
		&transaction.CreatedAt,
		&transaction.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("transaction not found: %w", err)
		}
		return nil, fmt.Errorf("error getting transaction: %w", err)
	}

	return &transaction, nil
}

// GetByUserID retrieves all transactions for a user
func (r *TransactionRepository) GetByUserID(ctx context.Context, userID string) ([]*domain.Transaction, error) {
	query := `
		SELECT 
			id, user_id, amount, description, category_id, type,
			payment_method_id, currency_id, date, created_at, updated_at
		FROM transactions
		WHERE user_id = $1
		ORDER BY date DESC
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("error querying transactions: %w", err)
	}
	defer rows.Close()

	return r.scanTransactions(rows)
}

// GetByCategoryID retrieves all transactions for a category
func (r *TransactionRepository) GetByCategoryID(ctx context.Context, categoryID string) ([]*domain.Transaction, error) {
	query := `
		SELECT 
			id, user_id, amount, description, category_id, type,
			payment_method_id, currency_id, date, created_at, updated_at
		FROM transactions
		WHERE category_id = $1
		ORDER BY date DESC
	`

	rows, err := r.db.QueryContext(ctx, query, categoryID)
	if err != nil {
		return nil, fmt.Errorf("error querying transactions by category: %w", err)
	}
	defer rows.Close()

	return r.scanTransactions(rows)
}

// GetByDateRange retrieves all transactions within a date range
func (r *TransactionRepository) GetByDateRange(ctx context.Context, userID string, startDate, endDate time.Time) ([]*domain.Transaction, error) {
	query := `
		SELECT 
			id, user_id, amount, description, category_id, type,
			payment_method_id, currency_id, date, created_at, updated_at
		FROM transactions
		WHERE user_id = $1 AND date BETWEEN $2 AND $3
		ORDER BY date DESC
	`

	rows, err := r.db.QueryContext(ctx, query, userID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("error querying transactions by date range: %w", err)
	}
	defer rows.Close()

	return r.scanTransactions(rows)
}

// Update updates a transaction's information
func (r *TransactionRepository) Update(ctx context.Context, transaction *domain.Transaction) error {
	query := `
		UPDATE transactions
		SET amount = $1, description = $2, category_id = $3, type = $4,
			payment_method_id = $5, currency_id = $6, date = $7, updated_at = $8
		WHERE id = $9 AND user_id = $10
	`

	now := time.Now()
	result, err := r.db.ExecContext(
		ctx,
		query,
		transaction.Amount,
		transaction.Description,
		transaction.CategoryID,
		transaction.Type,
		transaction.PaymentMethodID,
		transaction.CurrencyID,
		transaction.Date,
		now,
		transaction.ID,
		transaction.UserID,
	)

	if err != nil {
		return fmt.Errorf("error updating transaction: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("transaction not found or does not belong to the user")
	}

	transaction.UpdatedAt = now

	return nil
}

// Delete removes a transaction from the database
func (r *TransactionRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM transactions WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("error deleting transaction: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("transaction not found")
	}

	return nil
}

// scanTransactions is a helper method to scan multiple transaction rows
func (r *TransactionRepository) scanTransactions(rows *sql.Rows) ([]*domain.Transaction, error) {
	var transactions []*domain.Transaction
	for rows.Next() {
		var transaction domain.Transaction
		err := rows.Scan(
			&transaction.ID,
			&transaction.UserID,
			&transaction.Amount,
			&transaction.Description,
			&transaction.CategoryID,
			&transaction.Type,
			&transaction.PaymentMethodID,
			&transaction.CurrencyID,
			&transaction.Date,
			&transaction.CreatedAt,
			&transaction.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning transaction row: %w", err)
		}
		transactions = append(transactions, &transaction)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating transaction rows: %w", err)
	}

	return transactions, nil
}
