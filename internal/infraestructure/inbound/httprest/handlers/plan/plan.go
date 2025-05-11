package plan

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"MyMoneyBackend/internal/application/plan"
	"MyMoneyBackend/internal/domain"
)

// Handler maneja las solicitudes HTTP relacionadas con los planes
type Handler struct {
	service *plan.Service
}

// NewPlanHandler crea una nueva instancia de Handler
func NewPlanHandler(service *plan.Service) *Handler {
	return &Handler{
		service: service,
	}
}

// mapFeaturesToDomain convierte las características de la solicitud al dominio
func mapFeaturesToDomain(features []domain.PlanFeatureRequest) []domain.PlanFeature {
	result := make([]domain.PlanFeature, len(features))
	for i, f := range features {
		result[i] = domain.PlanFeature{
			Name:        f.Name,
			Description: f.Description,
			Value:       f.Value,
			Included:    f.Included,
		}
	}
	return result
}

// mapFeaturesToResponse convierte las características del dominio a la respuesta
func mapFeaturesToResponse(features []domain.PlanFeature) []domain.PlanFeatureResponse {
	result := make([]domain.PlanFeatureResponse, len(features))
	for i, f := range features {
		result[i] = domain.PlanFeatureResponse{
			Name:        f.Name,
			Description: f.Description,
			Value:       f.Value,
			Included:    f.Included,
		}
	}
	return result
}

// mapPlanToPlanResponse convierte un objeto Plan a PlanResponse
func mapPlanToPlanResponse(plan *domain.Plan) domain.PlanResponse {
	return domain.PlanResponse{
		ID:          plan.ID,
		Name:        plan.Name,
		Description: plan.Description,
		Price:       plan.Price,
		CurrencyID:  plan.CurrencyID,
		Interval:    string(plan.Interval),
		Features:    mapFeaturesToResponse(plan.Features),
		IsActive:    plan.IsActive,
		IsPublic:    plan.IsPublic,
		SortOrder:   plan.SortOrder,
		CreatedAt:   plan.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   plan.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

// GetAllPlans godoc
// @Summary Obtener todos los planes
// @Description Retorna una lista de todos los planes disponibles en el sistema
// @Tags plans
// @Accept json
// @Produce json
// @Success 200 {array} domain.PlanResponse
// @Router /plans [get]
func (h *Handler) GetAllPlans(c *gin.Context) {
	ctx := c.Request.Context()

	plans, err := h.service.GetAllPlans(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener planes: " + err.Error()})
		return
	}

	var response []domain.PlanResponse
	for _, p := range plans {
		response = append(response, mapPlanToPlanResponse(p))
	}

	c.JSON(http.StatusOK, response)
}

// GetPublicPlans godoc
// @Summary Obtener planes públicos
// @Description Retorna una lista de los planes públicos disponibles en el sistema
// @Tags plans
// @Accept json
// @Produce json
// @Success 200 {array} domain.PlanResponse
// @Router /plans/public [get]
func (h *Handler) GetPublicPlans(c *gin.Context) {
	ctx := c.Request.Context()

	plans, err := h.service.GetPublicPlans(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener planes públicos: " + err.Error()})
		return
	}

	var response []domain.PlanResponse
	for _, p := range plans {
		response = append(response, mapPlanToPlanResponse(p))
	}

	c.JSON(http.StatusOK, response)
}

// GetActivePlans godoc
// @Summary Obtener planes activos
// @Description Retorna una lista de los planes activos en el sistema
// @Tags plans
// @Accept json
// @Produce json
// @Success 200 {array} domain.PlanResponse
// @Router /plans/active [get]
func (h *Handler) GetActivePlans(c *gin.Context) {
	ctx := c.Request.Context()

	plans, err := h.service.GetActivePlans(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener planes activos: " + err.Error()})
		return
	}

	var response []domain.PlanResponse
	for _, p := range plans {
		response = append(response, mapPlanToPlanResponse(p))
	}

	c.JSON(http.StatusOK, response)
}

// GetPlanByID godoc
// @Summary Obtener un plan por ID
// @Description Retorna los detalles de un plan basado en su ID
// @Tags plans
// @Accept json
// @Produce json
// @Param id path string true "ID del plan"
// @Success 200 {object} domain.PlanResponse
// @Failure 404 {object} map[string]string
// @Router /plans/{id} [get]
func (h *Handler) GetPlanByID(c *gin.Context) {
	id := c.Param("id")
	ctx := c.Request.Context()

	plan, err := h.service.GetPlanByID(ctx, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Plan no encontrado: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, mapPlanToPlanResponse(plan))
}

// CreatePlan godoc
// @Summary Crear un nuevo plan
// @Description Crea un nuevo plan con los datos proporcionados
// @Tags plans
// @Accept json
// @Produce json
// @Param plan body domain.CreatePlanRequest true "Datos del plan"
// @Security Bearer
// @Success 201 {object} domain.PlanResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 409 {object} map[string]string
// @Router /plans [post]
func (h *Handler) CreatePlan(c *gin.Context) {
	var request domain.CreatePlanRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos: " + err.Error()})
		return
	}

	// Convertir el intervalo
	var interval domain.PlanInterval
	switch request.Interval {
	case "monthly":
		interval = domain.PlanIntervalMonthly
	case "yearly":
		interval = domain.PlanIntervalYearly
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Intervalo de facturación inválido. Debe ser 'monthly' o 'yearly'"})
		return
	}

	// Convertir las características
	features := mapFeaturesToDomain(request.Features)

	ctx := c.Request.Context()
	plan, err := h.service.CreatePlan(
		ctx,
		request.Name,
		request.Description,
		request.Price,
		request.CurrencyID,
		interval,
		features,
		request.IsActive,
		request.IsPublic,
		request.SortOrder,
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error al crear plan: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, mapPlanToPlanResponse(plan))
}

// UpdatePlan godoc
// @Summary Actualizar un plan existente
// @Description Actualiza los datos de un plan existente basado en su ID
// @Tags plans
// @Accept json
// @Produce json
// @Param id path string true "ID del plan"
// @Param plan body domain.UpdatePlanRequest true "Datos actualizados del plan"
// @Security Bearer
// @Success 200 {object} domain.PlanResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 404 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Router /plans/{id} [put]
func (h *Handler) UpdatePlan(c *gin.Context) {
	id := c.Param("id")
	var request domain.UpdatePlanRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos: " + err.Error()})
		return
	}

	// Convertir el intervalo
	var interval domain.PlanInterval
	switch request.Interval {
	case "monthly":
		interval = domain.PlanIntervalMonthly
	case "yearly":
		interval = domain.PlanIntervalYearly
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Intervalo de facturación inválido. Debe ser 'monthly' o 'yearly'"})
		return
	}

	// Convertir las características
	features := mapFeaturesToDomain(request.Features)

	ctx := c.Request.Context()
	plan, err := h.service.UpdatePlan(
		ctx,
		id,
		request.Name,
		request.Description,
		request.Price,
		request.CurrencyID,
		interval,
		features,
		request.IsActive,
		request.IsPublic,
		request.SortOrder,
	)

	if err != nil {
		statusCode := http.StatusBadRequest
		if err.Error() == "plan no encontrado" || err.Error() == "plan no encontrado con id: "+id {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, gin.H{"error": "Error al actualizar plan: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, mapPlanToPlanResponse(plan))
}

// DeletePlan godoc
// @Summary Eliminar un plan
// @Description Elimina un plan basado en su ID
// @Tags plans
// @Accept json
// @Produce json
// @Param id path string true "ID del plan"
// @Security Bearer
// @Success 204 "No Content"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 404 {object} map[string]string
// @Router /plans/{id} [delete]
func (h *Handler) DeletePlan(c *gin.Context) {
	id := c.Param("id")
	ctx := c.Request.Context()

	if err := h.service.DeletePlan(ctx, id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Error al eliminar plan: " + err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
