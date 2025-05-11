package swagger

import (
	"fmt"
	"log"
	"strings"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "MyMoneyBackend/internal/infraestructure/inbound/httprest/docs" // Importar documentación generada
)

// SetupSwaggerRoutes configura las rutas para la documentación Swagger
func SetupSwaggerRoutes(r *gin.Engine) {
	// Configuración de rutas de Swagger con opciones personalizadas
	swaggerHandler := ginSwagger.WrapHandler(swaggerFiles.Handler)

	// Middleware para interceptar y modificar solicitudes que vienen del Swagger UI
	// para asegurarse de que los tokens de autorización tengan el formato correcto
	authMiddleware := func(c *gin.Context) {
		// Solo interceptar solicitudes que no sean del UI de Swagger
		if !strings.Contains(c.Request.URL.Path, "/swagger/") {
			authHeader := c.GetHeader("Authorization")
			if authHeader != "" && !strings.HasPrefix(authHeader, "Bearer ") {
				// Si hay un token pero no tiene el prefijo Bearer, añadirlo
				c.Request.Header.Set("Authorization", "Bearer "+authHeader)
			}
		}
		c.Next()
	}

	// Aplicar el middleware a nivel global para afectar todas las solicitudes
	r.Use(authMiddleware)

	// Configurar Swagger UI
	r.GET("/swagger/*any", swaggerHandler)

	// Redirigir /swagger a /swagger/index.html para facilidad de uso
	r.GET("/swagger", func(c *gin.Context) {
		c.Redirect(301, "/swagger/index.html")
	})

	// Configurar redirecciones para que funcionen las operaciones desde Swagger UI
	r.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path
		method := c.Request.Method

		log.Printf("NoRoute handler: %s %s", method, path)

		// Verificar si la ruta ya comienza con /api para no duplicarlo
		if strings.HasPrefix(path, "/api/") {
			// La ruta ya tiene el prefijo /api, pero no se encontró
			// Esto probablemente significa que la ruta no existe realmente
			c.JSON(404, gin.H{
				"error":      fmt.Sprintf("Ruta no encontrada: %s %s", method, path),
				"suggestion": "Verifica que la ruta sea correcta y que el método HTTP sea el adecuado",
			})
			return
		}

		// Si la ruta no comienza con /api, intentamos redirigirla
		if len(path) > 0 && path[0] == '/' {
			apiPath := "/api" + path
			log.Printf("Redirigiendo a: %s", apiPath)

			// Modificar la URL de la petición para que apunte a /api/...
			c.Request.URL.Path = apiPath

			// Intentar manejar la petición con la nueva ruta
			r.HandleContext(c)
			return
		}

		// Si no se ha manejado, mostrar un mensaje de error personalizado
		c.JSON(404, gin.H{
			"error":      fmt.Sprintf("Ruta no encontrada: %s %s", method, path),
			"suggestion": "Verifica la documentación de la API para conocer las rutas disponibles",
		})
	})
}
