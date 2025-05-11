package api

import (
	"net/http"
	"time"

	transaction "MyMoneyBackend/internal/application/transaction"
	"MyMoneyBackend/internal/domain"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Constantes
const (
	// DefaultCurrencyUUID es el UUID por defecto para la moneda USD
	DefaultCurrencyUUID = "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11"
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
// @Summary Crear una nueva transacción
// @Description Crea una nueva transacción para el usuario autenticado
// @Tags transactions
// @Accept json
// @Produce json
// @Security Bearer
// @Param transaction body domain.CreateTransactionRequest true "Datos de la transacción a crear"
// @Success 201 {object} domain.Transaction
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /api/transactions [post]
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

	// Asegurarse de que currency_id sea un UUID válido
	currencyID := req.CurrencyID
	// Si el cliente está enviando un código de moneda o un UUID incorrecto, usar el UUID de USD
	if _, err := uuid.Parse(currencyID); err != nil {
		// UUID para USD obtenido de la migración
		currencyID = DefaultCurrencyUUID
	}

	transaction, err := h.transactionService.CreateTransaction(
		c.Request.Context(),
		req.Amount,
		req.Description,
		req.Date,
		req.CategoryID,
		req.PaymentMethodID,
		userID,
		currencyID,
		req.Type,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, transaction)
}

// GetUserTransactions returns all transactions for the current user
// @Summary Obtener todas las transacciones del usuario
// @Description Retorna todas las transacciones del usuario autenticado
// @Tags transactions
// @Produce json
// @Security Bearer
// @Success 200 {array} domain.Transaction
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/transactions [get]
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
// @Summary Obtener una transacción específica
// @Description Retorna una transacción específica por su ID
// @Tags transactions
// @Produce json
// @Security Bearer
// @Param id path string true "ID de la transacción"
// @Success 200 {object} domain.Transaction
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/transactions/{id} [get]
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
// @Summary Actualizar una transacción
// @Description Actualiza una transacción existente
// @Tags transactions
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "ID de la transacción"
// @Param transaction body domain.UpdateTransactionRequest true "Datos de la transacción a actualizar"
// @Success 200 {object} domain.Transaction
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/transactions/{id} [put]
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

	// Asegurarse de que currency_id sea un UUID válido si está presente
	currencyID := req.CurrencyID
	if currencyID != "" {
		if _, err := uuid.Parse(currencyID); err != nil {
			// UUID para USD obtenido de la migración
			currencyID = DefaultCurrencyUUID
		}
	}

	transaction, err := h.transactionService.UpdateTransaction(
		c.Request.Context(),
		transactionID,
		req.Amount,
		req.Description,
		req.Date,
		req.CategoryID,
		req.PaymentMethodID,
		currencyID,
		req.Type,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, transaction)
}

// DeleteTransaction deletes a transaction
// @Summary Eliminar una transacción
// @Description Elimina una transacción existente
// @Tags transactions
// @Security Bearer
// @Param id path string true "ID de la transacción"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /api/transactions/{id} [delete]
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
// @Summary Obtener transacciones por categoría
// @Description Retorna todas las transacciones de una categoría específica
// @Tags transactions
// @Produce json
// @Security Bearer
// @Param categoryId path string true "ID de la categoría"
// @Success 200 {array} domain.Transaction
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /api/transactions/category/{categoryId} [get]
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
// @Summary Obtener transacciones por rango de fechas
// @Description Retorna todas las transacciones dentro de un rango de fechas
// @Tags transactions
// @Produce json
// @Security Bearer
// @Param start_date query string true "Fecha de inicio (formato YYYY-MM-DD)"
// @Param end_date query string true "Fecha de fin (formato YYYY-MM-DD)"
// @Success 200 {array} domain.Transaction
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/transactions/date-range [get]
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
