package plan

import (
	"github.com/gin-gonic/gin"

	handler "mi-app-backend/internal/infraestructure/inbound/httprest/handlers/plan"
	middleware "mi-app-backend/internal/infraestructure/inbound/httprest/middlewares"
)

// SetupPlanRoutes configura las rutas relacionadas con los planes
func SetupPlanRoutes(r *gin.RouterGroup, planHandler *handler.Handler, authMiddleware *middleware.AuthMiddleware) {
	planRoutes := r.Group("/plans")
	{
		// Rutas públicas (solo lectura)
		planRoutes.GET("", planHandler.GetAllPlans)
		planRoutes.GET("/active", planHandler.GetActivePlans)
		planRoutes.GET("/public", planHandler.GetPublicPlans)
		planRoutes.GET("/:id", planHandler.GetPlanByID)

		// Rutas protegidas (requieren autenticación)
		protected := planRoutes.Group("/")
		protected.Use(authMiddleware.Authorize())
		{
			protected.POST("", planHandler.CreatePlan)
			protected.PUT("/:id", planHandler.UpdatePlan)
			protected.DELETE("/:id", planHandler.DeletePlan)
		}
	}
}
