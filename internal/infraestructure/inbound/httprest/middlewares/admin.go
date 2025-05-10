package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// AdminMiddleware maneja la verificación de permisos de administrador
type AdminMiddleware struct{}

// NewAdminMiddleware crea un nuevo AdminMiddleware
func NewAdminMiddleware() *AdminMiddleware {
	return &AdminMiddleware{}
}

// RequireAdmin es un middleware para requerir permisos de administrador
func (m *AdminMiddleware) RequireAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		// En una implementación real, verificaríamos si el usuario tiene rol de administrador
		// consultando en la base de datos o verificando claims del token.
		// Aquí simplemente validamos que el usuario esté autenticado y establecemos un valor

		_, exists := c.Get(UserIDKey)
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autenticado"})
			return
		}

		// Por ahora, establecemos admin=true para desarrollo
		// TODO: Implementar lógica real de verificación de admin
		c.Set(IsAdminKey, true)

		c.Next()
	}
}
