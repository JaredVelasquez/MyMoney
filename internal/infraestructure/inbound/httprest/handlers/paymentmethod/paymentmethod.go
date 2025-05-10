package handler

import (
	"net/http"

	payment_method "mi-app-backend/internal/application/paymentmethod"

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

// Create maneja la creación de un nuevo método de pago
func (h *PaymentMethodHandler) Create(c *gin.Context) {
	var req struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Obtener el ID de usuario del contexto (establecido por middleware de autenticación)
	userID, exists := c.Get("userID")
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

// GetByID obtiene un método de pago por su ID
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
	userID, exists := c.Get("userID")
	if !exists || paymentMethod.UserID != userID.(string) {
		c.JSON(http.StatusForbidden, gin.H{"error": "No tienes permiso para acceder a este método de pago"})
		return
	}

	c.JSON(http.StatusOK, paymentMethod)
}

// GetAll obtiene todos los métodos de pago del usuario autenticado
func (h *PaymentMethodHandler) GetAll(c *gin.Context) {
	userID, exists := c.Get("userID")
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

// Update actualiza un método de pago existente
func (h *PaymentMethodHandler) Update(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID no proporcionado"})
		return
	}

	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		IsActive    bool   `json:"is_active"`
	}

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

	userID, exists := c.Get("userID")
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

// Delete elimina un método de pago
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

	userID, exists := c.Get("userID")
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
