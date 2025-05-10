package api

import (
	"net/http"
	"time"

	transaction "mi-app-backend/internal/application/transaction"
	"mi-app-backend/internal/domain"

	"github.com/gin-gonic/gin"
)

// TransactionHandler handles transaction-related HTTP requests
type TransactionHandler struct {
	transactionService *transaction.Service
}

// NewTransactionHandler creates a new TransactionHandler
func NewTransactionHandler(transactionService *transaction.Service) *TransactionHandler {
	return &TransactionHandler{
		transactionService: transactionService,
	}
}

// CreateTransaction handles transaction creation
func (h *TransactionHandler) CreateTransaction(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req domain.CreateTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	transaction, err := h.transactionService.CreateTransaction(
		c.Request.Context(),
		req.Amount,
		req.Description,
		req.Date,
		req.CategoryID,
		req.PaymentMethodID,
		userID,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, transaction)
}

// GetUserTransactions returns all transactions for the current user
func (h *TransactionHandler) GetUserTransactions(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	transactions, err := h.transactionService.GetTransactionsByUserID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error retrieving transactions"})
		return
	}

	c.JSON(http.StatusOK, transactions)
}

// GetTransaction returns a specific transaction
func (h *TransactionHandler) GetTransaction(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	transactionID := c.Param("id")
	if transactionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "transaction ID is required"})
		return
	}

	transaction, err := h.transactionService.GetTransactionByID(c.Request.Context(), transactionID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "transaction not found"})
		return
	}

	c.JSON(http.StatusOK, transaction)
}

// UpdateTransaction updates a transaction
func (h *TransactionHandler) UpdateTransaction(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	transactionID := c.Param("id")
	if transactionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "transaction ID is required"})
		return
	}

	var req domain.UpdateTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	transaction, err := h.transactionService.UpdateTransaction(
		c.Request.Context(),
		transactionID,
		req.Amount,
		req.Description,
		req.Date,
		req.CategoryID,
		req.PaymentMethodID,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, transaction)
}

// DeleteTransaction deletes a transaction
func (h *TransactionHandler) DeleteTransaction(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	transactionID := c.Param("id")
	if transactionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "transaction ID is required"})
		return
	}

	err := h.transactionService.DeleteTransaction(c.Request.Context(), transactionID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "transaction deleted successfully"})
}

// GetTransactionsByCategory returns all transactions for a specific category
func (h *TransactionHandler) GetTransactionsByCategory(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	categoryID := c.Param("categoryId")
	if categoryID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "category ID is required"})
		return
	}

	transactions, err := h.transactionService.GetTransactionsByCategoryID(c.Request.Context(), categoryID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, transactions)
}

// GetTransactionsByDateRange returns all transactions within a date range
func (h *TransactionHandler) GetTransactionsByDateRange(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req domain.DateRangeRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid start date format, use YYYY-MM-DD"})
		return
	}

	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid end date format, use YYYY-MM-DD"})
		return
	}

	// Set endDate to the end of the day
	endDate = endDate.Add(23*time.Hour + 59*time.Minute + 59*time.Second)

	transactions, err := h.transactionService.GetTransactionsByDateRange(c.Request.Context(), userID, startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error retrieving transactions"})
		return
	}

	c.JSON(http.StatusOK, transactions)
}
