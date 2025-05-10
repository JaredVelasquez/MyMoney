package user

import (
	handler "mi-app-backend/internal/infraestructure/inbound/httprest/handlers/user"
	middleware "mi-app-backend/internal/infraestructure/inbound/httprest/middlewares"

	"github.com/gin-gonic/gin"
)

// SetupUserRoutes configura las rutas para los usuarios
func SetupUserRoutes(router *gin.RouterGroup, userHandler *handler.UserHandler, authMiddleware *middleware.AuthMiddleware) {
	// Rutas públicas (no requieren autenticación)
	router.POST("/auth/register", userHandler.Register)
	router.POST("/auth/login", userHandler.Login)

	// Rutas protegidas (requieren autenticación)
	users := router.Group("/users")
	users.Use(authMiddleware.Authorize())
	{
		users.GET("/me", userHandler.GetCurrentUser)
		users.PUT("/me", userHandler.UpdateUser)
		users.POST("/change-password", userHandler.ChangePassword)
	}
}
