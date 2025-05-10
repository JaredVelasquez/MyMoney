package routers

import (
	"database/sql"
	"log"

	"github.com/gin-gonic/gin"

	"mi-app-backend/db/config"
	"mi-app-backend/internal/application/auth"
	categoryService "mi-app-backend/internal/application/category"
	currencyService "mi-app-backend/internal/application/currency"
	paymentMethodService "mi-app-backend/internal/application/paymentmethod"
	planService "mi-app-backend/internal/application/plan"
	transactionService "mi-app-backend/internal/application/transaction"
	userService "mi-app-backend/internal/application/user"
	categoryHandler "mi-app-backend/internal/infraestructure/inbound/httprest/handlers/category"
	currencyHandler "mi-app-backend/internal/infraestructure/inbound/httprest/handlers/currency"
	healthHandler "mi-app-backend/internal/infraestructure/inbound/httprest/handlers/health"
	paymentMethodHandler "mi-app-backend/internal/infraestructure/inbound/httprest/handlers/paymentmethod"
	planHandler "mi-app-backend/internal/infraestructure/inbound/httprest/handlers/plan"
	transactionHandler "mi-app-backend/internal/infraestructure/inbound/httprest/handlers/transaction"
	userHandler "mi-app-backend/internal/infraestructure/inbound/httprest/handlers/user"
	middlewares "mi-app-backend/internal/infraestructure/inbound/httprest/middlewares"
	categoryRouter "mi-app-backend/internal/infraestructure/inbound/httprest/routers/category"
	currencyRouter "mi-app-backend/internal/infraestructure/inbound/httprest/routers/currency"
	healthRouter "mi-app-backend/internal/infraestructure/inbound/httprest/routers/health"
	paymentMethodRouter "mi-app-backend/internal/infraestructure/inbound/httprest/routers/paymentmethod"
	planRouter "mi-app-backend/internal/infraestructure/inbound/httprest/routers/plan"
	transactionRouter "mi-app-backend/internal/infraestructure/inbound/httprest/routers/transaction"
	userRouter "mi-app-backend/internal/infraestructure/inbound/httprest/routers/user"
	"mi-app-backend/internal/infraestructure/inbound/httprest/swagger"
	"mi-app-backend/internal/infraestructure/outbound/repository"
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

	// Inicializar servicios
	currencySvc := currencyService.NewService(currencyRepo)
	planSvc := planService.NewService(planRepo, currencyRepo)

	// Inicializar handlers
	currencyHdlr := currencyHandler.NewCurrencyHandler(currencySvc)
	planHdlr := planHandler.NewPlanHandler(planSvc)

	// Configurar grupo base de la API
	api := r.Group("/api")

	// Configurar rutas
	userRouter.SetupUserRoutes(api, userHdlr, authMiddleware)
	categoryRouter.SetupCategoryRoutes(api, categoryHdlr, authMiddleware)
	paymentMethodRouter.SetupPaymentMethodRoutes(api, paymentMethodHdlr, authMiddleware)
	transactionRouter.SetupTransactionRoutes(api, transactionHdlr, authMiddleware)
	currencyRouter.SetupCurrencyRoutes(api, currencyHdlr, authMiddleware)
	planRouter.SetupPlanRoutes(api, planHdlr, authMiddleware)

	// Configurar rutas de health check (no requieren autenticación)
	healthRouter.SetupHealthRoutes(api, healthHdlr)

	// Configurar las mismas rutas en la raíz para que coincidan con Swagger
	rootApi := r.Group("")
	healthRouter.SetupHealthRoutes(rootApi, healthHdlr)
	currencyRouter.SetupCurrencyRoutes(rootApi, currencyHdlr, authMiddleware)
	planRouter.SetupPlanRoutes(rootApi, planHdlr, authMiddleware)
}
