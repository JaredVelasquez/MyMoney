package middleware

import (
	"net/http"
	"strings"

	auth "MyMoneyBackend/internal/application/auth"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware handles authentication
type AuthMiddleware struct {
	tokenService *auth.TokenService
}

// NewAuthMiddleware creates a new AuthMiddleware
func NewAuthMiddleware(tokenService *auth.TokenService) *AuthMiddleware {
	return &AuthMiddleware{
		tokenService: tokenService,
	}
}

// Authorize is a middleware to authorize requests
func (m *AuthMiddleware) Authorize() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "authorization header is required"})
			return
		}

		// Check for Bearer prefix
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "authorization header format must be Bearer {token}"})
			return
		}

		// Parse and validate the token
		tokenString := parts[1]
		claims, err := m.tokenService.ValidateToken(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		// Set user ID and email in context
		c.Set(UserIDKey, claims.UserID)
		c.Set("email", claims.Email)

		c.Next()
	}
}
