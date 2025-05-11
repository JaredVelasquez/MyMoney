package currency

import (
	"github.com/gin-gonic/gin"

	handler "MyMoneyBackend/internal/infraestructure/inbound/httprest/handlers/currency"
	middleware "MyMoneyBackend/internal/infraestructure/inbound/httprest/middlewares"
)

// SetupCurrencyRoutes configura las rutas relacionadas con las monedas
func SetupCurrencyRoutes(r *gin.RouterGroup, currencyHandler *handler.Handler, authMiddleware *middleware.AuthMiddleware) {
	currencyRoutes := r.Group("/currencies")
	{
		// Rutas públicas (solo lectura)
		currencyRoutes.GET("", currencyHandler.GetAllCurrencies)
		currencyRoutes.GET("/active", currencyHandler.GetActiveCurrencies)
		currencyRoutes.GET("/:id", currencyHandler.GetCurrencyByID)
		currencyRoutes.GET("/code/:code", currencyHandler.GetCurrencyByCode)

		// Rutas protegidas (requieren autenticación)
		protected := currencyRoutes.Group("/")
		protected.Use(authMiddleware.Authorize())
		{
			protected.POST("", currencyHandler.CreateCurrency)
			protected.PUT("/:id", currencyHandler.UpdateCurrency)
			protected.DELETE("/:id", currencyHandler.DeleteCurrency)
		}
	}
}
