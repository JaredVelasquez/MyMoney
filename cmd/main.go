package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"mi-app-backend/db/config"
	"mi-app-backend/internal/application/auth"
	categoryService "mi-app-backend/internal/application/category"
	paymentMethodService "mi-app-backend/internal/application/paymentmethod"
	transactionService "mi-app-backend/internal/application/transaction"
	userService "mi-app-backend/internal/application/user"
	"mi-app-backend/internal/domain/ports/app"
	"mi-app-backend/internal/infraestructure/inbound/httprest/routers"
	"mi-app-backend/internal/infraestructure/outbound/repository"
)

// @title Mi App Backend API
// @version 1.0
// @description API para la aplicación de finanzas personales
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.miapp.com/support
// @contact.email support@miapp.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /
// @schemes http https

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Ingrese el token con el prefijo Bearer, por ejemplo: 'Bearer abcde12345'

func main() {
	// Cargar variables de entorno
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: Error loading .env file")
	}

	// Configurar modo de Gin
	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Inicializar conexión a la base de datos
	dbConn, err := config.NewConnection()
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	db := dbConn.GetDB()

	// Inicializar repositorios
	userRepo := repository.NewUserRepository(db)
	var categoryRepo app.CategoryRepository = repository.NewCategoryRepository(db)
	var paymentMethodRepo app.PaymentMethodRepository = repository.NewPaymentMethodRepository(db)
	var transactionRepo app.TransactionRepository = repository.NewTransactionRepository(db)

	// Inicializar servicios
	tokenService := auth.NewTokenService()
	userSvc := userService.NewUserService(userRepo)
	categorySvc := categoryService.NewService(categoryRepo)
	paymentMethodSvc := paymentMethodService.NewService(paymentMethodRepo)
	transactionSvc := transactionService.NewService(transactionRepo)

	// Inicializar router
	r := gin.Default()

	// Configurar rutas de la API
	routers.SetupRouter(r, userSvc, categorySvc, paymentMethodSvc, transactionSvc, tokenService)

	// Iniciar servidor
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	log.Printf("Swagger documentation available at http://localhost:%s/swagger/index.html", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
