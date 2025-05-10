package currency

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"mi-app-backend/internal/application/currency"
	"mi-app-backend/internal/domain"
)

// Handler maneja las solicitudes HTTP relacionadas con las monedas
type Handler struct {
	service *currency.Service
}

// NewCurrencyHandler crea una nueva instancia de Handler
func NewCurrencyHandler(service *currency.Service) *Handler {
	return &Handler{
		service: service,
	}
}

// mapCurrencyToCurrencyResponse convierte un objeto Currency a CurrencyResponse
func mapCurrencyToCurrencyResponse(currency *domain.Currency) domain.CurrencyResponse {
	return domain.CurrencyResponse{
		ID:        currency.ID,
		Code:      currency.Code,
		Name:      currency.Name,
		Symbol:    currency.Symbol,
		IsActive:  currency.IsActive,
		CreatedAt: currency.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: currency.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

// GetAllCurrencies godoc
// @Summary Obtener todas las monedas
// @Description Retorna una lista de todas las monedas disponibles en el sistema
// @Tags currencies
// @Accept json
// @Produce json
// @Success 200 {array} CurrencyResponse
// @Router /currencies [get]
func (h *Handler) GetAllCurrencies(c *gin.Context) {
	ctx := c.Request.Context()

	currencies, err := h.service.GetAllCurrencies(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener monedas: " + err.Error()})
		return
	}

	var response []domain.CurrencyResponse
	for _, curr := range currencies {
		response = append(response, mapCurrencyToCurrencyResponse(curr))
	}

	c.JSON(http.StatusOK, response)
}

// GetActiveCurrencies godoc
// @Summary Obtener monedas activas
// @Description Retorna una lista de las monedas activas en el sistema
// @Tags currencies
// @Accept json
// @Produce json
// @Success 200 {array} CurrencyResponse
// @Router /currencies/active [get]
func (h *Handler) GetActiveCurrencies(c *gin.Context) {
	ctx := c.Request.Context()

	currencies, err := h.service.GetActiveCurrencies(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener monedas activas: " + err.Error()})
		return
	}

	var response []domain.CurrencyResponse
	for _, curr := range currencies {
		response = append(response, mapCurrencyToCurrencyResponse(curr))
	}

	c.JSON(http.StatusOK, response)
}

// GetCurrencyByID godoc
// @Summary Obtener una moneda por ID
// @Description Retorna los detalles de una moneda basado en su ID
// @Tags currencies
// @Accept json
// @Produce json
// @Param id path string true "ID de la moneda"
// @Success 200 {object} CurrencyResponse
// @Failure 404 {object} map[string]string
// @Router /currencies/{id} [get]
func (h *Handler) GetCurrencyByID(c *gin.Context) {
	id := c.Param("id")
	ctx := c.Request.Context()

	currency, err := h.service.GetCurrencyByID(ctx, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Moneda no encontrada: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, mapCurrencyToCurrencyResponse(currency))
}

// GetCurrencyByCode godoc
// @Summary Obtener una moneda por código
// @Description Retorna los detalles de una moneda basado en su código
// @Tags currencies
// @Accept json
// @Produce json
// @Param code path string true "Código de la moneda"
// @Success 200 {object} CurrencyResponse
// @Failure 404 {object} map[string]string
// @Router /currencies/code/{code} [get]
func (h *Handler) GetCurrencyByCode(c *gin.Context) {
	code := c.Param("code")
	ctx := c.Request.Context()

	currency, err := h.service.GetCurrencyByCode(ctx, code)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Moneda no encontrada: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, mapCurrencyToCurrencyResponse(currency))
}

// CreateCurrency godoc
// @Summary Crear una nueva moneda
// @Description Crea una nueva moneda con los datos proporcionados
// @Tags currencies
// @Accept json
// @Produce json
// @Param currency body CreateCurrencyRequest true "Datos de la moneda"
// @Security Bearer
// @Success 201 {object} CurrencyResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 409 {object} map[string]string
// @Router /currencies [post]
func (h *Handler) CreateCurrency(c *gin.Context) {
	var request domain.CreateCurrencyRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos: " + err.Error()})
		return
	}

	ctx := c.Request.Context()
	currency, err := h.service.CreateCurrency(ctx, request.Code, request.Name, request.Symbol, request.IsActive)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Error al crear moneda: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, mapCurrencyToCurrencyResponse(currency))
}

// UpdateCurrency godoc
// @Summary Actualizar una moneda existente
// @Description Actualiza los datos de una moneda existente basado en su ID
// @Tags currencies
// @Accept json
// @Produce json
// @Param id path string true "ID de la moneda"
// @Param currency body UpdateCurrencyRequest true "Datos actualizados de la moneda"
// @Security Bearer
// @Success 200 {object} CurrencyResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 404 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Router /currencies/{id} [put]
func (h *Handler) UpdateCurrency(c *gin.Context) {
	id := c.Param("id")

	var request domain.UpdateCurrencyRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos: " + err.Error()})
		return
	}

	ctx := c.Request.Context()
	currency, err := h.service.UpdateCurrency(ctx, id, request.Code, request.Name, request.Symbol, request.IsActive)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Error al actualizar moneda: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, mapCurrencyToCurrencyResponse(currency))
}

// DeleteCurrency godoc
// @Summary Eliminar una moneda
// @Description Elimina una moneda basado en su ID
// @Tags currencies
// @Accept json
// @Produce json
// @Param id path string true "ID de la moneda"
// @Security Bearer
// @Success 204 "No Content"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 404 {object} map[string]string
// @Router /currencies/{id} [delete]
func (h *Handler) DeleteCurrency(c *gin.Context) {
	id := c.Param("id")
	ctx := c.Request.Context()

	err := h.service.DeleteCurrency(ctx, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Error al eliminar moneda: " + err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
