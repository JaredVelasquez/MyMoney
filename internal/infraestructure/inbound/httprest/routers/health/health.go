package health

import (
	"github.com/gin-gonic/gin"

	"mi-app-backend/internal/infraestructure/inbound/httprest/handlers/health"
)

// SetupHealthRoutes configura las rutas de health check
func SetupHealthRoutes(r *gin.RouterGroup, healthHandler *health.Handler) {
	// Rutas públicas (no requieren autenticación)
	healthRoutes := r.Group("/health")
	{
		healthRoutes.GET("", healthHandler.Status)
		healthRoutes.GET("/check", healthHandler.Check)
	}
}
