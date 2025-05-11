package routers

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"MyMoneyBackend/db/config"
	"MyMoneyBackend/internal/application/auth"
	categoryService "MyMoneyBackend/internal/application/category"
	currencyService "MyMoneyBackend/internal/application/currency"
	paymentMethodService "MyMoneyBackend/internal/application/paymentmethod"
	planService "MyMoneyBackend/internal/application/plan"
	transactionService "MyMoneyBackend/internal/application/transaction"
	userService "MyMoneyBackend/internal/application/user"
	userSubscriptionService "MyMoneyBackend/internal/application/user_subscription"
	categoryHandler "MyMoneyBackend/internal/infraestructure/inbound/httprest/handlers/category"
	currencyHandler "MyMoneyBackend/internal/infraestructure/inbound/httprest/handlers/currency"
	healthHandler "MyMoneyBackend/internal/infraestructure/inbound/httprest/handlers/health"
	paymentMethodHandler "MyMoneyBackend/internal/infraestructure/inbound/httprest/handlers/paymentmethod"
	planHandler "MyMoneyBackend/internal/infraestructure/inbound/httprest/handlers/plan"
	transactionHandler "MyMoneyBackend/internal/infraestructure/inbound/httprest/handlers/transaction"
	userHandler "MyMoneyBackend/internal/infraestructure/inbound/httprest/handlers/user"
	userSubscriptionHandler "MyMoneyBackend/internal/infraestructure/inbound/httprest/handlers/user_subscription"
	middlewares "MyMoneyBackend/internal/infraestructure/inbound/httprest/middlewares"
	categoryRouter "MyMoneyBackend/internal/infraestructure/inbound/httprest/routers/category"
	currencyRouter "MyMoneyBackend/internal/infraestructure/inbound/httprest/routers/currency"
	healthRouter "MyMoneyBackend/internal/infraestructure/inbound/httprest/routers/health"
	paymentMethodRouter "MyMoneyBackend/internal/infraestructure/inbound/httprest/routers/paymentmethod"
	planRouter "MyMoneyBackend/internal/infraestructure/inbound/httprest/routers/plan"
	transactionRouter "MyMoneyBackend/internal/infraestructure/inbound/httprest/routers/transaction"
	userRouter "MyMoneyBackend/internal/infraestructure/inbound/httprest/routers/user"
	userSubscriptionRouter "MyMoneyBackend/internal/infraestructure/inbound/httprest/routers/user_subscription"
	"MyMoneyBackend/internal/infraestructure/inbound/httprest/swagger"
	"MyMoneyBackend/internal/infraestructure/outbound/repository"
)

// SetupRouter configura todas las rutas de la API
func SetupRouter(
	r *gin.Engine,
	userSvc *userService.UserService,
	categorySvc *categoryService.Service,
	paymentMethodSvc *paymentMethodService.Service,
	transactionSvc *transactionService.Service,
	tokenSvc *auth.TokenService,
) {
	// Configurar CORS
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Configurar Swagger
	swagger.SetupSwaggerRoutes(r)

	// Crear middleware de autenticación
	authMiddleware := middlewares.NewAuthMiddleware(tokenSvc)
	adminMiddleware := middlewares.NewAdminMiddleware()

	// Crear handlers
	userHdlr := userHandler.NewUserHandler(*userSvc, *tokenSvc)
	categoryHdlr := categoryHandler.NewCategoryHandler(categorySvc)
	paymentMethodHdlr := paymentMethodHandler.NewPaymentMethodHandler(paymentMethodSvc)
	transactionHdlr := transactionHandler.NewTransactionHandler(transactionSvc)
	healthHdlr := healthHandler.NewHealthHandler()

	// Obtener conexión a la base de datos para los servicios adicionales
	var db *sql.DB
	dbConn, err := config.NewConnection()
	if err != nil {
		log.Printf("Error al conectar a la BD: %v", err)
	} else {
		db = dbConn.GetDB()
	}

	// Inicializar repositorios
	currencyRepo := repository.NewCurrencyRepository(db)
	planRepo := repository.NewPlanRepository(db)
	userSubscriptionRepo := repository.NewUserSubscriptionRepository(db)
	userRepo := repository.NewUserRepository(db)

	// Inicializar servicios
	currencySvc := currencyService.NewService(currencyRepo)
	planSvc := planService.NewService(planRepo, currencyRepo)
	userSubscriptionSvc := userSubscriptionService.NewService(userSubscriptionRepo, planRepo, userRepo)

	// Inicializar handlers
	currencyHdlr := currencyHandler.NewCurrencyHandler(currencySvc)
	planHdlr := planHandler.NewPlanHandler(planSvc)
	userSubscriptionHdlr := userSubscriptionHandler.NewUserSubscriptionHandler(userSubscriptionSvc)

	// Configurar grupo base de la API
	api := r.Group("/api")

	// Configurar rutas
	userRouter.SetupUserRoutes(api, userHdlr, authMiddleware)
	categoryRouter.SetupCategoryRoutes(api, categoryHdlr, authMiddleware)
	paymentMethodRouter.SetupPaymentMethodRoutes(api, paymentMethodHdlr, authMiddleware)
	transactionRouter.SetupTransactionRoutes(api, transactionHdlr, authMiddleware)
	currencyRouter.SetupCurrencyRoutes(api, currencyHdlr, authMiddleware)
	planRouter.SetupPlanRoutes(api, planHdlr, authMiddleware)

	// Configurar rutas de user_subscription
	userSubscriptionRouter.SetupUserSubscriptionRoutes(api, authMiddleware.Authorize(), adminMiddleware.RequireAdmin(), userSubscriptionHdlr)

	// Configurar rutas de health check (no requieren autenticación)
	healthRouter.SetupHealthRoutes(api, healthHdlr)

	// Configurar solo algunas rutas seleccionadas en la raíz para compatibilidad con Swagger
	rootApi := r.Group("")
	healthRouter.SetupHealthRoutes(rootApi, healthHdlr)
	currencyRouter.SetupCurrencyRoutes(rootApi, currencyHdlr, authMiddleware)
	planRouter.SetupPlanRoutes(rootApi, planHdlr, authMiddleware)
	userRouter.SetupUserRoutes(rootApi, userHdlr, authMiddleware)

	// Redireccionar peticiones a /categories hacia /api/categories para compatibilidad
	rootApi.GET("/categories", func(c *gin.Context) {
		c.Redirect(http.StatusTemporaryRedirect, "/api/categories")
	})
	rootApi.GET("/categories/:id", func(c *gin.Context) {
		id := c.Param("id")
		c.Redirect(http.StatusTemporaryRedirect, "/api/categories/"+id)
	})
	rootApi.GET("/categories/type/:type", func(c *gin.Context) {
		typeParam := c.Param("type")
		c.Redirect(http.StatusTemporaryRedirect, "/api/categories/type/"+typeParam)
	})

	// Redireccionar peticiones a /transactions hacia /api/transactions para compatibilidad
	rootApi.GET("/transactions", func(c *gin.Context) {
		c.Redirect(http.StatusTemporaryRedirect, "/api/transactions")
	})
	rootApi.POST("/transactions", func(c *gin.Context) {
		c.Redirect(http.StatusTemporaryRedirect, "/api/transactions")
	})
	rootApi.GET("/transactions/:id", func(c *gin.Context) {
		id := c.Param("id")
		c.Redirect(http.StatusTemporaryRedirect, "/api/transactions/"+id)
	})
	rootApi.PUT("/transactions/:id", func(c *gin.Context) {
		id := c.Param("id")
		c.Redirect(http.StatusTemporaryRedirect, "/api/transactions/"+id)
	})
	rootApi.DELETE("/transactions/:id", func(c *gin.Context) {
		id := c.Param("id")
		c.Redirect(http.StatusTemporaryRedirect, "/api/transactions/"+id)
	})

	// No incluir las rutas de categorías en la raíz para evitar respuestas duplicadas
	// categoryRouter.SetupCategoryRoutes(rootApi, categoryHdlr, authMiddleware)
	paymentMethodRouter.SetupPaymentMethodRoutes(rootApi, paymentMethodHdlr, authMiddleware)
	userSubscriptionRouter.SetupUserSubscriptionRoutes(rootApi, authMiddleware.Authorize(), adminMiddleware.RequireAdmin(), userSubscriptionHdlr)
}
