package user_subscription

import (
	"github.com/gin-gonic/gin"

	"mi-app-backend/internal/infraestructure/inbound/httprest/handlers/user_subscription"
)

// SetupUserSubscriptionRoutes configura las rutas para las suscripciones de usuarios
func SetupUserSubscriptionRoutes(
	router *gin.Engine,
	authMiddleware gin.HandlerFunc,
	adminMiddleware gin.HandlerFunc,
	handler *user_subscription.Handler,
) {
	// Grupo de rutas para suscripciones de usuarios
	subscriptionRoutes := router.Group("/subscriptions")
	{
		// Rutas que requieren autenticación
		authRoutes := subscriptionRoutes.Group("/")
		authRoutes.Use(authMiddleware)
		{
			// Obtener suscripción activa del usuario
			authRoutes.GET("/active", handler.GetActiveSubscription)

			// Listar todas las suscripciones del usuario
			authRoutes.GET("/", handler.GetUserSubscriptions)

			// Crear una nueva suscripción
			authRoutes.POST("/", handler.CreateSubscription)

			// Obtener una suscripción específica
			authRoutes.GET("/:id", handler.GetSubscriptionByID)

			// Cancelar una suscripción
			authRoutes.PUT("/:id/cancel", handler.CancelSubscription)

			// Cambiar el plan de una suscripción
			authRoutes.PUT("/:id/plan", handler.ChangePlan)

			// Actualizar el método de pago
			authRoutes.PUT("/:id/payment-method", handler.UpdatePaymentMethod)

			// Renovar una suscripción
			authRoutes.PUT("/:id/renew", handler.RenewSubscription)
		}

		// Rutas administrativas
		adminRoutes := subscriptionRoutes.Group("/admin")
		adminRoutes.Use(authMiddleware, adminMiddleware)
		{
			// Listar suscripciones por estado
			adminRoutes.GET("/status", handler.ListSubscriptionsByStatus)

			// Obtener suscripciones que expirarán pronto
			adminRoutes.GET("/expiring", handler.GetExpiringSubscriptions)

			// Obtener suscripciones pendientes de renovación
			adminRoutes.GET("/pending-renewals", handler.GetPendingRenewals)
		}
	}
}
