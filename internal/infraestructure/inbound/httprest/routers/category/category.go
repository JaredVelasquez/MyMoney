package category

import (
	handler "mi-app-backend/internal/infraestructure/inbound/httprest/handlers/category"
	middleware "mi-app-backend/internal/infraestructure/inbound/httprest/middlewares"

	"github.com/gin-gonic/gin"
)

// SetupCategoryRoutes configura las rutas para las categorías
func SetupCategoryRoutes(router *gin.RouterGroup, categoryHandler *handler.CategoryHandler, authMiddleware *middleware.AuthMiddleware) {
	// Todas las rutas de categorías requieren autenticación
	categories := router.Group("/categories")
	categories.Use(authMiddleware.Authorize())
	{
		categories.POST("", categoryHandler.CreateCategory)
		categories.GET("", categoryHandler.GetUserCategories)
		categories.GET("/:id", categoryHandler.GetCategory)
		categories.PUT("/:id", categoryHandler.UpdateCategory)
		categories.DELETE("/:id", categoryHandler.DeleteCategory)
		categories.GET("/type/:type", categoryHandler.GetCategoriesByType)
	}
}
