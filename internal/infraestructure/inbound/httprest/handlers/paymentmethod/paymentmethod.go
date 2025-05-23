package handler

import (
	"net/http"

	payment_method "MyMoneyBackend/internal/application/paymentmethod"
	"MyMoneyBackend/internal/domain"
	middleware "MyMoneyBackend/internal/infraestructure/inbound/httprest/middlewares"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// PaymentMethodHandler maneja las peticiones HTTP relacionadas con métodos de pago
type PaymentMethodHandler struct {
	service  *payment_method.Service
	validate *validator.Validate
}

// NewPaymentMethodHandler crea un nuevo handler para métodos de pago
func NewPaymentMethodHandler(service *payment_method.Service) *PaymentMethodHandler {
	return &PaymentMethodHandler{
		service:  service,
		validate: validator.New(),
	}
}

// Create godoc
// @Summary Crear un método de pago
// @Description Crea un nuevo método de pago para el usuario autenticado
// @Tags payment-methods
// @Accept json
// @Produce json
// @Param method body domain.CreatePaymentMethodRequest true "Datos del método de pago"
// @Security Bearer
// @Success 201 {object} domain.PaymentMethod
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /payment-methods [post]
func (h *PaymentMethodHandler) Create(c *gin.Context) {
	var req domain.CreatePaymentMethodRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Obtener el ID de usuario del contexto (establecido por middleware de autenticación)
	userID, exists := c.Get(middleware.UserIDKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autenticado"})
		return
	}

	paymentMethod, err := h.service.CreatePaymentMethod(c.Request.Context(), req.Name, req.Description, userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, paymentMethod)
}

// GetByID godoc
// @Summary Obtener un método de pago por ID
// @Description Retorna los detalles de un método de pago específico
// @Tags payment-methods
// @Accept json
// @Produce json
// @Param id path string true "ID del método de pago"
// @Security Bearer
// @Success 200 {object} domain.PaymentMethod
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /payment-methods/{id} [get]
func (h *PaymentMethodHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID no proporcionado"})
		return
	}

	paymentMethod, err := h.service.GetPaymentMethodByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if paymentMethod == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Método de pago no encontrado"})
		return
	}

	// Verificar que el método de pago pertenece al usuario autenticado
	userID, exists := c.Get(middleware.UserIDKey)
	if !exists || paymentMethod.UserID != userID.(string) {
		c.JSON(http.StatusForbidden, gin.H{"error": "No tienes permiso para acceder a este método de pago"})
		return
	}

	c.JSON(http.StatusOK, paymentMethod)
}

// GetAll godoc
// @Summary Obtener todos los métodos de pago
// @Description Retorna todos los métodos de pago del usuario autenticado
// @Tags payment-methods
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {array} domain.PaymentMethod
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /payment-methods [get]
func (h *PaymentMethodHandler) GetAll(c *gin.Context) {
	userID, exists := c.Get(middleware.UserIDKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autenticado"})
		return
	}

	paymentMethods, err := h.service.GetPaymentMethodsByUserID(c.Request.Context(), userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, paymentMethods)
}

// Update godoc
// @Summary Actualizar un método de pago
// @Description Actualiza los datos de un método de pago existente
// @Tags payment-methods
// @Accept json
// @Produce json
// @Param id path string true "ID del método de pago"
// @Param method body domain.UpdatePaymentMethodRequest true "Datos actualizados del método de pago"
// @Security Bearer
// @Success 200 {object} domain.PaymentMethod
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /payment-methods/{id} [put]
func (h *PaymentMethodHandler) Update(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID no proporcionado"})
		return
	}

	var req domain.UpdatePaymentMethodRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verificar que el método de pago existe y pertenece al usuario
	existingPaymentMethod, err := h.service.GetPaymentMethodByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if existingPaymentMethod == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Método de pago no encontrado"})
		return
	}

	userID, exists := c.Get(middleware.UserIDKey)
	if !exists || existingPaymentMethod.UserID != userID.(string) {
		c.JSON(http.StatusForbidden, gin.H{"error": "No tienes permiso para modificar este método de pago"})
		return
	}

	updatedPaymentMethod, err := h.service.UpdatePaymentMethod(
		c.Request.Context(),
		id,
		req.Name,
		req.Description,
		req.IsActive,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedPaymentMethod)
}

// Delete godoc
// @Summary Eliminar un método de pago
// @Description Elimina un método de pago existente
// @Tags payment-methods
// @Accept json
// @Produce json
// @Param id path string true "ID del método de pago"
// @Security Bearer
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /payment-methods/{id} [delete]
func (h *PaymentMethodHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID no proporcionado"})
		return
	}

	// Verificar que el método de pago existe y pertenece al usuario
	existingPaymentMethod, err := h.service.GetPaymentMethodByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if existingPaymentMethod == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Método de pago no encontrado"})
		return
	}

	userID, exists := c.Get(middleware.UserIDKey)
	if !exists || existingPaymentMethod.UserID != userID.(string) {
		c.JSON(http.StatusForbidden, gin.H{"error": "No tienes permiso para eliminar este método de pago"})
		return
	}

	if err := h.service.DeletePaymentMethod(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Método de pago eliminado correctamente"})
}
