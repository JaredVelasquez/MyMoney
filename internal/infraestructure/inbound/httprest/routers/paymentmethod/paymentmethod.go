package paymentmethod

import (
	handler "MyMoneyBackend/internal/infraestructure/inbound/httprest/handlers/paymentmethod"
	middleware "MyMoneyBackend/internal/infraestructure/inbound/httprest/middlewares"

	"github.com/gin-gonic/gin"
)

// SetupPaymentMethodRoutes configura las rutas para los métodos de pago
func SetupPaymentMethodRoutes(router *gin.RouterGroup, paymentMethodHandler *handler.PaymentMethodHandler, authMiddleware *middleware.AuthMiddleware) {
	// Todas las rutas de métodos de pago requieren autenticación
	paymentMethods := router.Group("/payment-methods")
	paymentMethods.Use(authMiddleware.Authorize())
	{
		paymentMethods.POST("", paymentMethodHandler.Create)
		paymentMethods.GET("", paymentMethodHandler.GetAll)
		paymentMethods.GET("/:id", paymentMethodHandler.GetByID)
		paymentMethods.PUT("/:id", paymentMethodHandler.Update)
		paymentMethods.DELETE("/:id", paymentMethodHandler.Delete)
	}
}
