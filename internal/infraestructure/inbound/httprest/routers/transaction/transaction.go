package transaction

import (
	handler "mi-app-backend/internal/infraestructure/inbound/httprest/handlers/transaction"
	middleware "mi-app-backend/internal/infraestructure/inbound/httprest/middlewares"

	"github.com/gin-gonic/gin"
)

// SetupTransactionRoutes configura las rutas para las transacciones
func SetupTransactionRoutes(router *gin.RouterGroup, transactionHandler *handler.TransactionHandler, authMiddleware *middleware.AuthMiddleware) {
	// Todas las rutas de transacciones requieren autenticación
	transactions := router.Group("/transactions")
	transactions.Use(authMiddleware.Authorize())
	{
		transactions.POST("", transactionHandler.CreateTransaction)
		transactions.GET("", transactionHandler.GetUserTransactions)
		transactions.GET("/:id", transactionHandler.GetTransaction)
		transactions.PUT("/:id", transactionHandler.UpdateTransaction)
		transactions.DELETE("/:id", transactionHandler.DeleteTransaction)
		transactions.GET("/category/:categoryId", transactionHandler.GetTransactionsByCategory)
		transactions.GET("/date-range", transactionHandler.GetTransactionsByDateRange)
	}
}
