package user_subscription

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"mi-app-backend/internal/application/user_subscription"
	"mi-app-backend/internal/domain"
	middleware "mi-app-backend/internal/infraestructure/inbound/httprest/middlewares"
)

// Handler maneja las solicitudes HTTP relacionadas con suscripciones de usuarios
type Handler struct {
	service *user_subscription.Service
}

// NewUserSubscriptionHandler crea una nueva instancia del controlador de suscripciones
func NewUserSubscriptionHandler(service *user_subscription.Service) *Handler {
	return &Handler{
		service: service,
	}
}

// CreateSubscription godoc
// @Summary Crear una suscripción
// @Description Crea una nueva suscripción para el usuario autenticado
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param subscription body CreateSubscriptionRequest true "Datos de la suscripción"
// @Security Bearer
// @Success 201 {object} SubscriptionResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Router /subscriptions [post]
func (h *Handler) CreateSubscription(c *gin.Context) {
	// Obtener el usuario del contexto
	userID, exists := c.Get(middleware.UserIDKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autenticado"})
		return
	}

	var req domain.CreateSubscriptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos de solicitud inválidos: " + err.Error()})
		return
	}

	// Si no se proporcionan fechas, establecer valores predeterminados
	now := time.Now()
	if req.StartDate.IsZero() {
		req.StartDate = now
	}

	if req.EndDate.IsZero() {
		// Por defecto, la suscripción dura un mes
		req.EndDate = now.AddDate(0, 1, 0)
	}

	// Crear la suscripción
	subscription, err := h.service.CreateSubscription(
		c.Request.Context(),
		userID.(string),
		req.PlanID,
		req.StartDate,
		req.EndDate,
		req.PaymentMethodID,
		req.Metadata,
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error al crear suscripción: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, domain.MapSubscriptionToResponse(subscription))
}

// GetActiveSubscription obtiene la suscripción activa del usuario
func (h *Handler) GetActiveSubscription(c *gin.Context) {
	// Obtener el usuario del contexto
	userID, exists := c.Get(middleware.UserIDKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autenticado"})
		return
	}

	// Obtener la suscripción activa
	subscription, err := h.service.GetActiveSubscription(c.Request.Context(), userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener suscripción: " + err.Error()})
		return
	}

	if subscription == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No hay suscripción activa"})
		return
	}

	c.JSON(http.StatusOK, domain.MapSubscriptionToResponse(subscription))
}

// GetSubscriptions godoc
// @Summary Obtener suscripciones del usuario
// @Description Retorna las suscripciones del usuario autenticado
// @Tags subscriptions
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {array} SubscriptionResponse
// @Failure 401 {object} map[string]string
// @Router /subscriptions [get]
func (h *Handler) GetUserSubscriptions(c *gin.Context) {
	// Obtener el usuario del contexto
	userID, exists := c.Get(middleware.UserIDKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autenticado"})
		return
	}

	// Obtener las suscripciones
	subscriptions, err := h.service.GetUserSubscriptions(c.Request.Context(), userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener suscripciones: " + err.Error()})
		return
	}

	// Mapear las suscripciones a respuestas
	var responses []domain.SubscriptionResponse
	for _, sub := range subscriptions {
		responses = append(responses, domain.MapSubscriptionToResponse(sub))
	}

	c.JSON(http.StatusOK, responses)
}

// GetSubscription godoc
// @Summary Obtener una suscripción
// @Description Retorna los detalles de una suscripción específica
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param id path string true "ID de la suscripción"
// @Security Bearer
// @Success 200 {object} SubscriptionResponse
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /subscriptions/{id} [get]
func (h *Handler) GetSubscriptionByID(c *gin.Context) {
	// Obtener el usuario del contexto
	userID, exists := c.Get(middleware.UserIDKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autenticado"})
		return
	}

	// Obtener el ID de la suscripción
	subscriptionID := c.Param("id")
	if subscriptionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de suscripción no proporcionado"})
		return
	}

	// Obtener la suscripción
	subscription, err := h.service.GetSubscriptionByID(c.Request.Context(), subscriptionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener suscripción: " + err.Error()})
		return
	}

	// Verificar que la suscripción pertenezca al usuario
	if subscription.UserID != userID.(string) {
		c.JSON(http.StatusForbidden, gin.H{"error": "No tiene permiso para acceder a esta suscripción"})
		return
	}

	c.JSON(http.StatusOK, domain.MapSubscriptionToResponse(subscription))
}

// CancelSubscription godoc
// @Summary Cancelar una suscripción
// @Description Cancela una suscripción activa
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param id path string true "ID de la suscripción"
// @Security Bearer
// @Success 200 {object} SubscriptionResponse
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Router /subscriptions/{id}/cancel [put]
func (h *Handler) CancelSubscription(c *gin.Context) {
	// Obtener el usuario del contexto
	userID, exists := c.Get(middleware.UserIDKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autenticado"})
		return
	}

	// Obtener el ID de la suscripción
	subscriptionID := c.Param("id")
	if subscriptionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de suscripción no proporcionado"})
		return
	}

	// Obtener la suscripción para verificar pertenencia
	subscription, err := h.service.GetSubscriptionByID(c.Request.Context(), subscriptionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener suscripción: " + err.Error()})
		return
	}

	// Verificar que la suscripción pertenezca al usuario
	if subscription.UserID != userID.(string) {
		c.JSON(http.StatusForbidden, gin.H{"error": "No tiene permiso para cancelar esta suscripción"})
		return
	}

	// Parsear la solicitud de cancelación
	var req domain.CancelSubscriptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// Si no hay cuerpo, usar una razón predeterminada
		req.Reason = "Cancelada por solicitud del usuario"
	}

	// Cancelar la suscripción
	err = h.service.CancelSubscription(c.Request.Context(), subscriptionID, req.Reason)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error al cancelar suscripción: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Suscripción cancelada correctamente"})
}

// ChangePlan cambia el plan de una suscripción
func (h *Handler) ChangePlan(c *gin.Context) {
	// Obtener el usuario del contexto
	userID, exists := c.Get(middleware.UserIDKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autenticado"})
		return
	}

	// Obtener el ID de la suscripción
	subscriptionID := c.Param("id")
	if subscriptionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de suscripción no proporcionado"})
		return
	}

	// Obtener la suscripción para verificar pertenencia
	subscription, err := h.service.GetSubscriptionByID(c.Request.Context(), subscriptionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener suscripción: " + err.Error()})
		return
	}

	// Verificar que la suscripción pertenezca al usuario
	if subscription.UserID != userID.(string) {
		c.JSON(http.StatusForbidden, gin.H{"error": "No tiene permiso para modificar esta suscripción"})
		return
	}

	// Parsear la solicitud de cambio de plan
	var req domain.ChangePlanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos de solicitud inválidos: " + err.Error()})
		return
	}

	// Cambiar el plan
	updatedSubscription, err := h.service.ChangeSubscriptionPlan(c.Request.Context(), subscriptionID, req.PlanID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error al cambiar plan: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, domain.MapSubscriptionToResponse(updatedSubscription))
}

// UpdatePaymentMethod actualiza el método de pago de una suscripción
func (h *Handler) UpdatePaymentMethod(c *gin.Context) {
	// Obtener el usuario del contexto
	userID, exists := c.Get(middleware.UserIDKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autenticado"})
		return
	}

	// Obtener el ID de la suscripción
	subscriptionID := c.Param("id")
	if subscriptionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de suscripción no proporcionado"})
		return
	}

	// Obtener la suscripción para verificar pertenencia
	subscription, err := h.service.GetSubscriptionByID(c.Request.Context(), subscriptionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener suscripción: " + err.Error()})
		return
	}

	// Verificar que la suscripción pertenezca al usuario
	if subscription.UserID != userID.(string) {
		c.JSON(http.StatusForbidden, gin.H{"error": "No tiene permiso para modificar esta suscripción"})
		return
	}

	// Parsear la solicitud
	var req domain.UpdatePaymentMethodRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos de solicitud inválidos: " + err.Error()})
		return
	}

	// Actualizar el método de pago
	updatedSubscription, err := h.service.UpdatePaymentMethod(c.Request.Context(), subscriptionID, req.PaymentMethodID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error al actualizar método de pago: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, domain.MapSubscriptionToResponse(updatedSubscription))
}

// GetAllSubscriptions godoc
// @Summary Obtener todas las suscripciones
// @Description Retorna todas las suscripciones (solo para administradores)
// @Tags subscriptions
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {array} SubscriptionResponse
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Router /admin/subscriptions [get]
func (h *Handler) ListSubscriptionsByStatus(c *gin.Context) {
	// Verificar si es administrador
	isAdmin, exists := c.Get(middleware.IsAdminKey)
	if !exists || !isAdmin.(bool) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Acceso denegado. Se requieren permisos de administrador"})
		return
	}

	// Obtener el estado de la consulta
	status := c.Query("status")
	if status == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parámetro de estado requerido"})
		return
	}

	// Obtener las suscripciones
	subscriptions, err := h.service.GetSubscriptionsByStatus(c.Request.Context(), domain.SubscriptionStatus(status))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener suscripciones: " + err.Error()})
		return
	}

	// Mapear las suscripciones a respuestas
	var responses []domain.SubscriptionResponse
	for _, sub := range subscriptions {
		responses = append(responses, domain.MapSubscriptionToResponse(sub))
	}

	c.JSON(http.StatusOK, responses)
}

// GetExpiringSubscriptions obtiene suscripciones que expirarán pronto (solo para administradores)
func (h *Handler) GetExpiringSubscriptions(c *gin.Context) {
	// Verificar si es administrador
	isAdmin, exists := c.Get(middleware.IsAdminKey)
	if !exists || !isAdmin.(bool) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Acceso denegado. Se requieren permisos de administrador"})
		return
	}

	// Parámetro de días (predeterminado: 7 días)
	days := 7
	// TODO: Implementar parseo de parámetro de días si es necesario

	// Obtener las suscripciones que expirarán pronto
	subscriptions, err := h.service.GetExpiringSubscriptions(c.Request.Context(), days)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener suscripciones: " + err.Error()})
		return
	}

	// Mapear las suscripciones a respuestas
	var responses []domain.SubscriptionResponse
	for _, sub := range subscriptions {
		responses = append(responses, domain.MapSubscriptionToResponse(sub))
	}

	c.JSON(http.StatusOK, responses)
}

// GetPendingRenewals obtiene suscripciones pendientes de renovación (solo para administradores)
func (h *Handler) GetPendingRenewals(c *gin.Context) {
	// Verificar si es administrador
	isAdmin, exists := c.Get(middleware.IsAdminKey)
	if !exists || !isAdmin.(bool) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Acceso denegado. Se requieren permisos de administrador"})
		return
	}

	// Parámetro de días (predeterminado: 7 días)
	days := 7
	// TODO: Implementar parseo de parámetro de días si es necesario

	// Obtener las suscripciones pendientes de renovación
	subscriptions, err := h.service.GetPendingRenewals(c.Request.Context(), days)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener renovaciones pendientes: " + err.Error()})
		return
	}

	// Mapear las suscripciones a respuestas
	var responses []domain.SubscriptionResponse
	for _, sub := range subscriptions {
		responses = append(responses, domain.MapSubscriptionToResponse(sub))
	}

	c.JSON(http.StatusOK, responses)
}

// GetSubscriptionsByStatus godoc
// @Summary Obtener suscripciones por estado
// @Description Retorna las suscripciones filtradas por estado (solo para administradores)
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param status path string true "Estado de suscripción (active, expired, cancelled)"
// @Security Bearer
// @Success 200 {array} SubscriptionResponse
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /admin/subscriptions/status/{status} [get]
func (h *Handler) GetSubscriptionsByStatus(c *gin.Context) {
	// ... existing code ...
}

// RenewSubscription godoc
// @Summary Renovar una suscripción
// @Description Renueva una suscripción expirada o cancelada
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param id path string true "ID de la suscripción"
// @Security Bearer
// @Success 200 {object} SubscriptionResponse
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Router /subscriptions/{id}/renew [put]
func (h *Handler) RenewSubscription(c *gin.Context) {
	// ... existing code ...
}
