package domain

import (
	"errors"
	"time"
)

// TransactionType representa el tipo de transacción (ingreso o gasto)
type TransactionType string

const (
	TransactionTypeIncome  TransactionType = "INCOME"
	TransactionTypeExpense TransactionType = "EXPENSE"
)

// Transaction representa una transacción financiera en el sistema
type Transaction struct {
	ID              string          `json:"id"`
	Amount          float64         `json:"amount"`
	Description     string          `json:"description"`
	Date            time.Time       `json:"date"`
	CategoryID      string          `json:"category_id"`
	PaymentMethodID string          `json:"payment_method_id"`
	UserID          string          `json:"user_id"`
	CreatedAt       time.Time       `json:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at"`
	Type            TransactionType `json:"type"`
	CurrencyID      string          `json:"currency_id"`
}

// Validate valida que los campos obligatorios estén presentes
func (t *Transaction) Validate() error {
	if t.Amount <= 0 {
		return ErrInvalidAmount
	}
	if t.CategoryID == "" {
		return ErrEmptyCategoryID
	}
	if t.UserID == "" {
		return ErrEmptyUserID
	}
	if t.Date.IsZero() {
		return errors.New("la fecha no puede estar vacía")
	}
	return nil
}

// createTransactionRequest represents the create transaction request
type CreateTransactionRequest struct {
	Amount          float64         `json:"amount" binding:"required,gt=0"`
	Description     string          `json:"description"`
	CategoryID      string          `json:"category_id" binding:"required"`
	Type            TransactionType `json:"type" binding:"required"`
	PaymentMethodID string          `json:"payment_method_id"`
	CurrencyID      string          `json:"currency_id" binding:"required"`
	Date            time.Time       `json:"date" binding:"required"`
}

// updateTransactionRequest represents the update transaction request
type UpdateTransactionRequest struct {
	Amount          float64         `json:"amount"`
	Description     string          `json:"description"`
	CategoryID      string          `json:"category_id"`
	Type            TransactionType `json:"type"`
	PaymentMethodID string          `json:"payment_method_id"`
	CurrencyID      string          `json:"currency_id"`
	Date            time.Time       `json:"date"`
}

// dateRangeRequest represents a date range query
type DateRangeRequest struct {
	StartDate string `form:"start_date" binding:"required"`
	EndDate   string `form:"end_date" binding:"required"`
}
